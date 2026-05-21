package photo

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/dmalch/go-geni"
	geniphoto "github.com/dmalch/go-geni/photo"
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

	photoResponse, err := r.getPhoto(ctx, state.ID.ValueString())
	if err != nil {
		if errors.Is(err, geni.ErrResourceNotFound) {
			resp.Diagnostics.AddWarning("Photo not found", "The photo was not found in the Geni API. Removing from state.")
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError("Error reading photo", err.Error())
		return
	}

	diags := ValueFrom(ctx, photoResponse, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)

	// Set data returned by API in identity
	identityData.ID = types.StringValue(photoResponse.ID)
	resp.Diagnostics.Append(resp.Identity.Set(ctx, identityData)...)
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

	resp.Diagnostics.Append(validatePhotoImportID(ctx, importID, r.getPhoto)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resource.ImportStatePassthroughWithIdentity(ctx, path.Root("id"), path.Root("id"), req, resp)
}

// getPhoto fetches a photo through the batch client, which coalesces the
// concurrent reads Terraform issues when refreshing many photos at once.
func (r *Resource) getPhoto(ctx context.Context, photoId string) (*geniphoto.Photo, error) {
	return r.batchClient.GetPhoto(ctx, photoId)
}

// validatePhotoImportID round-trips the API to confirm the imported photo exists
// on Geni. See the document resource's equivalent for why this is required. On
// success, state population is left to the framework's follow-up Read.
func validatePhotoImportID(
	ctx context.Context,
	id string,
	fetch func(context.Context, string) (*geniphoto.Photo, error),
) diag.Diagnostics {
	var diags diag.Diagnostics

	response, err := fetch(ctx, id)
	if err != nil {
		if errors.Is(err, geni.ErrResourceNotFound) {
			diags.AddError("Photo not found", fmt.Sprintf("No Geni photo with ID %q exists.", id))
			return diags
		}
		diags.AddError("Error reading photo for import", err.Error())
		return diags
	}
	// The Geni single-resource endpoint sometimes returns 200 with an empty
	// body for IDs that do not exist; treat a missing ID as not-found.
	if response == nil || response.ID == "" {
		diags.AddError("Photo not found", fmt.Sprintf("No Geni photo with ID %q exists.", id))
	}
	return diags
}
