package profile

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	. "github.com/onsi/gomega"

	"github.com/dmalch/go-geni"
	"github.com/dmalch/terraform-provider-genealogy/internal/resource/event"
)

func ptr[T any](s T) *T {
	return &s
}

func TestValueFrom(t *testing.T) {
	t.Run("Happy path, when a fully defined object is passed", func(t *testing.T) {
		RegisterTestingT(t)
		givenProfile := &geni.ProfileResponse{
			Id:          "123",
			Guid:        "abcdef0123456789abcdef0123456789",
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
			Unions:     []string{"union1", "union2"},
			Projects:   []string{"project-4505748", "project-4497781"},
			Title:      "Dr.",
			Occupation: "Astronaut",
			Suffix:     "Jr.",
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
		Expect(actualValue.Guid.ValueString()).To(Equal(givenProfile.Guid))
		Expect(actualValue.Title.ValueString()).To(Equal(givenProfile.Title))
		Expect(actualValue.Occupation.ValueString()).To(Equal(givenProfile.Occupation))
		Expect(actualValue.Suffix.ValueString()).To(Equal(givenProfile.Suffix))
		Expect(actualValue.Gender.ValueString()).To(Equal(*givenProfile.Gender))
		var actualAbout = make(map[string]string)
		Expect(actualValue.About.ElementsAs(t.Context(), &actualAbout, false).HasError()).To(BeFalse())
		Expect(actualAbout).To(HaveKeyWithValue("en-US", *givenProfile.AboutMe))
		Expect(actualValue.Public.ValueBool()).To(Equal(givenProfile.Public))
		Expect(actualValue.Alive.ValueBool()).To(Equal(givenProfile.IsAlive))
		var actualNames = make(map[string]NameModel)
		Expect(actualValue.Names.ElementsAs(t.Context(), &actualNames, false).HasError()).To(BeFalse())
		Expect(actualNames).To(HaveKeyWithValue("en", NameModel{
			FirstName:     types.StringPointerValue(givenProfile.Names["en"].FirstName),
			MiddleName:    types.StringPointerValue(givenProfile.Names["en"].MiddleName),
			LastName:      types.StringPointerValue(givenProfile.Names["en"].LastName),
			BirthLastName: types.StringPointerValue(givenProfile.Names["en"].MaidenName),
			DisplayName:   types.StringPointerValue(givenProfile.Names["en"].DisplayName),
			Nicknames:     types.SetNull(types.StringType),
		}))
		Expect(actualValue.CreatedAt.ValueString()).To(Equal(givenProfile.CreatedAt))

		var actualProjects []string
		Expect(actualValue.Projects.ElementsAs(t.Context(), &actualProjects, false).HasError()).To(BeFalse())
		Expect(actualProjects).To(ConsistOf("project-4505748", "project-4497781"))
	})

	t.Run("when the response has no guid, model.Guid is a null string", func(t *testing.T) {
		RegisterTestingT(t)
		given := &geni.ProfileResponse{Id: "profile-1"}
		actual := &ResourceModel{}

		Expect(ValueFrom(t.Context(), given, actual).HasError()).To(BeFalse())
		Expect(actual.Guid.IsNull()).To(BeTrue())
	})

	t.Run("when the response has no project memberships, Projects is a typed null Set[String]", func(t *testing.T) {
		RegisterTestingT(t)
		given := &geni.ProfileResponse{Id: "profile-1"}
		actual := &ResourceModel{}

		Expect(ValueFrom(t.Context(), given, actual).HasError()).To(BeFalse())
		Expect(actual.Projects.IsNull()).To(BeTrue())
		Expect(actual.Projects.ElementType(t.Context())).To(Equal(types.StringType))
	})

	t.Run("when about me in multiple languages is defiled", func(t *testing.T) {
		RegisterTestingT(t)
		givenProfile := &geni.ProfileResponse{
			DetailStrings: map[string]geni.DetailsString{
				"en-US": {
					AboutMe: ptr("This is a test profile in English"),
				},
				"fr-FR": {
					AboutMe: ptr("Ceci est un profil de test en français"),
				},
			},
		}

		actualValue := &ResourceModel{}
		diags := ValueFrom(t.Context(), givenProfile, actualValue)

		Expect(diags.HasError()).To(BeFalse())
		var actualAbout = make(map[string]string)
		Expect(actualValue.About.ElementsAs(t.Context(), &actualAbout, false).HasError()).To(BeFalse())
		Expect(actualAbout).To(HaveKeyWithValue("en-US", *givenProfile.DetailStrings["en-US"].AboutMe))
		Expect(actualAbout).To(HaveKeyWithValue("fr-FR", *givenProfile.DetailStrings["fr-FR"].AboutMe))
	})
}

func TestRequestFrom(t *testing.T) {
	t.Run("Happy path, with names, events, about, gender, alive and public flags", func(t *testing.T) {
		RegisterTestingT(t)
		givenModel := ResourceModel{
			Gender: types.StringValue("male"),
			Alive:  types.BoolValue(true),
			Public: types.BoolValue(true),
			Names: types.MapValueMust(types.ObjectType{AttrTypes: NameAttributeTypes()}, map[string]attr.Value{
				"en": types.ObjectValueMust(NameAttributeTypes(), map[string]attr.Value{
					"first_name":      types.StringValue("John"),
					"middle_name":     types.StringNull(),
					"last_name":       types.StringValue("Doe"),
					"birth_last_name": types.StringNull(),
					"display_name":    types.StringNull(),
					"nicknames":       types.SetNull(types.StringType),
				}),
			}),
			About: types.MapValueMust(types.StringType, map[string]attr.Value{
				"en-US": types.StringValue("A test profile"),
			}),
			Birth: types.ObjectValueMust(event.EventModelAttributeTypes(),
				map[string]attr.Value{
					"name":        types.StringNull(),
					"description": types.StringNull(),
					"date": types.ObjectValueMust(event.DateRangeModelAttributeTypes(),
						map[string]attr.Value{
							"range": types.StringNull(), "circa": types.BoolNull(),
							"day": types.Int32Value(1), "month": types.Int32Value(1), "year": types.Int32Value(2000),
							"end_circa": types.BoolNull(), "end_day": types.Int32Null(), "end_month": types.Int32Null(), "end_year": types.Int32Null(),
						}),
					"location": types.ObjectNull(event.LocationModelAttributeTypes()),
				}),
			Baptism:          types.ObjectNull(event.EventModelAttributeTypes()),
			Death:            types.ObjectNull(event.EventModelAttributeTypes()),
			Burial:           types.ObjectNull(event.EventModelAttributeTypes()),
			CauseOfDeath:     types.StringNull(),
			CurrentResidence: types.ObjectNull(event.LocationModelAttributeTypes()),
			Title:            types.StringValue("Dr."),
			Occupation:       types.StringValue("Astronaut"),
			Suffix:           types.StringValue("Jr."),
		}

		request, diags := RequestFrom(t.Context(), givenModel)

		Expect(diags.HasError()).To(BeFalse())
		Expect(request).ToNot(BeNil())
		Expect(request.Gender).To(HaveValue(Equal("male")))
		Expect(request.IsAlive).To(BeTrue())
		Expect(request.Public).To(BeTrue())
		Expect(request.Title).To(Equal("Dr."))
		Expect(request.Occupation).To(Equal("Astronaut"))
		Expect(request.Suffix).To(Equal("Jr."))
		Expect(request.Names).To(HaveLen(1))
		Expect(request.Names["en"].FirstName).To(HaveValue(Equal("John")))
		Expect(request.Names["en"].LastName).To(HaveValue(Equal("Doe")))
		Expect(request.DetailStrings).To(HaveLen(1))
		Expect(request.DetailStrings["en-US"].AboutMe).To(HaveValue(Equal("A test profile")))
		Expect(request.Birth).ToNot(BeNil())
		Expect(request.Birth.Date).ToNot(BeNil())
		Expect(request.Birth.Date.Day).To(HaveValue(Equal(int32(1))))
		// Null events should become empty EventElements (not nil)
		Expect(request.Baptism).ToNot(BeNil())
		Expect(request.Death).ToNot(BeNil())
		Expect(request.Burial).ToNot(BeNil())
	})

	t.Run("Minimal model with null events and empty names", func(t *testing.T) {
		RegisterTestingT(t)
		givenModel := ResourceModel{
			Gender:           types.StringNull(),
			Alive:            types.BoolValue(false),
			Public:           types.BoolValue(false),
			Names:            types.MapNull(types.ObjectType{AttrTypes: NameAttributeTypes()}),
			About:            types.MapNull(types.StringType),
			Birth:            types.ObjectNull(event.EventModelAttributeTypes()),
			Baptism:          types.ObjectNull(event.EventModelAttributeTypes()),
			Death:            types.ObjectNull(event.EventModelAttributeTypes()),
			Burial:           types.ObjectNull(event.EventModelAttributeTypes()),
			CauseOfDeath:     types.StringNull(),
			CurrentResidence: types.ObjectNull(event.LocationModelAttributeTypes()),
		}

		request, diags := RequestFrom(t.Context(), givenModel)

		Expect(diags.HasError()).To(BeFalse())
		Expect(request).ToNot(BeNil())
		Expect(request.Gender).To(BeNil())
		Expect(request.Names).To(BeEmpty())
		Expect(request.DetailStrings).To(BeEmpty())
	})
}

func TestUpdateComputedFields(t *testing.T) {
	t.Run("Updates ID, unions, events, about, deleted, merged_into, and created_at", func(t *testing.T) {
		RegisterTestingT(t)
		givenProfile := &geni.ProfileResponse{
			Id:     "profile-123",
			Unions: []string{"union-1", "union-2"},
			// Projects intentionally omitted: Create/Update pass the pre-link
			// API response, so it is stale and must not overwrite plan.Projects.
			Projects: []string{"project-stale"},
			Birth: &geni.EventElement{
				Name: "Birth",
				Date: &geni.DateElement{Year: ptr[int32](1990)},
			},
			DetailStrings: map[string]geni.DetailsString{
				"en-US": {AboutMe: ptr("Updated about")},
			},
			Deleted:    true,
			MergedInto: "profile-456",
			CreatedAt:  "1719709400",
		}

		planProjects, planProjectsDiags := types.SetValueFrom(t.Context(), types.StringType, []string{"project-from-plan"})
		Expect(planProjectsDiags.HasError()).To(BeFalse())

		model := &ResourceModel{
			Birth:            types.ObjectNull(event.EventModelAttributeTypes()),
			Baptism:          types.ObjectNull(event.EventModelAttributeTypes()),
			Death:            types.ObjectNull(event.EventModelAttributeTypes()),
			Burial:           types.ObjectNull(event.EventModelAttributeTypes()),
			CurrentResidence: types.ObjectNull(event.LocationModelAttributeTypes()),
			About:            types.MapNull(types.StringType),
			Projects:         planProjects,
		}

		diags := UpdateComputedFields(t.Context(), givenProfile, model)

		Expect(diags.HasError()).To(BeFalse())
		Expect(model.ID.ValueString()).To(Equal("profile-123"))
		Expect(model.Unions.Elements()).To(HaveLen(2))
		// Projects in the plan must survive: the API response passed to
		// UpdateComputedFields is captured before AddProfileToProject runs,
		// so its Projects field is stale and must not overwrite the plan.
		var actualProjects []string
		Expect(model.Projects.ElementsAs(t.Context(), &actualProjects, false).HasError()).To(BeFalse())
		Expect(actualProjects).To(ConsistOf("project-from-plan"))
		Expect(model.Deleted.ValueBool()).To(BeTrue())
		Expect(model.MergedInto.ValueString()).To(Equal("profile-456"))
		Expect(model.CreatedAt.ValueString()).To(Equal("1719709400"))

		var actualAbout = make(map[string]string)
		Expect(model.About.ElementsAs(t.Context(), &actualAbout, false).HasError()).To(BeFalse())
		Expect(actualAbout).To(HaveKeyWithValue("en-US", "Updated about"))
	})
}

func TestNameElementsFrom(t *testing.T) {
	t.Run("Converts Terraform names map to API NameElement map with nicknames", func(t *testing.T) {
		RegisterTestingT(t)
		namesMap := types.MapValueMust(types.ObjectType{AttrTypes: NameAttributeTypes()}, map[string]attr.Value{
			"en": types.ObjectValueMust(NameAttributeTypes(), map[string]attr.Value{
				"first_name":      types.StringValue("John"),
				"middle_name":     types.StringValue("A"),
				"last_name":       types.StringValue("Doe"),
				"birth_last_name": types.StringValue("Smith"),
				"display_name":    types.StringValue("John Doe"),
				"nicknames":       types.SetValueMust(types.StringType, []attr.Value{types.StringValue("Johnny"), types.StringValue("JD")}),
			}),
		})

		result, diags := NameElementsFrom(t.Context(), namesMap)

		Expect(diags.HasError()).To(BeFalse())
		Expect(result).To(HaveLen(1))
		Expect(result["en"].FirstName).To(HaveValue(Equal("John")))
		Expect(result["en"].MiddleName).To(HaveValue(Equal("A")))
		Expect(result["en"].LastName).To(HaveValue(Equal("Doe")))
		Expect(result["en"].MaidenName).To(HaveValue(Equal("Smith")))
		Expect(result["en"].DisplayName).To(HaveValue(Equal("John Doe")))
		// Nicknames should be joined as CSV
		Expect(result["en"].Nicknames).ToNot(BeNil())
		Expect(*result["en"].Nicknames).To(Or(Equal("Johnny,JD"), Equal("JD,Johnny")))
	})
}

func TestValueFromEdgeCases(t *testing.T) {
	t.Run("Empty profile with no names, no events", func(t *testing.T) {
		RegisterTestingT(t)
		givenProfile := &geni.ProfileResponse{
			Id: "profile-empty",
		}

		actualValue := &ResourceModel{}
		diags := ValueFrom(t.Context(), givenProfile, actualValue)

		Expect(diags.HasError()).To(BeFalse())
		Expect(actualValue.ID.ValueString()).To(Equal("profile-empty"))
		Expect(actualValue.Names.IsNull()).To(BeTrue())
		Expect(actualValue.Birth.IsNull()).To(BeTrue())
		Expect(actualValue.Death.IsNull()).To(BeTrue())
		Expect(actualValue.Baptism.IsNull()).To(BeTrue())
		Expect(actualValue.Burial.IsNull()).To(BeTrue())
	})

	t.Run("Profile with only DetailStrings", func(t *testing.T) {
		RegisterTestingT(t)
		givenProfile := &geni.ProfileResponse{
			Id: "profile-details",
			DetailStrings: map[string]geni.DetailsString{
				"en-US": {AboutMe: ptr("English bio")},
				"de-DE": {AboutMe: ptr("German bio")},
			},
		}

		actualValue := &ResourceModel{}
		diags := ValueFrom(t.Context(), givenProfile, actualValue)

		Expect(diags.HasError()).To(BeFalse())
		var actualAbout = make(map[string]string)
		Expect(actualValue.About.ElementsAs(t.Context(), &actualAbout, false).HasError()).To(BeFalse())
		Expect(actualAbout).To(HaveLen(2))
		Expect(actualAbout).To(HaveKeyWithValue("en-US", "English bio"))
		Expect(actualAbout).To(HaveKeyWithValue("de-DE", "German bio"))
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
				FirstName:     types.StringPointerValue(ptr("John")),
				MiddleName:    types.StringPointerValue(ptr("A")),
				LastName:      types.StringPointerValue(ptr("Doe")),
				BirthLastName: types.StringPointerValue(ptr("Smith")),
				DisplayName:   types.StringPointerValue(ptr("John A Doe")),
				Nicknames:     types.SetValueMust(types.StringType, []attr.Value{types.StringValue("A"), types.StringValue("B")}),
			},
			"fr": {
				FirstName:     types.StringPointerValue(ptr("Jean")),
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

func TestNamesWithFlatFallback(t *testing.T) {
	t.Run("Returns the API names map as-is when it has entries", func(t *testing.T) {
		RegisterTestingT(t)
		original := map[string]geni.NameElement{
			"fr": {FirstName: ptr("Jean"), LastName: ptr("Dupont")},
		}
		profile := &geni.ProfileResponse{
			Names:     original,
			FirstName: ptr("ignored-because-map-is-populated"),
		}

		result := namesWithFlatFallback(profile)

		Expect(result).To(Equal(original))
	})

	t.Run("Synthesizes an en-US entry from flat fields when the map is empty", func(t *testing.T) {
		RegisterTestingT(t)
		profile := &geni.ProfileResponse{
			FirstName:   ptr("John"),
			MiddleName:  ptr("A"),
			LastName:    ptr("Doe"),
			MaidenName:  ptr("Smith"),
			DisplayName: ptr("John A Doe"),
		}

		result := namesWithFlatFallback(profile)

		Expect(result).To(HaveLen(1))
		Expect(result).To(HaveKey("en-US"))
		Expect(result["en-US"].FirstName).To(HaveValue(Equal("John")))
		Expect(result["en-US"].MiddleName).To(HaveValue(Equal("A")))
		Expect(result["en-US"].LastName).To(HaveValue(Equal("Doe")))
		Expect(result["en-US"].MaidenName).To(HaveValue(Equal("Smith")))
		Expect(result["en-US"].DisplayName).To(HaveValue(Equal("John A Doe")))
		Expect(result["en-US"].Nicknames).To(BeNil())
	})

	t.Run("Joins flat nicknames into a comma-separated string", func(t *testing.T) {
		RegisterTestingT(t)
		profile := &geni.ProfileResponse{
			FirstName: ptr("John"),
			Nicknames: []string{"Johnny", "JD"},
		}

		result := namesWithFlatFallback(profile)

		Expect(result["en-US"].Nicknames).To(HaveValue(Equal("Johnny,JD")))
	})

	t.Run("Returns the original empty map when nothing populated", func(t *testing.T) {
		RegisterTestingT(t)
		profile := &geni.ProfileResponse{}

		result := namesWithFlatFallback(profile)

		Expect(result).To(BeEmpty())
	})

	t.Run("Returns the original empty map when only Nicknames slice is non-empty alongside nil flat fields", func(t *testing.T) {
		// Nicknames alone (without any other name field) is a legitimate signal
		// of presence — synthesize so the en-US entry exists for downstream
		// state writing.
		RegisterTestingT(t)
		profile := &geni.ProfileResponse{
			Nicknames: []string{"Jay"},
		}

		result := namesWithFlatFallback(profile)

		Expect(result).To(HaveLen(1))
		Expect(result["en-US"].Nicknames).To(HaveValue(Equal("Jay")))
		Expect(result["en-US"].FirstName).To(BeNil())
		Expect(result["en-US"].LastName).To(BeNil())
	})
}

func TestValueFrom_NamesFlatFallback(t *testing.T) {
	t.Run("ValueFrom hydrates names from flat fields when the API map is empty", func(t *testing.T) {
		RegisterTestingT(t)
		givenResponse := &geni.ProfileResponse{
			Id:        "profile-42",
			Public:    true,
			FirstName: ptr("John"),
			LastName:  ptr("Doe"),
		}

		var model ResourceModel
		diags := ValueFrom(t.Context(), givenResponse, &model)

		Expect(diags.HasError()).To(BeFalse())
		Expect(model.Names.IsNull()).To(BeFalse())

		var names map[string]NameModel
		Expect(model.Names.ElementsAs(t.Context(), &names, false).HasError()).To(BeFalse())
		Expect(names).To(HaveLen(1))
		Expect(names["en-US"].FirstName.ValueString()).To(Equal("John"))
		Expect(names["en-US"].LastName.ValueString()).To(Equal("Doe"))
	})

	t.Run("ValueFrom keeps API names when both the map and flat fields are populated", func(t *testing.T) {
		RegisterTestingT(t)
		givenResponse := &geni.ProfileResponse{
			Id:        "profile-43",
			Public:    true,
			FirstName: ptr("ignored"),
			LastName:  ptr("ignored"),
			Names: map[string]geni.NameElement{
				"fr": {FirstName: ptr("Jean"), LastName: ptr("Dupont")},
			},
		}

		var model ResourceModel
		diags := ValueFrom(t.Context(), givenResponse, &model)

		Expect(diags.HasError()).To(BeFalse())

		var names map[string]NameModel
		Expect(model.Names.ElementsAs(t.Context(), &names, false).HasError()).To(BeFalse())
		Expect(names).To(HaveLen(1))
		Expect(names).To(HaveKey("fr"))
		Expect(names).NotTo(HaveKey("en-US"))
		Expect(names["fr"].FirstName.ValueString()).To(Equal("Jean"))
	})

	t.Run("ValueFrom returns null names when neither the map nor flat fields are populated", func(t *testing.T) {
		RegisterTestingT(t)
		givenResponse := &geni.ProfileResponse{
			Id:     "profile-44",
			Public: true,
		}

		var model ResourceModel
		diags := ValueFrom(t.Context(), givenResponse, &model)

		Expect(diags.HasError()).To(BeFalse())
		Expect(model.Names.IsNull()).To(BeTrue())
	})
}
