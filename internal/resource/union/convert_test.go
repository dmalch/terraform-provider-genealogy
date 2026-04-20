package union

import (
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
	t.Run("Happy path, when a fully defined union response is passed", func(t *testing.T) {
		RegisterTestingT(t)
		givenResponse := &geni.UnionResponse{
			Id:       "union-123",
			Partners: []string{"profile-1", "profile-2"},
			Children: []string{"profile-3", "profile-4"},
			Marriage: &geni.EventElement{
				Name: "Marriage",
				Date: &geni.DateElement{
					Day:   ptr[int32](15),
					Month: ptr[int32](6),
					Year:  ptr[int32](1990),
				},
				Location: &geni.LocationElement{
					City:    ptr("Paris"),
					Country: ptr("France"),
				},
			},
			Divorce: &geni.EventElement{
				Name: "Divorce",
				Date: &geni.DateElement{
					Day:   ptr[int32](1),
					Month: ptr[int32](3),
					Year:  ptr[int32](2000),
				},
			},
		}

		model := &ResourceModel{
			Children: types.SetNull(types.StringType),
			Partners: types.SetNull(types.StringType),
		}
		diags := ValueFrom(t.Context(), givenResponse, model)

		Expect(diags.HasError()).To(BeFalse())
		Expect(model.ID.ValueString()).To(Equal("union-123"))

		Expect(model.Partners.Elements()).To(HaveLen(2))
		Expect(model.Children.Elements()).To(HaveLen(2))

		Expect(model.Marriage.IsNull()).To(BeFalse())
		Expect(model.Marriage.IsUnknown()).To(BeFalse())
		Expect(model.Divorce.IsNull()).To(BeFalse())
		Expect(model.Divorce.IsUnknown()).To(BeFalse())
	})

	t.Run("When partners and children are empty", func(t *testing.T) {
		RegisterTestingT(t)
		givenResponse := &geni.UnionResponse{
			Id: "union-456",
		}

		model := &ResourceModel{
			Children: types.SetNull(types.StringType),
			Partners: types.SetNull(types.StringType),
		}
		diags := ValueFrom(t.Context(), givenResponse, model)

		Expect(diags.HasError()).To(BeFalse())
		Expect(model.ID.ValueString()).To(Equal("union-456"))
		// When no partners/children, the original null sets should remain
		Expect(model.Partners.IsNull()).To(BeTrue())
		Expect(model.Children.IsNull()).To(BeTrue())
	})

	t.Run("When marriage and divorce are nil", func(t *testing.T) {
		RegisterTestingT(t)
		givenResponse := &geni.UnionResponse{
			Id:       "union-789",
			Partners: []string{"profile-1"},
		}

		model := &ResourceModel{
			Children: types.SetNull(types.StringType),
			Partners: types.SetNull(types.StringType),
		}
		diags := ValueFrom(t.Context(), givenResponse, model)

		Expect(diags.HasError()).To(BeFalse())
		Expect(model.Marriage.IsNull()).To(BeTrue())
		Expect(model.Divorce.IsNull()).To(BeTrue())
	})

	t.Run("Splits API children into biological, foster, and adopted buckets", func(t *testing.T) {
		RegisterTestingT(t)
		givenResponse := &geni.UnionResponse{
			Id:              "union-split",
			Partners:        []string{"profile-1", "profile-2"},
			Children:        []string{"profile-3", "profile-4", "profile-5", "profile-6"},
			FosterChildren:  []string{"profile-4"},
			AdoptedChildren: []string{"profile-5"},
		}

		model := &ResourceModel{
			Children:        types.SetNull(types.StringType),
			FosterChildren:  types.SetNull(types.StringType),
			AdoptedChildren: types.SetNull(types.StringType),
			Partners:        types.SetNull(types.StringType),
		}
		diags := ValueFrom(t.Context(), givenResponse, model)

		Expect(diags.HasError()).To(BeFalse())

		bio, d := convertToSlice(t.Context(), model.Children)
		Expect(d.HasError()).To(BeFalse())
		Expect(bio).To(ConsistOf("profile-3", "profile-6"))

		foster, d := convertToSlice(t.Context(), model.FosterChildren)
		Expect(d.HasError()).To(BeFalse())
		Expect(foster).To(ConsistOf("profile-4"))

		adopted, d := convertToSlice(t.Context(), model.AdoptedChildren)
		Expect(d.HasError()).To(BeFalse())
		Expect(adopted).To(ConsistOf("profile-5"))
	})

	t.Run("Leaves foster and adopted null when the response has no subsets", func(t *testing.T) {
		RegisterTestingT(t)
		givenResponse := &geni.UnionResponse{
			Id:       "union-bio-only",
			Children: []string{"profile-3"},
		}

		model := &ResourceModel{
			Children:        types.SetNull(types.StringType),
			FosterChildren:  types.SetNull(types.StringType),
			AdoptedChildren: types.SetNull(types.StringType),
			Partners:        types.SetNull(types.StringType),
		}
		diags := ValueFrom(t.Context(), givenResponse, model)

		Expect(diags.HasError()).To(BeFalse())
		Expect(model.Children.Elements()).To(HaveLen(1))
		Expect(model.FosterChildren.IsNull()).To(BeTrue())
		Expect(model.AdoptedChildren.IsNull()).To(BeTrue())
	})
}

