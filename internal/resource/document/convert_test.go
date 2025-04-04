package document

import (
	"math/big"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	. "github.com/onsi/gomega"

	"github.com/dmalch/terraform-provider-genealogy/internal/geni"
	"github.com/dmalch/terraform-provider-genealogy/internal/resource/event"
)

func ptr[T any](s T) *T {
	return &s
}

func TestValueFrom(t *testing.T) {
	t.Run("Happy path, when a fully defined object is passed", func(t *testing.T) {
		RegisterTestingT(t)
		givenResponse := geni.DocumentResponse{
			Id:          "123",
			Title:       "Test Document",
			Description: "This is a test document",
			ContentType: "application/pdf",
			Date: &geni.DateElement{
				Day:   ptr[int32](1),
				Month: ptr[int32](1),
				Year:  ptr[int32](2000),
			},
			Location: &geni.LocationElement{
				City:           ptr("City"),
				Country:        ptr("Country"),
				County:         ptr("County"),
				Latitude:       big.NewFloat(1.0),
				Longitude:      big.NewFloat(2.0),
				PlaceName:      ptr("Place Name"),
				State:          ptr("State"),
				StreetAddress1: ptr("Street Address 1"),
				StreetAddress2: ptr("Street Address 2"),
				StreetAddress3: ptr("Street Address 3"),
			},
			Tags:      []string{"tag1", "tag2"},
			Labels:    []string{"label1", "label2"},
			UpdatedAt: "12345679",
			CreatedAt: "12345678",
		}

		convertedModel := ResourceModel{}
		diags := ValueFrom(t.Context(), &givenResponse, &convertedModel)

		Expect(diags).To(BeEmpty())
		Expect(convertedModel.ID.ValueString()).To(Equal(givenResponse.Id))
		Expect(convertedModel.Title.ValueString()).To(Equal(givenResponse.Title))
		Expect(convertedModel.Description.ValueString()).To(Equal(givenResponse.Description))
		Expect(convertedModel.ContentType.ValueString()).To(Equal(givenResponse.ContentType))
		Expect(convertedModel.Date.IsNull()).ToNot(BeTrue())
		Expect(convertedModel.Date.IsUnknown()).ToNot(BeTrue())
		Expect(convertedModel.Location.IsNull()).ToNot(BeTrue())
		Expect(convertedModel.Location.IsUnknown()).ToNot(BeTrue())
		Expect(convertedModel.Profiles.IsNull()).ToNot(BeTrue())
		Expect(convertedModel.Profiles.IsUnknown()).ToNot(BeTrue())
		Expect(convertedModel.Profiles.Elements()).To(HaveLen(2))
		Expect(convertedModel.Labels).ToNot(BeNil())
		Expect(convertedModel.Labels.Elements()).To(HaveLen(2))
		Expect(convertedModel.CreatedAt.ValueString()).To(Equal(givenResponse.CreatedAt))
	})
}

func TestRequestFrom(t *testing.T) {
	t.Run("Happy path, when a fully defined object is passed", func(t *testing.T) {
		RegisterTestingT(t)
		givenResourceModel := ResourceModel{
			Title:       types.StringValue("Test Document"),
			Description: types.StringValue("This is a test document"),
			ContentType: types.StringValue("application/pdf"),
			Date: types.ObjectValueMust(event.DateModel{}.AttributeTypes(),
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
			Location: types.ObjectValueMust(event.LocationModel{}.AttributeTypes(),
				map[string]attr.Value{
					"city":            types.StringValue("City"),
					"country":         types.StringValue("Country"),
					"county":          types.StringValue("County"),
					"latitude":        types.NumberValue(big.NewFloat(1.0)),
					"longitude":       types.NumberValue(big.NewFloat(2.0)),
					"place_name":      types.StringValue("Place Name"),
					"state":           types.StringValue("State"),
					"street_address1": types.StringValue("Street Address 1"),
					"street_address2": types.StringValue("Street Address 2"),
					"street_address3": types.StringValue("Street Address 3"),
				}),
			Labels: types.SetValueMust(types.StringType, []attr.Value{types.StringValue("label1"), types.StringValue("label2")}),
		}

		documentRequest, diags := RequestFrom(t.Context(), givenResourceModel)

		Expect(diags).To(BeEmpty())
		Expect(documentRequest).ToNot(BeNil())
		Expect(documentRequest.Title).To(Equal(givenResourceModel.Title.ValueString()))
		Expect(documentRequest.Description).To(HaveValue(Equal(givenResourceModel.Description.ValueString())))
		Expect(documentRequest.ContentType).To(HaveValue(Equal(givenResourceModel.ContentType.ValueString())))
		Expect(documentRequest.Date).ToNot(BeNil())
		Expect(documentRequest.Date.Range).To(HaveValue(Equal("between")))
		Expect(documentRequest.Date.Circa).To(HaveValue(BeTrue()))
		Expect(documentRequest.Date.Day).To(HaveValue(Equal(int32(19))))
		Expect(documentRequest.Date.Month).To(HaveValue(Equal(int32(8))))
		Expect(documentRequest.Date.Year).To(HaveValue(Equal(int32(1922))))
		Expect(documentRequest.Date.EndCirca).To(HaveValue(BeFalse()))
		Expect(documentRequest.Date.EndDay).To(HaveValue(Equal(int32(20))))
		Expect(documentRequest.Date.EndMonth).To(HaveValue(Equal(int32(8))))
		Expect(documentRequest.Date.EndYear).To(HaveValue(Equal(int32(1922))))
		Expect(documentRequest.Location).ToNot(BeNil())
		Expect(documentRequest.Location.City).To(HaveValue(Equal("City")))
		Expect(documentRequest.Location.Country).To(HaveValue(Equal("Country")))
		Expect(documentRequest.Location.County).To(HaveValue(Equal("County")))
		Expect(documentRequest.Location.Latitude).To(Equal(big.NewFloat(1.0)))
		Expect(documentRequest.Location.Longitude).To(Equal(big.NewFloat(2.0)))
		Expect(documentRequest.Location.PlaceName).To(HaveValue(Equal("Place Name")))
		Expect(documentRequest.Location.State).To(HaveValue(Equal("State")))
		Expect(documentRequest.Location.StreetAddress1).To(HaveValue(Equal("Street Address 1")))
		Expect(documentRequest.Location.StreetAddress2).To(HaveValue(Equal("Street Address 2")))
		Expect(documentRequest.Location.StreetAddress3).To(HaveValue(Equal("Street Address 3")))
		Expect(documentRequest.Labels).To(HaveValue(Equal("label1,label2")))
	})
}
