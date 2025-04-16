package document

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/dmalch/terraform-provider-genealogy/internal/geni"
)

// Read reads the resource.
func (r *Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var documentResponse *geni.DocumentResponse
	var err error

	if r.useDocumentCache {
		documentResponse, err = r.client.GetDocumentFromCache(ctx, state.ID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Error reading document", err.Error())
			return
		}
	} else {
		documentResponse, err = r.client.GetDocument(ctx, state.ID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Error reading document", err.Error())
			return
		}
	}

	diags := ValueFrom(ctx, documentResponse, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
