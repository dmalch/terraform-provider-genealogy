package profile

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/dmalch/terraform-provider-genealogy/internal/geni"
)

// Read reads the resource.
func (r *Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var profileResponse *geni.ProfileResponse
	var err error

	profileResponse, err = r.getProfile(ctx, state.ID.ValueString())
	if err != nil {
		if errors.Is(err, geni.ErrResourceNotFound) {
			resp.Diagnostics.AddWarning("Resource not found", "The profile was not found in the Geni API.")
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError("Error reading profile", err.Error())
		return
	}

	// If the profile is deleted, check if it was merged into another profile and
	// read that profile instead. Iterate up to 10 times to find the merged profile.
	if state.AutoUpdateWhenMerged.ValueBool() && profileResponse.Deleted {
		for i := 0; i < 10 && profileResponse.Deleted && profileResponse.MergedInto != nil && *profileResponse.MergedInto != ""; i++ {
			profileResponse, err = r.client.GetProfile(ctx, *profileResponse.MergedInto)
			if err != nil {
				resp.Diagnostics.AddError("Error reading profile", err.Error())
				return
			}
		}
	}

	if profileResponse.Deleted {
		if profileResponse.MergedInto != nil && *profileResponse.MergedInto != "" {
			resp.Diagnostics.AddWarning(fmt.Sprintf("Resource %s is merged", profileResponse.Id),
				fmt.Sprintf("The profile %s was merged into another profile. Please use the merged profile ID=%s.", profileResponse.Id, *profileResponse.MergedInto))
		} else {
			resp.Diagnostics.AddWarning("Resource is deleted",
				fmt.Sprintf("The profile %s was deleted in the Geni API.", profileResponse.Id))
		}
		resp.State.RemoveResource(ctx)
		return
	}

	newState := state
	diags := ValueFrom(ctx, profileResponse, &newState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If names in the new state are empty, and the names in the state contain one
	// element for en-US, then use the state names.
	if len(newState.Names.Elements()) == 0 {
		// Get names from the current state
		names, diags := NameModelsFrom(ctx, state.Names)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		if _, ok := names["en-US"]; ok && len(names) == 1 {
			newState.Names = state.Names
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, newState)...)
}

func (r *Resource) getProfile(ctx context.Context, profileId string) (*geni.ProfileResponse, error) {
	if r.useProfileCache {
		return r.client.GetProfileFromCache(ctx, profileId)
	}

	return r.client.GetProfile(ctx, profileId)
}

func (r *Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
