package photo

import (
	"bytes"
	"context"
	"encoding/base64"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	geniphoto "github.com/dmalch/go-geni/photo"
	"github.com/dmalch/terraform-provider-genealogy/internal/tfset"
)

// Create creates the resource.
func (r *Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	profileIds, diags := tfset.Strings(ctx, plan.Profiles)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	raw, err := base64.StdEncoding.DecodeString(plan.File.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid photo file", "The file attribute must be a base64-encoded string: "+err.Error())
		return
	}

	var opts []geniphoto.CreateOption
	if !plan.Description.IsNull() {
		opts = append(opts, geniphoto.WithDescription(plan.Description.ValueString()))
	}
	if !plan.Date.IsNull() {
		opts = append(opts, geniphoto.WithDate(plan.Date.ValueString()))
	}
	if !plan.Album.IsNull() && !plan.Album.IsUnknown() {
		opts = append(opts, geniphoto.WithAlbum(plan.Album.ValueString()))
	}

	photoResponse, err := r.client.Photo().Create(ctx, plan.Title.ValueString(), plan.FileName.ValueString(), bytes.NewReader(raw), opts...)
	if err != nil {
		resp.Diagnostics.AddError("Error creating photo", err.Error())
		return
	}

	diags = UpdateComputedFields(ctx, photoResponse, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Tag the photo with the requested profiles.
	for _, profileId := range profileIds {
		if _, err := r.client.Photo().Tag(ctx, photoResponse.ID, profileId); err != nil {
			resp.Diagnostics.AddError("Error tagging photo", err.Error())
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
