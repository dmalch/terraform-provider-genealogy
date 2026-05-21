package photo

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/dmalch/terraform-provider-genealogy/internal/tfset"
)

// Update updates the resource.
func (r *Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state ResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read identity data
	var identityData ResourceIdentityModel
	if !req.Identity.Raw.IsNull() {
		resp.Diagnostics.Append(req.Identity.Get(ctx, &identityData)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	photoResponse, err := r.client.Photo().Update(ctx, plan.ID.ValueString(), RequestFrom(plan))
	if err != nil {
		resp.Diagnostics.AddError("Error updating photo", err.Error())
		return
	}

	diags := UpdateComputedFields(ctx, photoResponse, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Reconcile the tagged profiles.
	if !state.Profiles.Equal(plan.Profiles) {
		planProfileIds, diags := tfset.Strings(ctx, plan.Profiles)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		knownPlanProfileIds := tfset.Index(planProfileIds)

		stateProfileIds, diags := tfset.Strings(ctx, state.Profiles)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		knownStateProfileIds := tfset.Index(stateProfileIds)

		// Untag profiles that are no longer tagged in the photo.
		for profileId := range knownStateProfileIds {
			if _, ok := knownPlanProfileIds[profileId]; !ok {
				if _, err := r.client.Photo().Untag(ctx, photoResponse.ID, profileId); err != nil {
					resp.Diagnostics.AddError("Error untagging photo", err.Error())
					return
				}
			}
		}

		// Tag profiles that are now tagged in the photo.
		for profileId := range knownPlanProfileIds {
			if _, ok := knownStateProfileIds[profileId]; !ok {
				if _, err := r.client.Photo().Tag(ctx, photoResponse.ID, profileId); err != nil {
					resp.Diagnostics.AddError("Error tagging photo", err.Error())
					return
				}
			}
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)

	// Set data returned by API in identity
	identityData.ID = types.StringValue(photoResponse.ID)
	resp.Diagnostics.Append(resp.Identity.Set(ctx, identityData)...)
}
