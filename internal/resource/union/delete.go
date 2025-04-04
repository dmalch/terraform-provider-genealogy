package union

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

	// We can't delete a union, so we just remove the resource from the state

	resp.State.RemoveResource(ctx)
}
