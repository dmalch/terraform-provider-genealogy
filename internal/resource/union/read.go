package union

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/dmalch/terraform-provider-genealogy/internal/geni"
)

// Read reads the resource.
func (r *Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	unionResponse, err := r.client.GetUnionAsync(ctx, state.ID.ValueString())
	//unionResponse, err := r.client.GetUnion(ctx, state.ID.ValueString())
	if err != nil {
		if errors.Is(err, geni.ErrResourceNotFound) {
			resp.Diagnostics.AddWarning("Union not found", "The union was not found in the Geni API. Removing from state.")
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError("Error reading union", err.Error())
		return
	}

	// If the union doesn't have any partners or children, remove the resource from
	// the state
	if len(unionResponse.Partners) == 0 && len(unionResponse.Children) == 0 {
		existingUnionId, diags := r.findExistingUnion(ctx, state.Partners)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		if existingUnionId != "" {
			resp.Diagnostics.AddWarning("Found existing union", "The union in the state has no partners or children. Found an existing union with ID "+existingUnionId+".")
			unionResponse, err = r.client.GetUnion(ctx, existingUnionId)
			if err != nil {
				resp.Diagnostics.AddError("Error reading union", err.Error())
				return
			}
		} else {
			resp.Diagnostics.AddWarning("Union has no partners or children", "The union has no partners or children. Removing from state.")
			resp.State.RemoveResource(ctx)
			return
		}
	}

	diags := ValueFrom(ctx, unionResponse, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *Resource) findExistingUnion(ctx context.Context, partners types.Set) (string, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Attempt to find an existing union for partners in the state
	partnerIds, diags := convertToSlice(ctx, partners)
	if diags.HasError() {
		return "", diags
	}

	// If there are no partners, return an empty string
	if len(partnerIds) == 0 {
		return "", diags
	}

	// If there is only one partner, check if it has a union
	if len(partnerIds) == 1 {
		profile, err := r.client.GetProfile(ctx, partnerIds[0])
		if err != nil {
			diags.AddError("Error reading partner", err.Error())
			return "", diags
		}

		if len(profile.Unions) > 0 {
			return profile.Unions[0], diags
		}
	}

	// Get partners using the API
	profiles, err := r.client.GetProfiles(ctx, partnerIds)
	if err != nil {
		diags.AddError("Error reading partners", err.Error())
		return "", diags
	}

	if len(profiles.Results) < 2 {
		return "", diags
	}

	// Check if the partners have overlapping unions
	// Add first partner unions to a map
	unionMap := make(map[string]struct{})
	for _, union := range profiles.Results[0].Unions {
		unionMap[union] = struct{}{}
	}

	// Check if the second partner has any unions that are in the first partner
	for _, union := range profiles.Results[1].Unions {
		if _, ok := unionMap[union]; ok {
			return union, diags
		}
	}

	return "", diags
}

func (r *Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
