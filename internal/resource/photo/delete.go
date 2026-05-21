package photo

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/dmalch/go-geni"
)

// Delete deletes the resource.
func (r *Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.Photo().Delete(ctx, state.ID.ValueString())
	if err != nil && !errors.Is(err, geni.ErrResourceNotFound) {
		resp.Diagnostics.AddError("Error deleting photo", err.Error())
		return
	}

	resp.State.RemoveResource(ctx)
}
