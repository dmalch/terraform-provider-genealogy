package event

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Model struct {
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Date        types.Object `tfsdk:"date"`
	Location    types.Object `tfsdk:"location"`
}

func (m Model) AttributeTypes() map[string]attr.Type {
	return EventModelAttributeTypes()
}

func EventModelAttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name":        types.StringType,
		"description": types.StringType,
		"date": types.ObjectType{
			AttrTypes: DateRangeModelAttributeTypes(),
		},
		"location": types.ObjectType{
			AttrTypes: LocationModelAttributeTypes(),
		},
	}
}

type DateModel struct {
	Circa types.Bool  `tfsdk:"circa"`
	Day   types.Int32 `tfsdk:"day"`
	Month types.Int32 `tfsdk:"month"`
	Year  types.Int32 `tfsdk:"year"`
}

func (m DateModel) AttributeTypes() map[string]attr.Type {
	return DateModelAttributeTypes()
}

func DateModelAttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"circa": types.BoolType,
		"day":   types.Int32Type,
		"month": types.Int32Type,
		"year":  types.Int32Type,
	}
}

type DateRangeModel struct {
	DateModel
	Range    types.String `tfsdk:"range"`
	EndCirca types.Bool   `tfsdk:"end_circa"`
	EndDay   types.Int32  `tfsdk:"end_day"`
	EndMonth types.Int32  `tfsdk:"end_month"`
	EndYear  types.Int32  `tfsdk:"end_year"`
}

func (m DateRangeModel) AttributeTypes() map[string]attr.Type {
	return DateRangeModelAttributeTypes()
}

func DateRangeModelAttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"range":     types.StringType,
		"circa":     types.BoolType,
		"day":       types.Int32Type,
		"month":     types.Int32Type,
		"year":      types.Int32Type,
		"end_circa": types.BoolType,
		"end_day":   types.Int32Type,
		"end_month": types.Int32Type,
		"end_year":  types.Int32Type,
	}
}

type LocationModel struct {
	City           types.String  `tfsdk:"city"`
	Country        types.String  `tfsdk:"country"`
	County         types.String  `tfsdk:"county"`
	Latitude       types.Float64 `tfsdk:"latitude"`
	Longitude      types.Float64 `tfsdk:"longitude"`
	PlaceName      types.String  `tfsdk:"place_name"`
	State          types.String  `tfsdk:"state"`
	StreetAddress1 types.String  `tfsdk:"street_address1"`
	StreetAddress2 types.String  `tfsdk:"street_address2"`
	StreetAddress3 types.String  `tfsdk:"street_address3"`
}

func (m LocationModel) AttributeTypes() map[string]attr.Type {
	return LocationModelAttributeTypes()
}

func LocationModelAttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"city":            types.StringType,
		"country":         types.StringType,
		"county":          types.StringType,
		"latitude":        types.Float64Type,
		"longitude":       types.Float64Type,
		"place_name":      types.StringType,
		"state":           types.StringType,
		"street_address1": types.StringType,
		"street_address2": types.StringType,
		"street_address3": types.StringType,
	}
}
