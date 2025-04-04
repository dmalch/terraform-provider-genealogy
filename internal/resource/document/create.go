package document

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// Create creates the resource.
func (r *Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
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

	documentResponse, err := r.client.CreateDocument(ctx, documentRequest)
	if err != nil {
		resp.Diagnostics.AddError("Error creating document", err.Error())
		return
	}

	diags = UpdateComputedFields(ctx, documentResponse, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}
