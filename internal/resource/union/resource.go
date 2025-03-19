package union

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/dmalch/terraform-provider-geni/internal/config"
	"github.com/dmalch/terraform-provider-geni/internal/geni"
)

type Resource struct {
	resource.ResourceWithConfigure
	accessToken types.String
}

func NewUnionResource() resource.Resource {
	return &Resource{}
}

// Metadata provides the resource type name
func (r *Resource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "geni_union"
}

func (r *Resource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Always perform a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	cfg, ok := req.ProviderData.(*config.GeniProviderConfig)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *config.GeniProviderConfig, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.accessToken = cfg.AccessToken
}

type ResourceModel struct {
	ID       types.String `tfsdk:"id"`
	Children types.Set    `tfsdk:"children"`
	Partners types.Set    `tfsdk:"partners"`
}

// Create creates the resource
func (r *Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	partnerIds := make([]types.String, 0, len(plan.Partners.Elements()))
	diag := plan.Partners.ElementsAs(ctx, &partnerIds, false)
	if diag.HasError() {
		resp.Diagnostics.Append(diag...)
		return
	}

	// If there are two partners, we can create a union by calling the profile/add-partner API
	if len(plan.Partners.Elements()) == 2 {
		// It is impossible to create a union from two existing profiles using the API,
		// so we need to create a temporary partner profile and then merge it with the
		// existing second partner profile.

		tmpProfile, err := geni.AddPartner(r.accessToken.ValueString(), partnerIds[0].ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Error adding partner", err.Error())
			return
		}

		// Merge the temporary profile with the second partner
		if err := geni.MergeProfiles(r.accessToken.ValueString(), partnerIds[1].ValueString(), tmpProfile.Id); err != nil {
			resp.Diagnostics.AddError("Error merging profiles", err.Error())
			return
		}

		plan.ID = types.StringValue(tmpProfile.Unions[0])
	}

	// Set the children. If the union already exists and has children, we can set
	// them by calling the union/add-child API. If not, we can use profile/add-child
	// on a parent profile.
	if len(plan.Children.Elements()) > 0 {

		childrenIds := make([]types.String, 0, len(plan.Children.Elements()))
		diag := plan.Children.ElementsAs(ctx, &childrenIds, false)
		if diag.HasError() {
			resp.Diagnostics.Append(diag...)
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
				tmpProfile, err = geni.AddChild(r.accessToken.ValueString(), plan.ID.ValueString())
				if err != nil {
					resp.Diagnostics.AddError("Error adding child", err.Error())
					return
				}
			} else {
				// When one parent is known, we can add a child to the parent
				if len(partnerIds) > 0 {
					// It is impossible to add an existing child profile to a parent using the API,
					// so we need to create a temporary child profile and then merge it with the
					// existing child profile.
					var err error
					tmpProfile, err = geni.AddChild(r.accessToken.ValueString(), partnerIds[0].ValueString())
					if err != nil {
						resp.Diagnostics.AddError("Error adding child", err.Error())
						return
					}
				} else if len(partnerIds) == 0 && len(childrenIds) > 1 {
					// If there are no partners, we can add a child as a sibling to the first child
					// in the union using the union/add-sibling API.
					// It is impossible to add an existing child profile to a sibling using the API,
					// so we need to create a temporary child profile and then merge it with the
					// existing child profile.
					var err error
					tmpProfile, err = geni.AddSibling(r.accessToken.ValueString(), childrenIds[i+1].ValueString())
					if err != nil {
						resp.Diagnostics.AddError("Error adding child", err.Error())
						return
					}

					// Skip the next iteration because we already added the child
					skipNextIteration = true
				}
			}

			// Merge the temporary profile with the child profile
			if err := geni.MergeProfiles(r.accessToken.ValueString(), childId.ValueString(), tmpProfile.Id); err != nil {
				resp.Diagnostics.AddError("Error merging profiles", err.Error())
				return
			}

			if plan.ID.IsUnknown() || plan.ID.IsNull() {
				plan.ID = types.StringValue(tmpProfile.Unions[0])
			}
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

// Read reads the resource
func (r *Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	union, err := geni.GetUnion(r.accessToken.ValueString(), state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading union", err.Error())
		return
	}

	if union.Id != "" {
		state.ID = types.StringValue(union.Id)
	}
	if len(union.Children) > 0 {
		children, diag := types.SetValueFrom(ctx, types.StringType, union.Children)
		resp.Diagnostics.Append(diag...)
		if resp.Diagnostics.HasError() {
			return
		}
		state.Children = children
	}

	if len(union.Partners) > 0 {
		partners, diag := types.SetValueFrom(ctx, types.StringType, union.Partners)
		resp.Diagnostics.Append(diag...)
		if resp.Diagnostics.HasError() {
			return
		}
		state.Partners = partners
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Update updates the resource
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
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}
		knownPlanPartnerIds := hashMapFrom(planPartnerIds)

		statePartnerIds, diags := convertToSlice(ctx, state.Partners)
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}
		knownStatePartnerIds := hashMapFrom(statePartnerIds)

		for _, partnerId := range statePartnerIds {
			// If the partner is not in the plan, fail the update because we can't remove
			// partners from a union using the API
			if _, ok := knownPlanPartnerIds[partnerId.ValueString()]; !ok {
				resp.Diagnostics.AddError("Cannot remove partners", "Partners cannot be removed from a union")
				return
			}
		}

		for _, partnerId := range planPartnerIds {
			// If the partner is not in the state, we need to add it
			if _, ok := knownStatePartnerIds[partnerId.ValueString()]; !ok {
				// It is impossible to add an existing profile to a union using the API, so we
				// need to create a temporary profile and then merge it with the existing
				// profile.

				tmpProfile, err := geni.AddPartner(r.accessToken.ValueString(), plan.ID.ValueString())
				if err != nil {
					resp.Diagnostics.AddError("Error adding partner", err.Error())
					return
				}

				// Merge the temporary profile with the second partner
				if err := geni.MergeProfiles(r.accessToken.ValueString(), partnerId.ValueString(), tmpProfile.Id); err != nil {
					resp.Diagnostics.AddError("Error merging profiles", err.Error())
					return
				}
			}
		}
	}

	// Check if children were updated
	if !plan.Children.Equal(state.Children) {
		planChildIds, diags := convertToSlice(ctx, plan.Children)
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}
		knownPlanChildIds := hashMapFrom(planChildIds)

		stateChildIds, diags := convertToSlice(ctx, state.Children)
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}

		for _, childId := range stateChildIds {
			// If the child is not in the plan, fail the update because we can't remove
			// children from a union using the API
			if _, ok := knownPlanChildIds[childId.ValueString()]; !ok {
				resp.Diagnostics.AddError("Cannot remove children", "Children cannot be removed from a union")
				return
			}
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func hashMapFrom(slice []types.String) map[string]struct{} {
	hashMap := make(map[string]struct{}, len(slice))
	for _, elem := range slice {
		hashMap[elem.ValueString()] = struct{}{}
	}
	return hashMap
}

func convertToSlice(ctx context.Context, set types.Set) ([]types.String, diag.Diagnostics) {
	slice := make([]types.String, 0, len(set.Elements()))
	diags := set.ElementsAs(ctx, &slice, false)
	return slice, diags
}

// Delete deletes the resource
func (r *Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// We can't delete a union, so we just remove the resource from the state

	resp.State.RemoveResource(ctx)
}
