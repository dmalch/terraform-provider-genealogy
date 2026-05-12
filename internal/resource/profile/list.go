package profile

import (
	"context"
	"fmt"
	"iter"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/list"
	listschema "github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/dmalch/go-geni"
	"github.com/dmalch/terraform-provider-genealogy/internal/config"
	"github.com/dmalch/terraform-provider-genealogy/internal/listresource"
)

var _ list.ListResource = (*listResource)(nil)
var _ list.ListResourceWithConfigure = (*listResource)(nil)

type listResource struct {
	client *geni.Client
}

func NewListResource() list.ListResource {
	return &listResource{}
}

func (r *listResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	// Must match the managed resource type name so `terraform query` knows
	// which schema to use for the result Identity and Resource.
	resp.TypeName = req.ProviderTypeName + "_profile"
}

func (r *listResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, resp *list.ListResourceSchemaResponse) {
	// No filters in v1 — practitioners filter the query output with HCL
	// `for` expressions. Filtering surface can grow without breaking changes.
	resp.Schema = listschema.Schema{}
}

func (r *listResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	cfg, ok := req.ProviderData.(*config.ClientData)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected ListResource Configure Type",
			fmt.Sprintf("Expected *config.ClientData, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	r.client = cfg.Client
}

func (r *listResource) List(ctx context.Context, req list.ListRequest, stream *list.ListResultsStream) {
	stream.Results = streamManagedProfiles(ctx, r.client, req)
}

func streamManagedProfiles(ctx context.Context, c *geni.Client, req list.ListRequest) iter.Seq[list.ListResult] {
	return listresource.Paginate(ctx,
		func(ctx context.Context, page int) ([]geni.ProfileResponse, int, error) {
			bulk, err := c.GetManagedProfiles(ctx, page)
			if err != nil {
				return nil, 0, err
			}
			return bulk.Results, bulk.TotalCount, nil
		},
		func(err error) list.ListResult {
			return list.ListResult{Diagnostics: diag.Diagnostics{
				diag.NewErrorDiagnostic("Error listing profiles", err.Error()),
			}}
		},
		func(p geni.ProfileResponse) (list.ListResult, bool) {
			return buildListResult(ctx, &p, req)
		})
}

// displayNameFor produces a human-readable label for a profile in query output.
// The Geni API returns both a localized `names` map and flat top-level
// name fields; we prefer the localized en-US entry when present and fall
// back to the flat fields otherwise. When nothing is populated the bare ID
// is returned so the result is still visually identifiable.
func displayNameFor(p *geni.ProfileResponse) string {
	first, last := canonicalFirstLast(p)
	switch {
	case first != "" && last != "":
		return fmt.Sprintf("%s %s (%s)", first, last, p.Id)
	case first != "":
		return fmt.Sprintf("%s (%s)", first, p.Id)
	case last != "":
		return fmt.Sprintf("%s (%s)", last, p.Id)
	default:
		return p.Id
	}
}

func canonicalFirstLast(p *geni.ProfileResponse) (string, string) {
	if name, ok := p.Names["en-US"]; ok {
		return stringFromPtr(name.FirstName), stringFromPtr(name.LastName)
	}
	return stringFromPtr(p.FirstName), stringFromPtr(p.LastName)
}

func stringFromPtr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// buildListResult turns one API response into a list.ListResult whose Identity
// carries the profile ID under the managed resource's identity schema. When
// req.IncludeResource is true the Resource field is populated via ValueFrom —
// the same translator used by Read — so list output round-trips through
// `import { identity = ... }`.
//
// Returning false signals the caller (Paginate) to stop iteration; in that
// case the result's Diagnostics already carry the offending error.
func buildListResult(
	ctx context.Context,
	resp *geni.ProfileResponse,
	req list.ListRequest,
) (list.ListResult, bool) {
	result := req.NewListResult(ctx)

	identity := ResourceIdentityModel{
		ID: types.StringValue(resp.Id),
	}
	diags := result.Identity.Set(ctx, identity)
	result.Diagnostics.Append(diags...)
	if result.Diagnostics.HasError() {
		return result, false
	}

	result.DisplayName = displayNameFor(resp)

	if req.IncludeResource {
		// Seed every collection field with a typed null so ValueFrom can
		// overwrite the fields the API does return without leaving any
		// untouched attribute carrying a type-less zero `types.Set{}` /
		// `types.Map{}` / `types.Object{}`, which the framework rejects with
		// a MISSING TYPE conversion error.
		model := NewEmptyResourceModel()

		diags = ValueFrom(ctx, resp, &model)
		result.Diagnostics.Append(diags...)
		if result.Diagnostics.HasError() {
			return result, false
		}

		diags = result.Resource.Set(ctx, model)
		result.Diagnostics.Append(diags...)
		if result.Diagnostics.HasError() {
			return result, false
		}
	}

	return result, true
}
