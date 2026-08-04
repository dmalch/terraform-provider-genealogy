package profile

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	geniprofile "github.com/dmalch/go-geni/profile"
	"github.com/dmalch/terraform-provider-genealogy/internal/resource/event"
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

	planProjectIds, diags := tfset.Strings(ctx, plan.Projects)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	stateProjectIds, diags := tfset.Strings(ctx, state.Projects)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	projectIds := addedProjects(stateProjectIds, planProjectIds)

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

	// Link the profile to the projects the plan ADDS. See addedProjects.
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

// addedProjects returns the project ids the plan introduces — those not already
// linked in state.
//
// Re-asserting an EXISTING link is not a harmless no-op at the API. AddProfile
// is a write only the project's owner may perform, so re-sending a membership
// somebody else created fails with "access denied" and fails the whole update
// with it — even though the profile itself updated fine moments earlier.
//
// This is not hypothetical. A Нижняя Верея ancestor belongs to the public
// "Memorial: USSR political terror victims" project (project-28427), added
// there by Memorial's own curators. Terraform reads that membership on refresh,
// so it sits in both state and plan, and before this every later edit to the
// profile re-sent the link and errored. The profile's own update had already
// succeeded, which made the failure read as data loss when it was not.
//
// Create still links every project in the plan, correctly: a profile that did
// not exist a moment ago has no memberships to re-assert.
//
// Note this does not UNLINK projects the plan drops — Update never has, and
// changing that is a separate decision about who owns the field.
func addedProjects(state, plan []string) []string {
	linked := tfset.Index(state)

	var added []string
	for _, id := range plan {
		if _, exists := linked[id]; !exists {
			added = append(added, id)
			// Fold it in so a repeated id is not linked twice. plan.Projects
			// is a Set today, so this cannot bite from the Update path — but
			// the signature takes a slice, and one link op per id is the
			// contract regardless of who calls it.
			linked[id] = struct{}{}
		}
	}
	return added
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
