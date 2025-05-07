package profile

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
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
			Id:          "123",
			FirstName:   ptr("John"),
			LastName:    ptr("Doe"),
			MiddleName:  ptr("A"),
			MaidenName:  ptr("Smith"),
			DisplayName: ptr("John A Doe"),
			Gender:      ptr("male"),
			AboutMe:     ptr("This is a test profile"),
			Public:      true,
			IsAlive:     true,
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
		Expect(actualValue.Gender.ValueString()).To(Equal(*givenProfile.Gender))
		Expect(actualValue.About.ValueString()).To(Equal(*givenProfile.AboutMe))
		Expect(actualValue.Public.ValueBool()).To(Equal(givenProfile.Public))
		Expect(actualValue.Alive.ValueBool()).To(Equal(givenProfile.IsAlive))
		var actualNames = make(map[string]NameModel)
		Expect(actualValue.Names.ElementsAs(t.Context(), &actualNames, false).HasError()).To(BeFalse())
		Expect(actualNames).To(HaveKeyWithValue("en", NameModel{
			FistName:      types.StringPointerValue(givenProfile.Names["en"].FirstName),
			MiddleName:    types.StringPointerValue(givenProfile.Names["en"].MiddleName),
			LastName:      types.StringPointerValue(givenProfile.Names["en"].LastName),
			BirthLastName: types.StringPointerValue(givenProfile.Names["en"].MaidenName),
			DisplayName:   types.StringPointerValue(givenProfile.Names["en"].DisplayName),
			Nicknames:     types.SetNull(types.StringType),
		}))
		Expect(actualValue.CreatedAt.ValueString()).To(Equal(givenProfile.CreatedAt))
	})
}

func TestNameValueFrom(t *testing.T) {
	t.Run("Happy path, when a fully defined object is passed", func(t *testing.T) {
		RegisterTestingT(t)
		givenNames := map[string]geni.NameElement{
			"en": {
				FirstName:   ptr("John"),
				MiddleName:  ptr("A"),
				LastName:    ptr("Doe"),
				MaidenName:  ptr("Smith"),
				DisplayName: ptr("John A Doe"),
				Nicknames:   ptr("A,B"),
			},
			"fr": {
				FirstName:   ptr("Jean"),
				MiddleName:  ptr("B"),
				LastName:    ptr("Dupont"),
				MaidenName:  ptr("Bernard"),
				DisplayName: ptr("Jean B Bernard"),
				Nicknames:   ptr("C,D"),
			},
		}

		expectedNames := map[string]NameModel{
			"en": {
				FistName:      types.StringPointerValue(ptr("John")),
				MiddleName:    types.StringPointerValue(ptr("A")),
				LastName:      types.StringPointerValue(ptr("Doe")),
				BirthLastName: types.StringPointerValue(ptr("Smith")),
				DisplayName:   types.StringPointerValue(ptr("John A Doe")),
				Nicknames:     types.SetValueMust(types.StringType, []attr.Value{types.StringValue("A"), types.StringValue("B")}),
			},
			"fr": {
				FistName:      types.StringPointerValue(ptr("Jean")),
				MiddleName:    types.StringPointerValue(ptr("B")),
				LastName:      types.StringPointerValue(ptr("Dupont")),
				BirthLastName: types.StringPointerValue(ptr("Bernard")),
				DisplayName:   types.StringPointerValue(ptr("Jean B Bernard")),
				Nicknames:     types.SetValueMust(types.StringType, []attr.Value{types.StringValue("C"), types.StringValue("D")}),
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
