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
		Expect(element.Description).To(Equal("Event Description"))
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
		Expect(element.Description).To(Equal("Event Description"))
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
}
