package union

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
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

	// Check if parents were updated
	if !plan.Partners.Equal(state.Partners) {
		planPartnerIds, diags := convertToSlice(ctx, plan.Partners)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		knownPlanPartnerIds := hashMapFrom(planPartnerIds)

		statePartnerIds, diags := convertToSlice(ctx, state.Partners)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		knownStatePartnerIds := hashMapFrom(statePartnerIds)

		for _, partnerId := range statePartnerIds {
			// If the partner is not in the plan, fail the update because we can't remove
			// partners from a union using the API
			if _, ok := knownPlanPartnerIds[partnerId]; !ok {
				resp.Diagnostics.AddAttributeWarning(path.Root(fieldPartners), "Cannot remove partners", "Partners cannot be removed from a union using terraform unless the profile is deleted")
			}
		}

		for _, partnerId := range planPartnerIds {
			// If the partner is not in the state, we need to add it
			if _, ok := knownStatePartnerIds[partnerId]; !ok {
				// It is impossible to add an existing profile to a union using the API, so we
				// need to create a temporary profile and then merge it with the existing
				// profile.

				tmpProfile, err := r.client.AddPartner(ctx, plan.ID.ValueString())
				if err != nil {
					resp.Diagnostics.AddAttributeError(path.Root(fieldPartners), "Error adding partner", err.Error())
					return
				}

				// Merge the temporary profile with the second partner
				if err := r.client.MergeProfiles(ctx, partnerId, tmpProfile.Id); err != nil {
					resp.Diagnostics.AddAttributeError(path.Root(fieldPartners), "Error merging profiles", err.Error())
					return
				}
			}
		}
	}

	// Check if children were updated
	if !plan.Children.Equal(state.Children) {
		planChildIds, diags := convertToSlice(ctx, plan.Children)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		knownPlanChildIds := hashMapFrom(planChildIds)

		stateChildIds, diags := convertToSlice(ctx, state.Children)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		knownStateChildIds := hashMapFrom(stateChildIds)

		for _, childId := range stateChildIds {
			// If the child is not in the plan, fail the update because we can't remove
			// children from a union using the API
			if _, ok := knownPlanChildIds[childId]; !ok {
				resp.Diagnostics.AddAttributeWarning(path.Root(fieldChildren), "Cannot remove children", "Children cannot be removed from a union using terraform unless the profile is deleted")
			}
		}

		for _, childId := range planChildIds {
			// If the child is not in the state, we need to add it
			if _, ok := knownStateChildIds[childId]; !ok {
				// It is impossible to add an existing profile to a union using the API, so we
				// need to create a temporary profile and then merge it with the existing
				// profile.

				tmpProfile, err := r.client.AddChild(ctx, plan.ID.ValueString())
				if err != nil {
					resp.Diagnostics.AddAttributeError(path.Root(fieldChildren), "Error adding child with ID="+childId, err.Error())
					return
				}

				// Merge the temporary profile with the child profile
				if err := r.client.MergeProfiles(ctx, childId, tmpProfile.Id); err != nil {
					resp.Diagnostics.AddAttributeError(path.Root(fieldChildren), "Error merging profiles", err.Error())
					return
				}
			}
		}
	}

	// Check if marriage or divorce were updated
	if !plan.Marriage.Equal(state.Marriage) || !plan.Divorce.Equal(state.Divorce) {
		unionRequest, diags := RequestFrom(ctx, plan)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		unionResponse, err := r.client.UpdateUnion(ctx, plan.ID.ValueString(), unionRequest)
		if err != nil {
			resp.Diagnostics.AddError("Error updating union", err.Error())
			return
		}

		diags = UpdateComputedFields(ctx, unionResponse, &plan)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}
