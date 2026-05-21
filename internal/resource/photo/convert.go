package photo

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	geniphoto "github.com/dmalch/go-geni/photo"
	"github.com/dmalch/terraform-provider-genealogy/internal/resource/event"
)

// ValueFrom populates model from a Geni photo response. File and FileName are
// left untouched — Geni never returns the uploaded bytes — so they survive from
// prior state.
func ValueFrom(ctx context.Context, response *geniphoto.Photo, model *ResourceModel) diag.Diagnostics {
	var d diag.Diagnostics

	model.ID = types.StringValue(response.ID)
	model.Title = types.StringValue(response.Title)
	model.Description = stringOrNull(response.Description)
	model.Date = stringOrNull(response.Date)
	model.Album = stringOrNull(response.AlbumId)
	model.Guid = stringOrNull(response.Guid)
	model.ContentType = stringOrNull(response.ContentType)
	model.Attribution = stringOrNull(response.Attribution)
	model.URL = stringOrNull(response.Url)
	model.CreatedAt = stringOrNull(response.CreatedAt)
	model.UpdatedAt = stringOrNull(response.UpdatedAt)

	profiles, diags := types.SetValueFrom(ctx, types.StringType, response.Tags)
	d.Append(diags...)
	model.Profiles = profiles

	sizes, diags := types.MapValueFrom(ctx, types.StringType, response.Sizes)
	d.Append(diags...)
	model.Sizes = sizes

	location, diags := event.LocationValueFrom(ctx, response.Location)
	d.Append(diags...)
	model.Location = location

	return d
}

// RequestFrom builds the Geni update request from model. The photo file is set
// only at creation (the schema marks it RequiresReplace), and location is
// read-only, so neither is sent.
func RequestFrom(model ResourceModel) *geniphoto.Request {
	return &geniphoto.Request{
		Title:       model.Title.ValueString(),
		Description: model.Description.ValueString(),
		Date:        model.Date.ValueString(),
	}
}

// UpdateComputedFields fills the computed attributes of model from a Geni photo
// response while preserving the values the caller already set from the plan.
func UpdateComputedFields(ctx context.Context, response *geniphoto.Photo, model *ResourceModel) diag.Diagnostics {
	var d diag.Diagnostics

	if model.ID.IsNull() || model.ID.IsUnknown() {
		model.ID = types.StringValue(response.ID)
	}
	if model.Album.IsNull() || model.Album.IsUnknown() {
		model.Album = stringOrNull(response.AlbumId)
	}
	if model.Profiles.IsNull() || model.Profiles.IsUnknown() {
		profiles, diags := types.SetValueFrom(ctx, types.StringType, response.Tags)
		d.Append(diags...)
		model.Profiles = profiles
	}

	model.Guid = stringOrNull(response.Guid)
	model.ContentType = stringOrNull(response.ContentType)
	model.Attribution = stringOrNull(response.Attribution)
	model.URL = stringOrNull(response.Url)
	model.CreatedAt = stringOrNull(response.CreatedAt)
	model.UpdatedAt = stringOrNull(response.UpdatedAt)

	sizes, diags := types.MapValueFrom(ctx, types.StringType, response.Sizes)
	d.Append(diags...)
	model.Sizes = sizes

	// Location is read-only — always taken from the API response.
	location, diags := event.LocationValueFrom(ctx, response.Location)
	d.Append(diags...)
	model.Location = location

	return d
}

// stringOrNull maps an empty Geni string field to a null value, so an absent
// optional attribute does not round-trip as an empty string and flap the plan.
func stringOrNull(s string) types.String {
	if s == "" {
		return types.StringNull()
	}
	return types.StringValue(s)
}
