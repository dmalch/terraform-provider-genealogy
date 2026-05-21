package tfset

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	. "github.com/onsi/gomega"
)

func TestStrings(t *testing.T) {
	t.Run("Returns the elements of a populated set", func(t *testing.T) {
		RegisterTestingT(t)
		set := types.SetValueMust(types.StringType, []attr.Value{
			types.StringValue("a"),
			types.StringValue("b"),
			types.StringValue("c"),
		})

		result, diags := Strings(t.Context(), set)

		Expect(diags.HasError()).To(BeFalse())
		Expect(result).To(ConsistOf("a", "b", "c"))
	})

	t.Run("Returns a non-nil empty slice for an empty set", func(t *testing.T) {
		RegisterTestingT(t)
		set := types.SetValueMust(types.StringType, []attr.Value{})

		result, diags := Strings(t.Context(), set)

		Expect(diags.HasError()).To(BeFalse())
		Expect(result).ToNot(BeNil())
		Expect(result).To(BeEmpty())
	})

	t.Run("Returns a non-nil empty slice for a null set", func(t *testing.T) {
		RegisterTestingT(t)
		result, diags := Strings(t.Context(), types.SetNull(types.StringType))

		Expect(diags.HasError()).To(BeFalse())
		Expect(result).ToNot(BeNil())
		Expect(result).To(BeEmpty())
	})
}

func TestIndex(t *testing.T) {
	t.Run("Indexes every value for membership lookup", func(t *testing.T) {
		RegisterTestingT(t)
		result := Index([]string{"a", "b", "c"})

		Expect(result).To(HaveLen(3))
		Expect(result).To(HaveKey("a"))
		Expect(result).To(HaveKey("b"))
		Expect(result).To(HaveKey("c"))
	})

	t.Run("Collapses duplicate values", func(t *testing.T) {
		RegisterTestingT(t)
		result := Index([]string{"a", "a", "b"})

		Expect(result).To(HaveLen(2))
	})

	t.Run("Returns an empty index for an empty slice", func(t *testing.T) {
		RegisterTestingT(t)
		Expect(Index(nil)).To(BeEmpty())
		Expect(Index([]string{})).To(BeEmpty())
	})
}
