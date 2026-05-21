package profile

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	geniprofile "github.com/dmalch/go-geni/profile"
	"github.com/dmalch/terraform-provider-genealogy/internal/resource/event"
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
			profileRequest.DetailStrings[removedAboutKey] = geniprofile.DetailsString{}
		}
	}

	projectIds, diags := convertToSlice(ctx, plan.Projects)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Geni's PATCH deep-merges nested objects per-key, so clearing individual
	// date sub-fields needs a wipe-then-rewrite (#94). Issue the pre-wipe only
	// for events where the plan keeps the date but clears at least one sub-field.
	wipeEvents := planDateWipes(state, plan)
	if len(wipeEvents) > 0 {
		if err := r.client.Profile().WipeEventDates(ctx, plan.ID.ValueString(), wipeEvents); err != nil {
			resp.Diagnostics.AddError("Error clearing date fields", err.Error())
			return
		}
	}

	profileResponse, err := r.client.Profile().Update(ctx, plan.ID.ValueString(), profileRequest)
	if err != nil {
		resp.Diagnostics.AddError("Error updating profile", err.Error())
		return
	}

	// Link the profile to the projects if specified.
	for _, projectId := range projectIds {
		if _, err := r.client.Project().AddProfile(ctx, profileResponse.ID, projectId); err != nil {
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
	identityData.ID = types.StringValue(profileResponse.ID)
	resp.Diagnostics.Append(resp.Identity.Set(ctx, identityData)...)
}

// planDateWipes lists the event keys whose date sub-object the Update path
// must pre-wipe before sending the regular PATCH. Order is fixed for
// determinism in tests and to match the API payload shape.
func planDateWipes(state, plan ResourceModel) []string {
	var wipes []string
	for _, e := range []struct {
		name        string
		state, plan types.Object
	}{
		{"birth", state.Birth, plan.Birth},
		{"baptism", state.Baptism, plan.Baptism},
		{"death", state.Death, plan.Death},
		{"burial", state.Burial, plan.Burial},
	} {
		if event.EventNeedsDatePreWipe(e.state, e.plan) {
			wipes = append(wipes, e.name)
		}
	}
	return wipes
}

func findRemovedKeys(stateAbout types.Map, planAbout types.Map) []string {
	removedKeys := []string{}
	for locale := range stateAbout.Elements() {
		if _, ok := planAbout.Elements()[locale]; !ok {
			removedKeys = append(removedKeys, locale)
		}
	}
	return removedKeys
}
