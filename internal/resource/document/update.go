package document

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

	documentRequest, diags := RequestFrom(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	documentResponse, err := r.client.UpdateDocument(ctx, plan.ID.ValueString(), documentRequest)
	if err != nil {
		resp.Diagnostics.AddError("Error updating document", err.Error())
		return
	}

	diags = UpdateComputedFields(ctx, documentResponse, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}
