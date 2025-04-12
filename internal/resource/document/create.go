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

	// Get the planned profiles to tag the document with
	profileIds, diags := convertToSlice(ctx, plan.Profiles)
	resp.Diagnostics.Append(diags...)
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

	// Tag the document with the person IDs
	for _, profileId := range profileIds {
		if _, err = r.client.TagDocument(ctx, documentResponse.Id, profileId); err != nil {
			resp.Diagnostics.AddError("Error tagging document", err.Error())
			return
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}
