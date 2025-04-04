package document

import (
	"math/big"
	"testing"

	. "github.com/onsi/gomega"

	"github.com/dmalch/terraform-provider-genealogy/internal/geni"
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
