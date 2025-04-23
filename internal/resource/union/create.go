package union

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/dmalch/terraform-provider-genealogy/internal/geni"
)

// Create creates the resource.
func (r *Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	partnerIds, diags := convertToSlice(ctx, plan.Partners)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If there are two partners, we can create a union by calling the profile/add-partner API
	if len(plan.Partners.Elements()) == 2 {
		// It is impossible to create a union from two existing profiles using the API,
		// so we need to create a temporary partner profile and then merge it with the
		// existing second partner profile.

		tmpProfile, err := r.client.AddPartner(ctx, partnerIds[0])
		if err != nil {
			resp.Diagnostics.AddAttributeError(path.Root(fieldPartners), "Error adding partner", err.Error())
			return
		}

		// Merge the temporary profile with the second partner
		if err := r.client.MergeProfiles(ctx, partnerIds[1], tmpProfile.Id); err != nil {
			resp.Diagnostics.AddAttributeError(path.Root(fieldPartners), "Error merging profiles", err.Error())
			return
		}

		plan.ID = types.StringValue(tmpProfile.Unions[0])
	}

	// Set the children. If the union already exists and has children, we can set
	// them by calling the union/add-child API. If not, we can use profile/add-child
	// on a parent profile.
	if len(plan.Children.Elements()) > 0 {

		childrenIds, diags := convertToSlice(ctx, plan.Children)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		var skipNextIteration bool
		for i, childId := range childrenIds {
			if skipNextIteration {
				skipNextIteration = false
				continue
			}

			var tmpProfile *geni.ProfileResponse

			// If the union already exists, we can add children to it
			if !plan.ID.IsUnknown() && !plan.ID.IsNull() {
				// It is impossible to add an existing child profile to a union using the API, so
				// we need to create a temporary child profile and then merge it with the
				// existing child profile.
				var err error
				tmpProfile, err = r.client.AddChild(ctx, plan.ID.ValueString())
				if err != nil {
					resp.Diagnostics.AddAttributeError(path.Root(fieldChildren), "Error adding child with ID="+childId, err.Error())
					return
				}
			} else {
				// When one parent is known, we can add a child to the parent
				if len(partnerIds) > 0 {
					// It is impossible to add an existing child profile to a parent using the API,
					// so we need to create a temporary child profile and then merge it with the
					// existing child profile.
					var err error
					tmpProfile, err = r.client.AddChild(ctx, partnerIds[0])
					if err != nil {
						resp.Diagnostics.AddAttributeError(path.Root(fieldChildren), "Error adding child with ID="+childId, err.Error())
						return
					}
				} else if len(partnerIds) == 0 && len(childrenIds) > 1 {
					// If there are no partners, we can add a child as a sibling to the first child
					// in the union using the union/add-sibling API.
					// It is impossible to add an existing child profile to a sibling using the API,
					// so we need to create a temporary child profile and then merge it with the
					// existing child profile.
					var err error
					tmpProfile, err = r.client.AddSibling(ctx, childrenIds[i+1])
					if err != nil {
						resp.Diagnostics.AddAttributeError(path.Root(fieldChildren), "Error adding child with ID="+childId, err.Error())
						return
					}

					// Skip the next iteration because we already added the child
					skipNextIteration = true
				}
			}

			// Merge the temporary profile with the child profile
			if err := r.client.MergeProfiles(ctx, childId, tmpProfile.Id); err != nil {
				resp.Diagnostics.AddAttributeError(path.Root(fieldChildren), "Error merging profiles", err.Error())
				return
			}

			if plan.ID.IsUnknown() || plan.ID.IsNull() {
				plan.ID = types.StringValue(tmpProfile.Unions[0])
			}
		}
	}

	if !plan.Marriage.IsUnknown() && !plan.Marriage.IsNull() {
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
