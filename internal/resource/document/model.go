package document

import (
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/dmalch/terraform-provider-genealogy/internal/resource/event"
)

// NewEmptyResourceModel returns a ResourceModel whose collection fields are
// initialized to typed null values matching the schema. Use it when building
// a model from scratch with no prior state to seed from — e.g. when assembling
// a list-resource query result. ValueFrom can then overwrite the fields the
// API does populate without leaving any collection attribute carrying a
// type-less zero value that would later fail framework state writes.
func NewEmptyResourceModel() ResourceModel {
	return ResourceModel{
		Date:     types.ObjectNull(event.DateModelAttributeTypes()),
		Location: types.ObjectNull(event.LocationModelAttributeTypes()),
		Profiles: types.SetNull(types.StringType),
		Projects: types.SetNull(types.StringType),
		Labels:   types.SetNull(types.StringType),
	}
}

type ResourceModel struct {
	ID          types.String `tfsdk:"id"`
	Title       types.String `tfsdk:"title"`
	Description types.String `tfsdk:"description"`
	ContentType types.String `tfsdk:"content_type"`
	Text        types.String `tfsdk:"text"`
	File        types.String `tfsdk:"file"`
	FileName    types.String `tfsdk:"file_name"`
	SourceUrl   types.String `tfsdk:"source_url"`
	Date        types.Object `tfsdk:"date"`
	Location    types.Object `tfsdk:"location"`
	Profiles    types.Set    `tfsdk:"profiles"`
	Projects    types.Set    `tfsdk:"projects"`
	Labels      types.Set    `tfsdk:"labels"`
	CreatedAt   types.String `tfsdk:"created_at"`
}

type ResourceIdentityModel struct {
	ID types.String `tfsdk:"id"`
}
