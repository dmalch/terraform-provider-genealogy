package profile

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/dmalch/terraform-provider-genealogy/internal/geni"
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

	profileRequest, diags := RequestFrom(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Check if about fields were removed in the plan
	if !plan.About.Equal(state.About) {
		removedAboutKeys := findRemovedKeys(state.About, plan.About)

		for _, removedAboutKey := range removedAboutKeys {
			profileRequest.DetailStrings[removedAboutKey] = geni.DetailsString{}
		}
	}

	projectIds, diags := convertToSlice(ctx, plan.Projects)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	profileResponse, err := r.client.UpdateProfile(ctx, plan.ID.ValueString(), profileRequest)
	if err != nil {
		resp.Diagnostics.AddError("Error updating profile", err.Error())
		return
	}

	// Link the profile to the projects if specified.
	for _, projectId := range projectIds {
		if _, err := r.client.AddProfileToProject(ctx, profileResponse.Id, projectId); err != nil {
			resp.Diagnostics.AddError("Error linking profile to project", err.Error())
			return
		}
	}

	diags = UpdateComputedFields(ctx, profileResponse, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)

	// Set data returned by API in identity
	identityData.ID = types.StringValue(profileResponse.Id)
	resp.Diagnostics.Append(resp.Identity.Set(ctx, identityData)...)
}

func findRemovedKeys(stateAbout types.Map, planAbout types.Map) []string {
	removedKeys := []string{}
	for locale, _ := range stateAbout.Elements() {
		if _, ok := planAbout.Elements()[locale]; !ok {
			removedKeys = append(removedKeys, locale)
		}
	}
	return removedKeys
}
