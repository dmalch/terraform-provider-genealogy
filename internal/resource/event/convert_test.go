package event

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	. "github.com/onsi/gomega"

	"github.com/dmalch/terraform-provider-genealogy/internal/geni"
)

func ptr[T any](s T) *T {
	return &s
}

func TestElementFrom(t *testing.T) {
	t.Run("regular case", func(t *testing.T) {
		RegisterTestingT(t)
		eventObject := types.ObjectValueMust(Model{}.AttributeTypes(),
			map[string]attr.Value{
				"name":        types.StringValue("Event Name"),
				"description": types.StringValue("Event Description"),
				"date": types.ObjectValueMust(DateRangeModel{}.AttributeTypes(),
					map[string]attr.Value{
						"range":     types.StringValue("between"),
						"circa":     types.BoolValue(true),
						"day":       types.Int32Value(19),
						"month":     types.Int32Value(8),
						"year":      types.Int32Value(1922),
						"end_circa": types.BoolValue(false),
						"end_day":   types.Int32Value(20),
						"end_month": types.Int32Value(8),
						"end_year":  types.Int32Value(1922),
					}),
				"location": types.ObjectValueMust(LocationModel{}.AttributeTypes(),
					map[string]attr.Value{
						"city":            types.StringValue("City"),
						"country":         types.StringValue("Country"),
						"county":          types.StringValue("County"),
						"latitude":        types.Float64Value(1.0),
						"longitude":       types.Float64Value(2.0),
						"place_name":      types.StringValue("Place Name"),
						"state":           types.StringValue("State"),
						"street_address1": types.StringValue("Street Address 1"),
						"street_address2": types.StringValue("Street Address 2"),
						"street_address3": types.StringValue("Street Address 3"),
					}),
			},
		)

		element, diags := ElementFrom(t.Context(), eventObject)

		Expect(diags).To(BeEmpty())
		Expect(element.Name).To(Equal("Event Name"))
		Expect(element.Description).To(HaveValue(Equal("Event Description")))
		Expect(element.Date).ToNot(BeNil())
		Expect(element.Date.Range).To(HaveValue(Equal("between")))
		Expect(element.Date.Circa).To(HaveValue(BeTrue()))
		Expect(element.Date.Day).To(HaveValue(Equal(int32(19))))
		Expect(element.Date.Month).To(HaveValue(Equal(int32(8))))
		Expect(element.Date.Year).To(HaveValue(Equal(int32(1922))))
		Expect(element.Date.EndCirca).To(HaveValue(BeFalse()))
		Expect(element.Date.EndDay).To(HaveValue(Equal(int32(20))))
		Expect(element.Date.EndMonth).To(HaveValue(Equal(int32(8))))
		Expect(element.Date.EndYear).To(HaveValue(Equal(int32(1922))))
		Expect(element.Location).ToNot(BeNil())
		Expect(element.Location.City).To(HaveValue(Equal("City")))
		Expect(element.Location.Country).To(HaveValue(Equal("Country")))
		Expect(element.Location.County).To(HaveValue(Equal("County")))
		Expect(element.Location.Latitude).To(HaveValue(Equal(1.0)))
		Expect(element.Location.Longitude).To(HaveValue(Equal(2.0)))
		Expect(element.Location.PlaceName).To(HaveValue(Equal("Place Name")))
		Expect(element.Location.State).To(HaveValue(Equal("State")))
		Expect(element.Location.StreetAddress1).To(HaveValue(Equal("Street Address 1")))
		Expect(element.Location.StreetAddress2).To(HaveValue(Equal("Street Address 2")))
		Expect(element.Location.StreetAddress3).To(HaveValue(Equal("Street Address 3")))
	})

	t.Run("when date and location are nulls", func(t *testing.T) {
		RegisterTestingT(t)
		eventObject := types.ObjectValueMust(EventModelAttributeTypes(),
			map[string]attr.Value{
				"name":        types.StringValue("Event Name"),
				"description": types.StringValue("Event Description"),
				"date":        types.ObjectNull(DateRangeModel{}.AttributeTypes()),
				"location":    types.ObjectNull(LocationModel{}.AttributeTypes()),
			})

		element, diags := ElementFrom(t.Context(), eventObject)

		Expect(diags).To(BeEmpty())
		Expect(element.Name).To(Equal("Event Name"))
		Expect(element.Description).To(HaveValue(Equal("Event Description")))
		Expect(element.Date).To(BeNil())
		Expect(element.Location).To(BeNil())
	})

	t.Run("when event is null", func(t *testing.T) {
		RegisterTestingT(t)
		eventObject := types.ObjectNull(Model{}.AttributeTypes())

		element, diags := ElementFrom(t.Context(), eventObject)

		Expect(diags).To(BeEmpty())
		Expect(element).To(BeNil())
	})
}

