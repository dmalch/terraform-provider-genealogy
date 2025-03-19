package event

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Model struct {
	Description types.String `json:"description"`
	Name        types.String `json:"name"`
	Date        types.Object `json:"date"`
	Location    types.Object `json:"location"`
}

func (m Model) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"description": types.StringType,
		"name":        types.StringType,
		"date": types.ObjectType{
			AttrTypes: dateModelObjectType(),
		},
		"location": types.ObjectType{
			AttrTypes: locationModelObjectType(),
		},
	}
}

type DateModel struct {
	Range    types.String `json:"range"`
	Circa    types.Bool   `json:"circa"`
	Day      types.Number `json:"day"`
	Month    types.Number `json:"month"`
	Year     types.Number `json:"year"`
	EndCirca types.Bool   `json:"end_circa"`
	EndDay   types.Number `json:"end_day"`
	EndMonth types.Number `json:"end_month"`
	EndYear  types.Number `json:"end_year"`
}

func (m DateModel) AttributeTypes() map[string]attr.Type {
	return dateModelObjectType()
}

func dateModelObjectType() map[string]attr.Type {
	return map[string]attr.Type{
		"range":     types.StringType,
		"circa":     types.BoolType,
		"day":       types.NumberType,
		"month":     types.NumberType,
		"year":      types.NumberType,
		"end_circa": types.BoolType,
		"end_day":   types.NumberType,
		"end_month": types.NumberType,
		"end_year":  types.NumberType,
	}
}

type LocationModel struct {
	City           types.String `json:"city"`
	Country        types.String `json:"country"`
	County         types.String `json:"county"`
	Latitude       types.Number `json:"latitude"`
	Longitude      types.Number `json:"longitude"`
	PlaceName      types.String `json:"place_name"`
	State          types.String `json:"state"`
	StreetAddress1 types.String `json:"street_address1"`
	StreetAddress2 types.String `json:"street_address2"`
	StreetAddress3 types.String `json:"street_address3"`
}

func (m LocationModel) AttributeTypes() map[string]attr.Type {
	return locationModelObjectType()
}

func locationModelObjectType() map[string]attr.Type {
	return map[string]attr.Type{
		"city":            types.StringType,
		"country":         types.StringType,
		"county":          types.StringType,
		"latitude":        types.NumberType,
		"longitude":       types.NumberType,
		"place_name":      types.StringType,
		"state":           types.StringType,
		"street_address1": types.StringType,
		"street_address2": types.StringType,
		"street_address3": types.StringType,
	}
}
