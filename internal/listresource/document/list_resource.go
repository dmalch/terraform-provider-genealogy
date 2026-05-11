package documentlist

import (
	"context"
	"fmt"
	"iter"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/dmalch/terraform-provider-genealogy/internal/config"
	"github.com/dmalch/terraform-provider-genealogy/internal/geni"
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
	resp.Schema = schema.Schema{}
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
		func(ctx context.Context, page int) ([]geni.DocumentResponse, int, error) {
			bulk, err := c.GetUploadedDocuments(ctx, page)
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
		func(d geni.DocumentResponse) (list.ListResult, bool) {
			return buildDocumentListResult(ctx, &d, req)
		})
}
