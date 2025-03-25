package profile

import (
	"testing"

	"github.com/dmalch/terraform-provider-genealogy/internal/geni"

	. "github.com/onsi/gomega"
)

func ptr[T any](s T) *T {
	return &s
}

func TestValueFrom(t *testing.T) {
	t.Run("Happy path, when a fully defined object is passed", func(t *testing.T) {
		RegisterTestingT(t)
		givenProfile := &geni.ProfileResponse{
			Id:        "123",
			FirstName: "John",
			LastName:  "Doe",
			Gender:    "male",
			Names: map[string]geni.NameElement{
				"en": {
					FirstName:  ptr("John"),
					MiddleName: ptr("A"),
					LastName:   ptr("Doe"),
				},
			},
			Unions: []string{"union1", "union2"},
			Birth: &geni.EventElement{
				Date: &geni.DateElement{
					Day:   ptr[int32](1),
					Month: ptr[int32](1),
					Year:  ptr[int32](2000),
				},
			},
			Baptism: &geni.EventElement{
				Date: &geni.DateElement{
					Day:   ptr[int32](1),
					Month: ptr[int32](1),
					Year:  ptr[int32](2000),
				},
			},
			Death: &geni.EventElement{
				Date: &geni.DateElement{
					Day:   ptr[int32](1),
					Month: ptr[int32](1),
					Year:  ptr[int32](2000),
				},
			},
			Burial: &geni.EventElement{
				Date: &geni.DateElement{
					Day:   ptr[int32](1),
					Month: ptr[int32](1),
					Year:  ptr[int32](2000),
				},
			},
			CreatedAt: "1719709400",
		}

		actualValue := &ResourceModel{}
		diags := ValueFrom(t.Context(), givenProfile, actualValue)

		Expect(diags.HasError()).To(BeFalse())
		Expect(actualValue.ID.ValueString()).To(Equal(givenProfile.Id))
		Expect(actualValue.FirstName.ValueString()).To(Equal(givenProfile.FirstName))
		Expect(actualValue.LastName.ValueString()).To(Equal(givenProfile.LastName))
		Expect(actualValue.Gender.ValueString()).To(Equal(givenProfile.Gender))
		Expect(actualValue.CreatedAt.ValueString()).To(Equal(givenProfile.CreatedAt))
	})
}
