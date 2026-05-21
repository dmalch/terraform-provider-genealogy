package photo

import (
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/dmalch/terraform-provider-genealogy/internal/resource/event"
)

type ResourceModel struct {
	ID          types.String `tfsdk:"id"`
	Title       types.String `tfsdk:"title"`
	Description types.String `tfsdk:"description"`
	Date        types.String `tfsdk:"date"`
	File        types.String `tfsdk:"file"`
	FileName    types.String `tfsdk:"file_name"`
	Album       types.String `tfsdk:"album"`
	Profiles    types.Set    `tfsdk:"profiles"`
	Location    types.Object `tfsdk:"location"`
	Guid        types.String `tfsdk:"guid"`
	ContentType types.String `tfsdk:"content_type"`
	Attribution types.String `tfsdk:"attribution"`
	URL         types.String `tfsdk:"url"`
	Sizes       types.Map    `tfsdk:"sizes"`
	CreatedAt   types.String `tfsdk:"created_at"`
	UpdatedAt   types.String `tfsdk:"updated_at"`
}

type ResourceIdentityModel struct {
	ID types.String `tfsdk:"id"`
}

// NewEmptyResourceModel returns a ResourceModel whose collection fields are
// initialized to typed null values matching the schema. Use it when building a
// model from scratch with no prior state to seed from — e.g. when assembling a
// list-resource query result — so no collection attribute carries a type-less
// zero value that would later fail framework state writes.
func NewEmptyResourceModel() ResourceModel {
	return ResourceModel{
		Profiles: types.SetNull(types.StringType),
		Location: types.ObjectNull(event.LocationModelAttributeTypes()),
		Sizes:    types.MapNull(types.StringType),
	}
}
