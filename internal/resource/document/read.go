package document

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
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
	state, identity, diags := resolveDocumentImport(ctx, req.ID, r.getDocument)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
	resp.Diagnostics.Append(resp.Identity.Set(ctx, identity)...)
}

// resolveDocumentImport validates that the imported document exists on Geni and
// produces the state and identity values to write. Returning a not-found diagnostic
// here — rather than silently passing the user-supplied ID through to a later Read —
// enforces the domain rule that an imported resource must exist; otherwise the
// concurrent batch-read path would write a zombie state row that fails refresh
// forever (see GitHub issue #80).
func resolveDocumentImport(
	ctx context.Context,
	id string,
	fetch func(context.Context, string) (*geni.DocumentResponse, error),
) (ResourceModel, ResourceIdentityModel, diag.Diagnostics) {
	var state ResourceModel
	var identity ResourceIdentityModel
	var diags diag.Diagnostics

	documentResponse, err := fetch(ctx, id)
	if err != nil {
		if errors.Is(err, geni.ErrResourceNotFound) {
			diags.AddError(
				"Document not found",
				fmt.Sprintf("No Geni document with ID %q exists.", id),
			)
			return state, identity, diags
		}
		diags.AddError("Error reading document for import", err.Error())
		return state, identity, diags
	}

	diags.Append(ValueFrom(ctx, documentResponse, &state)...)
	if diags.HasError() {
		return state, identity, diags
	}

	identity.ID = types.StringValue(documentResponse.Id)
	return state, identity, diags
}
