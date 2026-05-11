package profile

import (
	"context"
	"errors"
	"fmt"

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

	// Read identity data
	var identityData ResourceIdentityModel
	if !req.Identity.Raw.IsNull() {
		resp.Diagnostics.Append(req.Identity.Get(ctx, &identityData)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	profileResponse, err := r.getProfile(ctx, state.ID.ValueString())
	if err != nil {
		if errors.Is(err, geni.ErrResourceNotFound) {
			resp.Diagnostics.AddWarning("Profile not found", "The profile was not found in the Geni API.")
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError("Error reading profile", err.Error())
		return
	}

	// If the profile is deleted, check if it was merged into another profile and
	// read that profile instead. Iterate up to 10 times to find the merged profile.
	if profileResponse.Deleted {
		if profileResponse.MergedInto == "" {
			resp.Diagnostics.AddWarning("Resource is deleted",
				fmt.Sprintf("The profile %s was deleted in the Geni API.", profileResponse.Id))
		}

		if r.autoUpdateMergedProfiles {
			profileResponse, err = FollowMergedInto(ctx, profileResponse, r.batchClient.GetProfile, 10)
			if err != nil {
				resp.Diagnostics.AddError("Error reading profile", err.Error())
				return
			}
			if profileResponse.Deleted {
				resp.State.RemoveResource(ctx)
				return
			}
		} else {
			resp.State.RemoveResource(ctx)
			return
		}
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

	// Set data returned by API in identity
	identityData.ID = types.StringValue(profileResponse.Id)
	resp.Diagnostics.Append(resp.Identity.Set(ctx, identityData)...)
}

func (r *Resource) getProfile(ctx context.Context, profileId string) (*geni.ProfileResponse, error) {
	if r.useProfileCache {
		return r.cacheClient.GetProfile(ctx, profileId)
	}

	return r.batchClient.GetProfile(ctx, profileId)
}

func (r *Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Plannable imports (Terraform 1.12+ `import { identity = { id = "..." } }`)
	// pass the ID via the typed Identity field instead of the legacy string ID.
	importID := req.ID
	if req.Identity != nil && !req.Identity.Raw.IsNull() {
		var identity ResourceIdentityModel
		resp.Diagnostics.Append(req.Identity.Get(ctx, &identity)...)
		if resp.Diagnostics.HasError() {
			return
		}
		if identity.ID.ValueString() != "" {
			importID = identity.ID.ValueString()
		}
	}

	resp.Diagnostics.Append(validateProfileImportID(ctx, importID, r.getProfile)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resource.ImportStatePassthroughWithIdentity(ctx, path.Root("id"), path.Root("id"), req, resp)
}

// validateProfileImportID round-trips the API to confirm the imported profile
// exists on Geni. See GitHub issue #80 for why this validation is required. On
// success, state population is left to the framework's follow-up Read so that
// schema-aware null defaults for collection fields are preserved.
func validateProfileImportID(
	ctx context.Context,
	id string,
	fetch func(context.Context, string) (*geni.ProfileResponse, error),
) diag.Diagnostics {
	var diags diag.Diagnostics

	response, err := fetch(ctx, id)
	if err != nil {
		if errors.Is(err, geni.ErrResourceNotFound) {
			diags.AddError(
				"Profile not found",
				fmt.Sprintf("No Geni profile with ID %q exists.", id),
			)
			return diags
		}
		diags.AddError("Error reading profile for import", err.Error())
		return diags
	}
	// The Geni single-resource endpoint sometimes returns 200 with an empty
	// body for IDs that do not exist (observed on sandbox), instead of 404.
	// Treat a missing Id field as the same domain signal as ErrResourceNotFound.
	if response == nil || response.Id == "" {
		diags.AddError(
			"Profile not found",
			fmt.Sprintf("No Geni profile with ID %q exists.", id),
		)
	}
	return diags
}