func TestUpdateComputedFields(t *testing.T) {
	t.Run("Updates computed fields with events", func(t *testing.T) {
		RegisterTestingT(t)
		givenResponse := &geni.UnionResponse{
			Id: "union-123",
			Marriage: &geni.EventElement{
				Name:        "Marriage",
				Description: ptr("A beautiful ceremony"),
				Date: &geni.DateElement{
					Year: ptr[int32](1990),
				},
				Location: &geni.LocationElement{
					City: ptr("Paris"),
				},
			},
			Divorce: &geni.EventElement{
				Name: "Divorce",
				Date: &geni.DateElement{
					Year: ptr[int32](2000),
				},
			},
		}

		// Create a model with existing event objects that have unknown computed fields
		marriageObj, _ := event.ValueFrom(t.Context(), &geni.EventElement{
			Date: &geni.DateElement{
				Year: ptr[int32](1990),
			},
		})
		divorceObj, _ := event.ValueFrom(t.Context(), &geni.EventElement{
			Date: &geni.DateElement{
				Year: ptr[int32](2000),
			},
		})

		model := &ResourceModel{
			Marriage: marriageObj,
			Divorce:  divorceObj,
		}

		diags := UpdateComputedFields(t.Context(), givenResponse, model)

		Expect(diags.HasError()).To(BeFalse())
		Expect(model.ID.ValueString()).To(Equal("union-123"))
		Expect(model.Marriage.IsNull()).To(BeFalse())
		Expect(model.Divorce.IsNull()).To(BeFalse())
	})

	t.Run("Updates computed fields without events", func(t *testing.T) {
		RegisterTestingT(t)
		givenResponse := &geni.UnionResponse{
			Id: "union-456",
		}

		model := &ResourceModel{
			Marriage: types.ObjectNull(event.EventModelAttributeTypes()),
			Divorce:  types.ObjectNull(event.EventModelAttributeTypes()),
		}

		diags := UpdateComputedFields(t.Context(), givenResponse, model)

		Expect(diags.HasError()).To(BeFalse())
		Expect(model.ID.ValueString()).To(Equal("union-456"))
	})
}