func TestValueFrom(t *testing.T) {
	t.Run("full event with date and location", func(t *testing.T) {
		RegisterTestingT(t)
		eventElement := &geni.EventElement{
			Name:        "Birth",
			Description: ptr("Born in city"),
			Date: &geni.DateElement{
				Circa: ptr(false),
				Day:   ptr(int32(15)),
				Month: ptr(int32(3)),
				Year:  ptr(int32(1990)),
			},
			Location: &geni.LocationElement{
				City:    ptr("Springfield"),
				Country: ptr("US"),
			},
		}

		result, diags := ValueFrom(t.Context(), eventElement)

		Expect(diags).To(BeEmpty())
		Expect(result.IsNull()).To(BeFalse())

		var model Model
		diags = result.As(t.Context(), &model, basetypes.ObjectAsOptions{})
		Expect(diags).To(BeEmpty())
		Expect(model.Name.ValueString()).To(Equal("Birth"))
		Expect(model.Description.ValueString()).To(Equal("Born in city"))
		Expect(model.Date.IsNull()).To(BeFalse())
		Expect(model.Location.IsNull()).To(BeFalse())
	})

	t.Run("nil event returns null object", func(t *testing.T) {
		RegisterTestingT(t)

		result, diags := ValueFrom(t.Context(), nil)

		Expect(diags).To(BeEmpty())
		Expect(result.IsNull()).To(BeTrue())
	})
}

func TestDateValueFrom(t *testing.T) {
	t.Run("full date", func(t *testing.T) {
		RegisterTestingT(t)
		dateElement := &geni.DateElement{
			Circa: ptr(true),
			Day:   ptr(int32(25)),
			Month: ptr(int32(12)),
			Year:  ptr(int32(1900)),
		}

		result, diags := DateValueFrom(t.Context(), dateElement)

		Expect(diags).To(BeEmpty())
		Expect(result.IsNull()).To(BeFalse())

		var model DateModel
		diags = result.As(t.Context(), &model, basetypes.ObjectAsOptions{})
		Expect(diags).To(BeEmpty())
		Expect(model.Circa.ValueBool()).To(BeTrue())
		Expect(model.Day.ValueInt32()).To(Equal(int32(25)))
		Expect(model.Month.ValueInt32()).To(Equal(int32(12)))
		Expect(model.Year.ValueInt32()).To(Equal(int32(1900)))
	})

	t.Run("nil returns null", func(t *testing.T) {
		RegisterTestingT(t)

		result, diags := DateValueFrom(t.Context(), nil)

		Expect(diags).To(BeEmpty())
		Expect(result.IsNull()).To(BeTrue())
	})
}

