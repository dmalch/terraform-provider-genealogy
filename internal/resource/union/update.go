package union

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"

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

	// Check if parents were updated
	if !plan.Partners.Equal(state.Partners) {
		planPartnerIds, diags := tfset.Strings(ctx, plan.Partners)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		knownPlanPartnerIds := tfset.Index(planPartnerIds)

		statePartnerIds, diags := tfset.Strings(ctx, state.Partners)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		knownStatePartnerIds := tfset.Index(statePartnerIds)

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
				// It is impossible to add an existing profile to a union using the
				// API, so create a temporary profile and merge it with the existing
				// one. addAndMerge deletes the temp profile if the merge fails.
				if _, err := r.addAndMerge(ctx, partnerId, func(ctx context.Context) (*geniprofile.Profile, error) {
					return r.client.Union().AddPartner(ctx, plan.ID.ValueString())
				}); err != nil {
					resp.Diagnostics.AddAttributeError(path.Root(fieldPartners), "Error adding partner", err.Error())
					return
				}
			}
		}
	}

	// Warn on modifier changes: Geni has no endpoint to re-tag an existing child.
	for _, m := range childrenWithChangedModifier(ctx, plan, state) {
		resp.Diagnostics.AddAttributeWarning(path.Root(fieldChildren),
			"Cannot change relationship modifier",
			"Profile "+m.id+" cannot be moved from "+m.from+" to "+m.to+
				" via the Geni API. Re-tag the relationship on Geni.com, then re-run terraform.",
		)
	}

	// Check if any of the three child sets were updated
	if !plan.Children.Equal(state.Children) ||
		!plan.FosterChildren.Equal(state.FosterChildren) ||
		!plan.AdoptedChildren.Equal(state.AdoptedChildren) {

		planBio, diags := tfset.Strings(ctx, plan.Children)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		planFoster, diags := tfset.Strings(ctx, plan.FosterChildren)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		planAdopted, diags := tfset.Strings(ctx, plan.AdoptedChildren)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		stateBio, diags := tfset.Strings(ctx, state.Children)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		stateFoster, diags := tfset.Strings(ctx, state.FosterChildren)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		stateAdopted, diags := tfset.Strings(ctx, state.AdoptedChildren)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		planAll := append(append(append([]string{}, planBio...), planFoster...), planAdopted...)
		stateAll := append(append(append([]string{}, stateBio...), stateFoster...), stateAdopted...)

		knownPlanAll := tfset.Index(planAll)
		knownStateAll := tfset.Index(stateAll)

		for _, childId := range stateAll {
			// If the child is not in the plan, fail the update because we can't remove
			// children from a union using the API
			if _, ok := knownPlanAll[childId]; !ok {
				resp.Diagnostics.AddAttributeWarning(path.Root(fieldChildren), "Cannot remove children", "Children cannot be removed from a union using terraform unless the profile is deleted")
			}
		}

		fosterSet := tfset.Index(planFoster)
		adoptedSet := tfset.Index(planAdopted)

		for _, childId := range planAll {
			// If the child is not in the state, we need to add it
			if _, ok := knownStateAll[childId]; !ok {
				// It is impossible to add an existing profile to a union using the
				// API, so create a temporary profile and merge it with the existing
				// one. addAndMerge deletes the temp profile if the merge fails.
				modifier := modifierFor(childId, fosterSet, adoptedSet)
				if _, err := r.addAndMerge(ctx, childId, func(ctx context.Context) (*geniprofile.Profile, error) {
					return r.client.Union().AddChild(ctx, plan.ID.ValueString(), geniprofile.WithModifier(modifier))
				}); err != nil {
					resp.Diagnostics.AddAttributeError(path.Root(fieldChildren), "Error adding child with ID="+childId, err.Error())
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

		// Pre-wipe event dates that the plan partially clears — Geni's PATCH
		// deep-merges nested objects so per-key nulls inside `date` are no-ops
		// (#94). Cost is one extra PATCH per Update that hits this case.
		var wipeEvents []string
		if event.EventNeedsDatePreWipe(state.Marriage, plan.Marriage) {
			wipeEvents = append(wipeEvents, "marriage")
		}
		if event.EventNeedsDatePreWipe(state.Divorce, plan.Divorce) {
			wipeEvents = append(wipeEvents, "divorce")
		}
		if len(wipeEvents) > 0 {
			if err := r.client.Profile().WipeEventDates(ctx, plan.ID.ValueString(), wipeEvents); err != nil {
				resp.Diagnostics.AddError("Error clearing date fields", err.Error())
				return
			}
		}

		unionResponse, err := r.client.Union().Update(ctx, plan.ID.ValueString(), unionRequest)
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

	// Set data returned by API in identity
	identityData.ID = plan.ID
	resp.Diagnostics.Append(resp.Identity.Set(ctx, identityData)...)
}
