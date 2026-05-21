package photo

import (
	"context"
	"fmt"
	"iter"
	"math"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/list"
	listschema "github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/dmalch/go-geni"
	geniphoto "github.com/dmalch/go-geni/photo"
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
	resp.TypeName = req.ProviderTypeName + "_photo"
}

func (r *listResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, resp *list.ListResourceSchemaResponse) {
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
	stream.Results = streamUploadedPhotos(ctx, r.client, req)
}

func streamUploadedPhotos(ctx context.Context, c *geni.Client, req list.ListRequest) iter.Seq[list.ListResult] {
	return listresource.Paginate(ctx,
		func(ctx context.Context, page int) ([]geniphoto.Photo, int, error) {
			bulk, err := c.User().UploadedPhotos(ctx, page)
			if err != nil {
				return nil, 0, err
			}
			// The uploaded-photos endpoint reports no total count, so return a
			// sentinel maximum; Paginate stops when a page comes back empty.
			return bulk.Results, math.MaxInt, nil
		},
		func(err error) list.ListResult {
			return list.ListResult{Diagnostics: diag.Diagnostics{
				diag.NewErrorDiagnostic("Error listing photos", err.Error()),
			}}
		},
		func(p geniphoto.Photo) (list.ListResult, bool) {
			return buildListResult(ctx, &p, req)
		})
}

// displayNameFor produces a human-readable label for a photo in query output.
// The Title is the obvious choice; it falls back to the bare ID when absent.
func displayNameFor(p *geniphoto.Photo) string {
	if p.Title == "" {
		return p.ID
	}
	return fmt.Sprintf("%s (%s)", p.Title, p.ID)
}

// buildListResult turns one API response into a list.ListResult whose Identity
// carries the photo ID under the managed resource's identity schema. When
// req.IncludeResource is true the Resource field is populated via ValueFrom —
// the same translator used by Read — so list output round-trips through
// `import { identity = ... }`.
func buildListResult(ctx context.Context, resp *geniphoto.Photo, req list.ListRequest) (list.ListResult, bool) {
	result := req.NewListResult(ctx)

	identity := ResourceIdentityModel{ID: types.StringValue(resp.ID)}
	diags := result.Identity.Set(ctx, identity)
	result.Diagnostics.Append(diags...)
	if result.Diagnostics.HasError() {
		return result, false
	}

	result.DisplayName = displayNameFor(resp)

	if req.IncludeResource {
		// Seed collection fields with typed nulls so ValueFrom can overwrite
		// only what the API returns without leaving a type-less zero value.
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