func TestDateRangeValueFrom(t *testing.T) {
	t.Run("full range date", func(t *testing.T) {
		RegisterTestingT(t)
		dateElement := &geni.DateElement{
			Range:    ptr("between"),
			Circa:    ptr(true),
			Day:      ptr(int32(1)),
			Month:    ptr(int32(1)),
			Year:     ptr(int32(1800)),
			EndCirca: ptr(false),
			EndDay:   ptr(int32(31)),
			EndMonth: ptr(int32(12)),
			EndYear:  ptr(int32(1850)),
		}

		result, diags := DateRangeValueFrom(t.Context(), dateElement)

		Expect(diags).To(BeEmpty())
		Expect(result.IsNull()).To(BeFalse())

		var model DateRangeModel
		diags = result.As(t.Context(), &model, basetypes.ObjectAsOptions{})
		Expect(diags).To(BeEmpty())
		Expect(model.Range.ValueString()).To(Equal("between"))
		Expect(model.Circa.ValueBool()).To(BeTrue())
		Expect(model.Day.ValueInt32()).To(Equal(int32(1)))
		Expect(model.Month.ValueInt32()).To(Equal(int32(1)))
		Expect(model.Year.ValueInt32()).To(Equal(int32(1800)))
		Expect(model.EndCirca.ValueBool()).To(BeFalse())
		Expect(model.EndDay.ValueInt32()).To(Equal(int32(31)))
		Expect(model.EndMonth.ValueInt32()).To(Equal(int32(12)))
		Expect(model.EndYear.ValueInt32()).To(Equal(int32(1850)))
	})

	t.Run("nil returns null", func(t *testing.T) {
		RegisterTestingT(t)

		result, diags := DateRangeValueFrom(t.Context(), nil)

		Expect(diags).To(BeEmpty())
		Expect(result.IsNull()).To(BeTrue())
	})
}

func TestLocationValueFrom(t *testing.T) {
	t.Run("full location", func(t *testing.T) {
		RegisterTestingT(t)
		locationElement := &geni.LocationElement{
			City:           ptr("London"),
			Country:        ptr("UK"),
			County:         ptr("Greater London"),
			Latitude:       ptr(51.5074),
			Longitude:      ptr(-0.1278),
			PlaceName:      ptr("Westminster"),
			State:          ptr("England"),
			StreetAddress1: ptr("10 Downing St"),
			StreetAddress2: ptr("Suite 1"),
			StreetAddress3: ptr("Floor 2"),
		}

		result, diags := LocationValueFrom(t.Context(), locationElement)

		Expect(diags).To(BeEmpty())
		Expect(result.IsNull()).To(BeFalse())

		var model LocationModel
		diags = result.As(t.Context(), &model, basetypes.ObjectAsOptions{})
		Expect(diags).To(BeEmpty())
		Expect(model.City.ValueString()).To(Equal("London"))
		Expect(model.Country.ValueString()).To(Equal("UK"))
		Expect(model.County.ValueString()).To(Equal("Greater London"))
		Expect(model.Latitude.ValueFloat64()).To(Equal(51.5074))
		Expect(model.Longitude.ValueFloat64()).To(Equal(-0.1278))
		Expect(model.PlaceName.ValueString()).To(Equal("Westminster"))
		Expect(model.State.ValueString()).To(Equal("England"))
		Expect(model.StreetAddress1.ValueString()).To(Equal("10 Downing St"))
		Expect(model.StreetAddress2.ValueString()).To(Equal("Suite 1"))
		Expect(model.StreetAddress3.ValueString()).To(Equal("Floor 2"))
	})

	t.Run("nil returns null", func(t *testing.T) {
		RegisterTestingT(t)

		result, diags := LocationValueFrom(t.Context(), nil)

		Expect(diags).To(BeEmpty())
		Expect(result.IsNull()).To(BeTrue())
	})
}

func TestDateElementFrom(t *testing.T) {
	t.Run("full model", func(t *testing.T) {
		RegisterTestingT(t)
		model := &DateModel{
			Circa: types.BoolValue(true),
			Day:   types.Int32Value(5),
			Month: types.Int32Value(6),
			Year:  types.Int32Value(1950),
		}

		result := DateElementFrom(model)

		Expect(result).ToNot(BeNil())
		Expect(result.Circa).To(HaveValue(BeTrue()))
		Expect(result.Day).To(HaveValue(Equal(int32(5))))
		Expect(result.Month).To(HaveValue(Equal(int32(6))))
		Expect(result.Year).To(HaveValue(Equal(int32(1950))))
	})

	t.Run("nil returns nil", func(t *testing.T) {
		RegisterTestingT(t)

		result := DateElementFrom(nil)

		Expect(result).To(BeNil())
	})
}

