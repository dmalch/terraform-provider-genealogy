package union

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	geniprofile "github.com/dmalch/go-geni/profile"
	"github.com/dmalch/terraform-provider-genealogy/internal/tfset"
)

// Create creates the resource.
func (r *Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// A union is built by creating temporary profiles and merging them into the
	// real ones (Geni cannot link existing profiles directly). If a step fails
	// after the union exists, persist it so Terraform tracks the partial union
	// instead of stranding it and creating another on the next apply.
	defer func() {
		if resp.Diagnostics.HasError() {
			persistPartialUnion(ctx, resp, plan)
		}
	}()

	partnerIds, diags := tfset.Strings(ctx, plan.Partners)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If there are two partners, we can create a union by calling the profile/add-partner API
	if len(plan.Partners.Elements()) == 2 {
		// It is impossible to create a union from two existing profiles using the API,
		// so we create a temporary partner profile and merge it with the existing
		// second partner profile.
		tmpProfile, err := r.addAndMerge(ctx, partnerIds[1], func(ctx context.Context) (*geniprofile.Profile, error) {
			return r.client.Profile().AddPartner(ctx, partnerIds[0])
		})
		plan.ID = unionIDFrom(plan.ID, tmpProfile)
		if err != nil {
			resp.Diagnostics.AddAttributeError(path.Root(fieldPartners), "Error adding partner", err.Error())
			return
		}
	}

	// Set the children. If the union already exists and has children, we can set
	// them by calling the union/add-child API. If not, we can use profile/add-child
	// on a parent profile.
	childrenIds, diags := tfset.Strings(ctx, plan.Children)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	fosterIds, diags := tfset.Strings(ctx, plan.FosterChildren)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	adoptedIds, diags := tfset.Strings(ctx, plan.AdoptedChildren)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	allChildrenIds := make([]string, 0, len(childrenIds)+len(fosterIds)+len(adoptedIds))
	allChildrenIds = append(allChildrenIds, childrenIds...)
	allChildrenIds = append(allChildrenIds, fosterIds...)
	allChildrenIds = append(allChildrenIds, adoptedIds...)

	fosterSet := tfset.Index(fosterIds)
	adoptedSet := tfset.Index(adoptedIds)

	if len(allChildrenIds) > 0 {
		// When the union already exists, fetch its live membership once so a
		// child that Geni already migrated onto it (e.g. via an auto-merge of a
		// duplicate union) is skipped instead of re-added — a redundant add fails
		// with "access denied" and taints the resource (#138).
		var liveResolvedID string
		var liveChildren map[string]struct{}
		membersFetched := false

		var skipNextIteration bool
		for i, childId := range allChildrenIds {
			if skipNextIteration {
				skipNextIteration = false
				continue
			}

			modifier := modifierFor(childId, fosterSet, adoptedSet)

			// It is impossible to add an existing child profile to a union using
			// the API, so create a temporary child profile and merge it with the
			// existing child profile. addAndMerge deletes the temp profile if the
			// merge fails.
			var tmpProfile *geniprofile.Profile
			var err error
			switch {
			case !plan.ID.IsUnknown() && !plan.ID.IsNull():
				// The union already exists — add the child to it.
				if !membersFetched {
					var diags diag.Diagnostics
					liveResolvedID, _, liveChildren, diags = r.currentUnionMembers(ctx, plan.ID.ValueString(), plan.Partners)
					resp.Diagnostics.Append(diags...)
					if resp.Diagnostics.HasError() {
						return
					}
					membersFetched = true
				}
				if _, ok := liveChildren[childId]; ok {
					// Already a child of the live union — nothing to add.
					continue
				}
				tmpProfile, err = r.addAndMerge(ctx, childId, func(ctx context.Context) (*geniprofile.Profile, error) {
					return r.client.Union().AddChild(ctx, liveResolvedID, geniprofile.WithModifier(modifier))
				})
			case len(partnerIds) > 0:
				// When one parent is known, add the child to the parent.
				tmpProfile, err = r.addAndMerge(ctx, childId, func(ctx context.Context) (*geniprofile.Profile, error) {
					return r.client.Profile().AddChild(ctx, partnerIds[0], geniprofile.WithModifier(modifier))
				})
			case len(allChildrenIds) > 1:
				// With no partners, add the child as a sibling of the next child.
				tmpProfile, err = r.addAndMerge(ctx, childId, func(ctx context.Context) (*geniprofile.Profile, error) {
					return r.client.Profile().AddSibling(ctx, allChildrenIds[i+1], geniprofile.WithModifier(modifier))
				})
				// Skip the next iteration because we already added that child.
				skipNextIteration = true
			default:
				// A single child with no parents cannot be attached; ValidateConfig
				// rejects this configuration before Create runs.
				continue
			}

			plan.ID = unionIDFrom(plan.ID, tmpProfile)
			if err != nil {
				resp.Diagnostics.AddAttributeError(path.Root(fieldChildren), "Error adding child with ID="+childId, err.Error())
				return
			}
		}
	}

	if !plan.Marriage.IsUnknown() && !plan.Marriage.IsNull() {
		unionRequest, diags := RequestFrom(ctx, plan)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
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
	identity := ResourceIdentityModel{
		ID: plan.ID,
	}
	resp.Diagnostics.Append(resp.Identity.Set(ctx, identity)...)
}