func TestRequestFrom(t *testing.T) {
	t.Run("Happy path, with marriage and divorce events", func(t *testing.T) {
		RegisterTestingT(t)
		marriageObj := types.ObjectValueMust(event.EventModelAttributeTypes(),
			map[string]attr.Value{
				"name":        types.StringValue("Marriage"),
				"description": types.StringNull(),
				"date": types.ObjectValueMust(event.DateRangeModelAttributeTypes(),
					map[string]attr.Value{
						"range":     types.StringNull(),
						"circa":     types.BoolNull(),
						"day":       types.Int32Value(15),
						"month":     types.Int32Value(6),
						"year":      types.Int32Value(1990),
						"end_circa": types.BoolNull(),
						"end_day":   types.Int32Null(),
						"end_month": types.Int32Null(),
						"end_year":  types.Int32Null(),
					}),
				"location": types.ObjectNull(event.LocationModelAttributeTypes()),
			})

		plan := ResourceModel{
			Marriage: marriageObj,
			Divorce:  types.ObjectNull(event.EventModelAttributeTypes()),
		}

		request, diags := RequestFrom(t.Context(), plan)

		Expect(diags.HasError()).To(BeFalse())
		Expect(request).ToNot(BeNil())
		Expect(request.Marriage).ToNot(BeNil())
		Expect(request.Marriage.Date).ToNot(BeNil())
		Expect(request.Marriage.Date.Day).To(HaveValue(Equal(int32(15))))
		Expect(request.Marriage.Date.Month).To(HaveValue(Equal(int32(6))))
		Expect(request.Marriage.Date.Year).To(HaveValue(Equal(int32(1990))))
		Expect(request.Divorce).To(BeNil())
	})

	t.Run("When both events are null", func(t *testing.T) {
		RegisterTestingT(t)
		plan := ResourceModel{
			Marriage: types.ObjectNull(event.EventModelAttributeTypes()),
			Divorce:  types.ObjectNull(event.EventModelAttributeTypes()),
		}

		request, diags := RequestFrom(t.Context(), plan)

		Expect(diags.HasError()).To(BeFalse())
		Expect(request).ToNot(BeNil())
		Expect(request.Marriage).To(BeNil())
		Expect(request.Divorce).To(BeNil())
	})
}

func TestHashMapFrom(t *testing.T) {
	t.Run("Creates a hash map from a slice", func(t *testing.T) {
		RegisterTestingT(t)
		result := hashMapFrom([]string{"a", "b", "c"})

		Expect(result).To(HaveLen(3))
		Expect(result).To(HaveKey("a"))
		Expect(result).To(HaveKey("b"))
		Expect(result).To(HaveKey("c"))
	})

	t.Run("Deduplicates entries", func(t *testing.T) {
		RegisterTestingT(t)
		result := hashMapFrom([]string{"a", "a", "b"})

		Expect(result).To(HaveLen(2))
	})

	t.Run("Returns empty map for empty slice", func(t *testing.T) {
		RegisterTestingT(t)
		result := hashMapFrom([]string{})

		Expect(result).To(BeEmpty())
	})
}

func TestModifierFor(t *testing.T) {
	RegisterTestingT(t)
	plan := ResourceModel{
		Children:        types.SetValueMust(types.StringType, []attr.Value{types.StringValue("bio-1")}),
		FosterChildren:  types.SetValueMust(types.StringType, []attr.Value{types.StringValue("foster-1")}),
		AdoptedChildren: types.SetValueMust(types.StringType, []attr.Value{types.StringValue("adopted-1")}),
	}

	foster, err := convertToSlice(t.Context(), plan.FosterChildren)
	Expect(err.HasError()).To(BeFalse())
	adopted, err := convertToSlice(t.Context(), plan.AdoptedChildren)
	Expect(err.HasError()).To(BeFalse())

	fosterSet := hashMapFrom(foster)
	adoptedSet := hashMapFrom(adopted)

	Expect(modifierFor("foster-1", fosterSet, adoptedSet)).To(Equal("foster"))
	Expect(modifierFor("adopted-1", fosterSet, adoptedSet)).To(Equal("adopt"))
	Expect(modifierFor("bio-1", fosterSet, adoptedSet)).To(Equal(""))
	Expect(modifierFor("unknown", fosterSet, adoptedSet)).To(Equal(""))
}

func TestConvertToSlice(t *testing.T) {
	t.Run("Converts a set to a slice", func(t *testing.T) {
		RegisterTestingT(t)
		set := types.SetValueMust(types.StringType, []attr.Value{
			types.StringValue("a"),
			types.StringValue("b"),
			types.StringValue("c"),
		})

		result, diags := convertToSlice(t.Context(), set)

		Expect(diags.HasError()).To(BeFalse())
		Expect(result).To(HaveLen(3))
		Expect(result).To(ContainElements("a", "b", "c"))
	})

	t.Run("Returns empty slice for empty set", func(t *testing.T) {
		RegisterTestingT(t)
		set := types.SetValueMust(types.StringType, []attr.Value{})

		result, diags := convertToSlice(t.Context(), set)

		Expect(diags.HasError()).To(BeFalse())
		Expect(result).To(BeEmpty())
	})
}
