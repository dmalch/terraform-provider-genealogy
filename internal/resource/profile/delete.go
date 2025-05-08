package profile

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/dmalch/terraform-provider-genealogy/internal/geni"
)

// Delete deletes the resource.
func (r *Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteProfile(ctx, state.ID.ValueString())
	if err != nil {
		if errors.Is(err, geni.ErrAccessDenied) {
			// Check if the profile has been deleted in Geni.
			profile, err := r.client.GetProfile(ctx, state.ID.ValueString())
			if err != nil {
				resp.Diagnostics.AddError("Error reading profile", err.Error())
				return
			}

			// If the profile has been deleted in Geni, it is safe to remove it from the state.
			if profile.Deleted {
				resp.State.RemoveResource(ctx)
				return
			}
		}

		if !errors.Is(err, geni.ErrResourceNotFound) {
			resp.Diagnostics.AddError("Error deleting profile", err.Error())
			return
		}
	}

	resp.State.RemoveResource(ctx)
}
