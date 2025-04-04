package profile

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// Update updates the resource.
func (r *Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	profileRequest, diags := RequestFrom(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	profileResponse, err := r.client.UpdateProfile(ctx, plan.ID.ValueString(), profileRequest)
	if err != nil {
		resp.Diagnostics.AddError("Error updating profile", err.Error())
		return
	}

	diags = UpdateComputedFields(ctx, profileResponse, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}
