package document

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// Update updates the resource.
func (r *Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state ResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
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

	// Check if profile IDs have changed
	if !state.Profiles.Equal(plan.Profiles) {
		planProfileIds, diags := convertToSlice(ctx, plan.Profiles)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		knownPlanProfileIds := hashMapFrom(planProfileIds)

		stateProfileIds, diags := convertToSlice(ctx, state.Profiles)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		knownStateProfileIds := hashMapFrom(stateProfileIds)

		// Untag profiles that are no longer associated with the document
		for profileId := range knownStateProfileIds {
			if _, ok := knownPlanProfileIds[profileId]; !ok {
				if _, err = r.client.UntagDocument(ctx, documentResponse.Id, profileId); err != nil {
					resp.Diagnostics.AddError("Error untagging document", err.Error())
					return
				}
			}
		}

		// Tag profiles that are now associated with the document
		for profileId := range knownPlanProfileIds {
			if _, ok := knownStateProfileIds[profileId]; !ok {
				if _, err = r.client.TagDocument(ctx, documentResponse.Id, profileId); err != nil {
					resp.Diagnostics.AddError("Error tagging document", err.Error())
					return
				}
			}
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}