func TestDateRangeElementFrom(t *testing.T) {
	t.Run("full model with range", func(t *testing.T) {
		RegisterTestingT(t)
		model := &DateRangeModel{
			DateModel: DateModel{
				Circa: types.BoolValue(false),
				Day:   types.Int32Value(1),
				Month: types.Int32Value(1),
				Year:  types.Int32Value(1800),
			},
			Range:    types.StringValue("between"),
			EndCirca: types.BoolValue(true),
			EndDay:   types.Int32Value(31),
			EndMonth: types.Int32Value(12),
			EndYear:  types.Int32Value(1899),
		}

		result := DateRangeElementFrom(model)

		Expect(result).ToNot(BeNil())
		Expect(result.Range).To(HaveValue(Equal("between")))
		Expect(result.Circa).To(HaveValue(BeFalse()))
		Expect(result.Day).To(HaveValue(Equal(int32(1))))
		Expect(result.Month).To(HaveValue(Equal(int32(1))))
		Expect(result.Year).To(HaveValue(Equal(int32(1800))))
		Expect(result.EndCirca).To(HaveValue(BeTrue()))
		Expect(result.EndDay).To(HaveValue(Equal(int32(31))))
		Expect(result.EndMonth).To(HaveValue(Equal(int32(12))))
		Expect(result.EndYear).To(HaveValue(Equal(int32(1899))))
	})

	t.Run("nil returns nil", func(t *testing.T) {
		RegisterTestingT(t)

		result := DateRangeElementFrom(nil)

		Expect(result).To(BeNil())
	})
}

func TestLocationElementFrom(t *testing.T) {
	t.Run("full model", func(t *testing.T) {
		RegisterTestingT(t)
		model := &LocationModel{
			City:           types.StringValue("Paris"),
			Country:        types.StringValue("France"),
			County:         types.StringValue("Île-de-France"),
			Latitude:       types.Float64Value(48.8566),
			Longitude:      types.Float64Value(2.3522),
			PlaceName:      types.StringValue("Eiffel Tower"),
			State:          types.StringValue("Île-de-France"),
			StreetAddress1: types.StringValue("5 Avenue Anatole"),
			StreetAddress2: types.StringValue("Apt 1"),
			StreetAddress3: types.StringValue("Building A"),
		}

		result := LocationElementFrom(model)

		Expect(result).ToNot(BeNil())
		Expect(result.City).To(HaveValue(Equal("Paris")))
		Expect(result.Country).To(HaveValue(Equal("France")))
		Expect(result.County).To(HaveValue(Equal("Île-de-France")))
		Expect(result.Latitude).To(HaveValue(Equal(48.8566)))
		Expect(result.Longitude).To(HaveValue(Equal(2.3522)))
		Expect(result.PlaceName).To(HaveValue(Equal("Eiffel Tower")))
		Expect(result.State).To(HaveValue(Equal("Île-de-France")))
		Expect(result.StreetAddress1).To(HaveValue(Equal("5 Avenue Anatole")))
		Expect(result.StreetAddress2).To(HaveValue(Equal("Apt 1")))
		Expect(result.StreetAddress3).To(HaveValue(Equal("Building A")))
	})

	t.Run("nil returns nil", func(t *testing.T) {
		RegisterTestingT(t)

		result := LocationElementFrom(nil)

		Expect(result).To(BeNil())
	})
}

