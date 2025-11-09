package profile

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Create creates the resource.
func (r *Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	profileRequest, diags := RequestFrom(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	projectIds, diags := convertToSlice(ctx, plan.Projects)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	profileResponse, err := r.client.CreateProfile(ctx, profileRequest)
	if err != nil {
		resp.Diagnostics.AddError("Error creating profile", err.Error())
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
	identity := ResourceIdentityModel{
		ID: types.StringValue(profileResponse.Id),
	}
	resp.Diagnostics.Append(resp.Identity.Set(ctx, identity)...)
}
