package document

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/dmalch/terraform-provider-genealogy/internal/geni"
)

// Read reads the resource.
func (r *Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read identity data
	var identityData ResourceIdentityModel
	if !req.Identity.Raw.IsNull() {
		resp.Diagnostics.Append(req.Identity.Get(ctx, &identityData)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	documentResponse, err := r.getDocument(ctx, state.ID.ValueString())
	if err != nil {
		if errors.Is(err, geni.ErrResourceNotFound) {
			resp.Diagnostics.AddWarning("Document not found", "The document was not found in the Geni API.")
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError("Error reading document", err.Error())
		return
	}

	diags := ValueFrom(ctx, documentResponse, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)

	// Set data returned by API in identity
	identityData.ID = types.StringValue(documentResponse.Id)
	resp.Diagnostics.Append(resp.Identity.Set(ctx, identityData)...)
}

func (r *Resource) getDocument(ctx context.Context, documentId string) (*geni.DocumentResponse, error) {
	if r.useDocumentCache {
		return r.cacheClient.GetDocument(ctx, documentId)
	}

	return r.batchClient.GetDocument(ctx, documentId)
}

func (r *Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughWithIdentity(ctx, path.Root("id"), path.Root("id"), req, resp)
}