func TestUpdateComputedFieldsInEvent(t *testing.T) {
	t.Run("unknown name and description get populated from response", func(t *testing.T) {
		RegisterTestingT(t)
		eventObject := types.ObjectValueMust(EventModelAttributeTypes(),
			map[string]attr.Value{
				"name":        types.StringUnknown(),
				"description": types.StringUnknown(),
				"date":        types.ObjectNull(DateRangeModel{}.AttributeTypes()),
				"location":    types.ObjectNull(LocationModel{}.AttributeTypes()),
			})
		eventElement := &geni.EventElement{
			Name:        "Marriage",
			Description: ptr("Wedding ceremony"),
		}

		result, diags := UpdateComputedFieldsInEvent(t.Context(), eventObject, eventElement)

		Expect(diags).To(BeEmpty())
		Expect(result.IsNull()).To(BeFalse())

		var model Model
		diags = result.As(t.Context(), &model, basetypes.ObjectAsOptions{})
		Expect(diags).To(BeEmpty())
		Expect(model.Name.ValueString()).To(Equal("Marriage"))
		Expect(model.Description.ValueString()).To(Equal("Wedding ceremony"))
	})

	t.Run("null event object returns null", func(t *testing.T) {
		RegisterTestingT(t)
		eventObject := types.ObjectNull(EventModelAttributeTypes())

		result, diags := UpdateComputedFieldsInEvent(t.Context(), eventObject, nil)

		Expect(diags).To(BeEmpty())
		Expect(result.IsNull()).To(BeTrue())
	})

	t.Run("nil event element nulls computed fields", func(t *testing.T) {
		RegisterTestingT(t)
		eventObject := types.ObjectValueMust(EventModelAttributeTypes(),
			map[string]attr.Value{
				"name":        types.StringUnknown(),
				"description": types.StringUnknown(),
				"date":        types.ObjectNull(DateRangeModel{}.AttributeTypes()),
				"location":    types.ObjectNull(LocationModel{}.AttributeTypes()),
			})

		result, diags := UpdateComputedFieldsInEvent(t.Context(), eventObject, nil)

		Expect(diags).To(BeEmpty())
		Expect(result.IsNull()).To(BeFalse())

		var model Model
		diags = result.As(t.Context(), &model, basetypes.ObjectAsOptions{})
		Expect(diags).To(BeEmpty())
		Expect(model.Name.IsNull()).To(BeTrue())
		Expect(model.Description.IsNull()).To(BeTrue())
	})
}

