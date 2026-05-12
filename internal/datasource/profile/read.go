package profile

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/dmalch/go-geni"
	resourceprofile "github.com/dmalch/terraform-provider-genealogy/internal/resource/profile"
)

const maxMergeHops = 10

func (d *DataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data resourceprofile.ResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	lookup := data.ID.ValueString()
	if lookup == "" {
		// Geni's URL scheme addresses profiles by guid as `profile-g<guid>`
		// (e.g. profile-g598352); the bare guid is not a routable resource path.
		lookup = "profile-g" + data.Guid.ValueString()
	}

	response, err := d.getProfile(ctx, lookup)
	if err != nil {
		if errors.Is(err, geni.ErrResourceNotFound) {
			resp.Diagnostics.AddError(
				"Profile not found",
				fmt.Sprintf("No Geni profile with identifier %q exists.", lookup),
			)
			return
		}
		resp.Diagnostics.AddError("Error reading profile", err.Error())
		return
	}
	if response == nil || response.Id == "" {
		resp.Diagnostics.AddError(
			"Profile not found",
			fmt.Sprintf("No Geni profile with identifier %q exists.", lookup),
		)
		return
	}

	if response.Deleted {
		if !d.autoUpdateMergedProfiles {
			resp.Diagnostics.AddError(
				"Profile is deleted",
				fmt.Sprintf("Profile %q is deleted on Geni. Set the provider's `auto_update_merged_profiles = true` to follow merge chains automatically.", response.Id),
			)
			return
		}
		response, err = resourceprofile.FollowMergedInto(ctx, response, d.batchClient.GetProfile, maxMergeHops)
		if err != nil {
			resp.Diagnostics.AddError("Error following merge chain", err.Error())
			return
		}
		if response.Deleted {
			resp.Diagnostics.AddError(
				"Profile is deleted with no live merge target",
				fmt.Sprintf("The merge chain starting at %q did not resolve to a live profile within %d hops.", lookup, maxMergeHops),
			)
			return
		}
	}

	state := resourceprofile.NewEmptyResourceModel()
	resp.Diagnostics.Append(resourceprofile.ValueFrom(ctx, response, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (d *DataSource) getProfile(ctx context.Context, idOrGuid string) (*geni.ProfileResponse, error) {
	if d.useProfileCache {
		return d.cacheClient.GetProfile(ctx, idOrGuid)
	}
	return d.batchClient.GetProfile(ctx, idOrGuid)
}
