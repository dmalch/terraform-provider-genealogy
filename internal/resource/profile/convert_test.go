package profile

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	. "github.com/onsi/gomega"

	"github.com/dmalch/terraform-provider-genealogy/internal/geni"
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
		var actualNames = make(map[string]NameModel)
		Expect(actualValue.Names.ElementsAs(t.Context(), &actualNames, false).HasError()).To(BeFalse())
		Expect(actualNames).To(HaveKeyWithValue("en", NameModel{
			FistName:   types.StringPointerValue(givenProfile.Names["en"].FirstName),
			MiddleName: types.StringPointerValue(givenProfile.Names["en"].MiddleName),
			LastName:   types.StringPointerValue(givenProfile.Names["en"].LastName),
		}))
		Expect(actualValue.CreatedAt.ValueString()).To(Equal(givenProfile.CreatedAt))
	})
}

func TestNameValueFrom(t *testing.T) {
	t.Run("Happy path, when a fully defined object is passed", func(t *testing.T) {
		RegisterTestingT(t)
		givenNames := map[string]geni.NameElement{
			"en": {
				FirstName:  ptr("John"),
				MiddleName: ptr("A"),
				LastName:   ptr("Doe"),
			},
			"fr": {
				FirstName:  ptr("Jean"),
				MiddleName: ptr("B"),
				LastName:   ptr("Dupont"),
			},
		}

		expectedNames := map[string]NameModel{
			"en": {
				FistName:   types.StringPointerValue(ptr("John")),
				MiddleName: types.StringPointerValue(ptr("A")),
				LastName:   types.StringPointerValue(ptr("Doe")),
			},
			"fr": {
				FistName:   types.StringPointerValue(ptr("Jean")),
				MiddleName: types.StringPointerValue(ptr("B")),
				LastName:   types.StringPointerValue(ptr("Dupont")),
			},
		}

		actualValue, diags := NameValueFrom(t.Context(), givenNames)

		Expect(diags.HasError()).To(BeFalse())
		Expect(actualValue.Elements()).To(HaveLen(len(expectedNames)))

		var actualNames map[string]NameModel
		Expect(actualValue.ElementsAs(t.Context(), &actualNames, false).HasError()).To(BeFalse())
		Expect(actualNames).To(Equal(expectedNames))
	})
}
