package document

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
	genidocument "github.com/dmalch/go-geni/document"
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
	resp.TypeName = req.ProviderTypeName + "_document"
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
	stream.Results = streamUploadedDocuments(ctx, r.client, req)
}

func streamUploadedDocuments(ctx context.Context, c *geni.Client, req list.ListRequest) iter.Seq[list.ListResult] {
	return listresource.Paginate(ctx,
		func(ctx context.Context, page int) ([]genidocument.Document, int, error) {
			bulk, err := c.User().UploadedDocuments(ctx, page)
			if err != nil {
				return nil, 0, err
			}
			return bulk.Results, bulk.TotalCount, nil
		},
		func(err error) list.ListResult {
			return list.ListResult{Diagnostics: diag.Diagnostics{
				diag.NewErrorDiagnostic("Error listing documents", err.Error()),
			}}
		},
		func(d genidocument.Document) (list.ListResult, bool) {
			return buildListResult(ctx, &d, req)
		})
}

// displayNameFor produces a human-readable label for a document in query
// output. The Title is the obvious choice; we fall back to the bare ID when
// it is absent so the result remains visually identifiable.
func displayNameFor(d *genidocument.Document) string {
	if d.Title == "" {
		return d.ID
	}
	return fmt.Sprintf("%s (%s)", d.Title, d.ID)
}

// buildListResult turns one API response into a list.ListResult whose Identity
// carries the document ID under the managed resource's identity schema. When
// req.IncludeResource is true the Resource field is populated via ValueFrom —
// the same translator used by Read — so list output round-trips through
// `import { identity = ... }`.
func buildListResult(
	ctx context.Context,
	resp *genidocument.Document,
	req list.ListRequest,
) (list.ListResult, bool) {
	result := req.NewListResult(ctx)

	identity := ResourceIdentityModel{
		ID: types.StringValue(resp.ID),
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
		// untouched attribute carrying a type-less zero value (see the
		// equivalent comment in the profile package for the failure mode).
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
