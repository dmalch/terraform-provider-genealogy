package document

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
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
	resp.Diagnostics.Append(validateDocumentImportID(ctx, req.ID, r.getDocument)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resource.ImportStatePassthroughWithIdentity(ctx, path.Root("id"), path.Root("id"), req, resp)
}

// validateDocumentImportID round-trips the API to confirm the imported document
// exists on Geni. Without this, the framework's bare passthrough would write the
// user-supplied ID into state unchecked; combined with the concurrent batch-read
// path, that produces a zombie state row that fails refresh forever (see GitHub
// issue #80). On success, state population is left to the framework's follow-up
// Read so that schema-aware null defaults for collection fields are preserved.
func validateDocumentImportID(
	ctx context.Context,
	id string,
	fetch func(context.Context, string) (*geni.DocumentResponse, error),
) diag.Diagnostics {
	var diags diag.Diagnostics

	response, err := fetch(ctx, id)
	if err != nil {
		if errors.Is(err, geni.ErrResourceNotFound) {
			diags.AddError(
				"Document not found",
				fmt.Sprintf("No Geni document with ID %q exists.", id),
			)
			return diags
		}
		diags.AddError("Error reading document for import", err.Error())
		return diags
	}
	// The Geni single-resource endpoint sometimes returns 200 with an empty
	// body for IDs that do not exist (observed on sandbox), instead of 404.
	// Treat a missing Id field as the same domain signal as ErrResourceNotFound.
	if response == nil || response.Id == "" {
		diags.AddError(
			"Document not found",
			fmt.Sprintf("No Geni document with ID %q exists.", id),
		)
	}
	return diags
}
