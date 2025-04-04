package document

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// Delete deletes the resource.
func (r *Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteDocument(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting document", err.Error())
		return
	}

	resp.State.RemoveResource(ctx)
}