func TestUpdateComputedFieldsInLocationObject(t *testing.T) {
	t.Run("regular case when all elements are defined", func(t *testing.T) {
		RegisterTestingT(t)

		givenLocationObject := types.ObjectValueMust(LocationModelAttributeTypes(),
			map[string]attr.Value{
				"city":            types.StringValue("City"),
				"country":         types.StringValue("Country"),
				"county":          types.StringValue("County"),
				"latitude":        types.Float64Value(1.0),
				"longitude":       types.Float64Value(2.0),
				"place_name":      types.StringValue("Place Name"),
				"state":           types.StringValue("State"),
				"street_address1": types.StringValue("Street Address 1"),
				"street_address2": types.StringValue("Street Address 2"),
				"street_address3": types.StringValue("Street Address 3"),
			})
		givenLocationResponse := &geni.LocationElement{
			City:           ptr("City Response"),
			Country:        ptr("Country Response"),
			County:         ptr("County Response"),
			Latitude:       ptr(1.1),
			Longitude:      ptr(2.1),
			PlaceName:      ptr("Place Name Response"),
			State:          ptr("State Response"),
			StreetAddress1: ptr("Street Address 1 Response"),
			StreetAddress2: ptr("Street Address 2 Response"),
			StreetAddress3: ptr("Street Address 3 Response"),
		}
		updatedLocationObject, diags := UpdateComputedFieldsInLocationObject(t.Context(), givenLocationObject, givenLocationResponse)
		Expect(diags).To(BeEmpty())
		Expect(updatedLocationObject).ToNot(BeNil())
		Expect(updatedLocationObject).To(Equal(givenLocationObject))
	})

	t.Run("when latitude and longitude is unknown", func(t *testing.T) {
		RegisterTestingT(t)

		givenLocationObject := types.ObjectValueMust(LocationModelAttributeTypes(),
			map[string]attr.Value{
				"city":            types.StringNull(),
				"country":         types.StringNull(),
				"county":          types.StringNull(),
				"latitude":        types.Float64Unknown(),
				"longitude":       types.Float64Unknown(),
				"place_name":      types.StringNull(),
				"state":           types.StringNull(),
				"street_address1": types.StringNull(),
				"street_address2": types.StringNull(),
				"street_address3": types.StringNull(),
			})
		givenLocationResponse := &geni.LocationElement{
			City:           ptr("City Response"),
			Country:        ptr("Country Response"),
			County:         ptr("County Response"),
			Latitude:       ptr(1.1),
			Longitude:      ptr(2.1),
			PlaceName:      ptr("Place Name Response"),
			State:          ptr("State Response"),
			StreetAddress1: ptr("Street Address 1 Response"),
			StreetAddress2: ptr("Street Address 2 Response"),
			StreetAddress3: ptr("Street Address 3 Response"),
		}
		updatedLocationObject, diags := UpdateComputedFieldsInLocationObject(t.Context(), givenLocationObject, givenLocationResponse)
		Expect(diags).To(BeEmpty())
		Expect(updatedLocationObject).ToNot(BeNil())

		var updatedLocationModel LocationModel
		diags = updatedLocationObject.As(t.Context(), &updatedLocationModel, basetypes.ObjectAsOptions{})
		Expect(diags).To(BeEmpty())
		Expect(updatedLocationModel.Latitude.ValueFloat64()).To(Equal(*givenLocationResponse.Latitude))
		Expect(updatedLocationModel.Longitude.ValueFloat64()).To(Equal(*givenLocationResponse.Longitude))
	})

	t.Run("when latitude and longitude is unknown and response is a zero", func(t *testing.T) {
		RegisterTestingT(t)

		givenLocationObject := types.ObjectValueMust(LocationModelAttributeTypes(),
			map[string]attr.Value{
				"city":            types.StringNull(),
				"country":         types.StringNull(),
				"county":          types.StringNull(),
				"latitude":        types.Float64Unknown(),
				"longitude":       types.Float64Unknown(),
				"place_name":      types.StringNull(),
				"state":           types.StringNull(),
				"street_address1": types.StringNull(),
				"street_address2": types.StringNull(),
				"street_address3": types.StringNull(),
			})
		givenLocationResponse := &geni.LocationElement{
			City:           ptr("City Response"),
			Country:        ptr("Country Response"),
			County:         ptr("County Response"),
			Latitude:       ptr(0.0),
			Longitude:      ptr(0.0),
			PlaceName:      ptr("Place Name Response"),
			State:          ptr("State Response"),
			StreetAddress1: ptr("Street Address 1 Response"),
			StreetAddress2: ptr("Street Address 2 Response"),
			StreetAddress3: ptr("Street Address 3 Response"),
		}
		updatedLocationObject, diags := UpdateComputedFieldsInLocationObject(t.Context(), givenLocationObject, givenLocationResponse)
		Expect(diags).To(BeEmpty())
		Expect(updatedLocationObject).ToNot(BeNil())

		var updatedLocationModel LocationModel
		diags = updatedLocationObject.As(t.Context(), &updatedLocationModel, basetypes.ObjectAsOptions{})
		Expect(diags).To(BeEmpty())
		Expect(updatedLocationModel.Latitude.IsNull()).To(BeTrue())
		Expect(updatedLocationModel.Longitude.IsNull()).To(BeTrue())
	})

	t.Run("when latitude and longitude is unknown and response is a null", func(t *testing.T) {
		RegisterTestingT(t)

		givenLocationObject := types.ObjectValueMust(LocationModelAttributeTypes(),
			map[string]attr.Value{
				"city":            types.StringNull(),
				"country":         types.StringNull(),
				"county":          types.StringNull(),
				"latitude":        types.Float64Unknown(),
				"longitude":       types.Float64Unknown(),
				"place_name":      types.StringNull(),
				"state":           types.StringNull(),
				"street_address1": types.StringNull(),
				"street_address2": types.StringNull(),
				"street_address3": types.StringNull(),
			})
		updatedLocationObject, diags := UpdateComputedFieldsInLocationObject(t.Context(), givenLocationObject, nil)
		Expect(diags).To(BeEmpty())
		Expect(updatedLocationObject).ToNot(BeNil())

		var updatedLocationModel LocationModel
		diags = updatedLocationObject.As(t.Context(), &updatedLocationModel, basetypes.ObjectAsOptions{})
		Expect(diags).To(BeEmpty())
		Expect(updatedLocationModel.Latitude.IsNull()).To(BeTrue())
		Expect(updatedLocationModel.Longitude.IsNull()).To(BeTrue())
	})
}
