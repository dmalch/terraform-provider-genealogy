package geniplanmodifier

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	. "github.com/onsi/gomega"
)

func TestContains(t *testing.T) {
	t.Run("Returns true when element is found", func(t *testing.T) {
		RegisterTestingT(t)
		elements := []attr.Value{
			types.StringValue("a"),
			types.StringValue("b"),
			types.StringValue("c"),
		}

		Expect(contains(elements, types.StringValue("b"))).To(BeTrue())
	})

	t.Run("Returns false when element is not found", func(t *testing.T) {
		RegisterTestingT(t)
		elements := []attr.Value{
			types.StringValue("a"),
			types.StringValue("b"),
		}

		Expect(contains(elements, types.StringValue("z"))).To(BeFalse())
	})

	t.Run("Returns false for empty slice", func(t *testing.T) {
		RegisterTestingT(t)
		Expect(contains([]attr.Value{}, types.StringValue("a"))).To(BeFalse())
	})
}

func TestValuesAreRemovedFromState(t *testing.T) {
	t.Run("Returns early when config is null", func(t *testing.T) {
		RegisterTestingT(t)
		resp := &setplanmodifier.RequiresReplaceIfFuncResponse{}
		req := planmodifier.SetRequest{
			ConfigValue: types.SetNull(types.StringType),
			StateValue:  types.SetValueMust(types.StringType, []attr.Value{types.StringValue("a")}),
			PlanValue:   types.SetValueMust(types.StringType, []attr.Value{types.StringValue("a")}),
		}

		ValuesAreRemovedFromState(t.Context(), req, resp)

		Expect(resp.RequiresReplace).To(BeFalse())
	})

	t.Run("Does not require replace when all state values are in plan", func(t *testing.T) {
		RegisterTestingT(t)
		resp := &setplanmodifier.RequiresReplaceIfFuncResponse{}
		req := planmodifier.SetRequest{
			ConfigValue: types.SetValueMust(types.StringType, []attr.Value{types.StringValue("a"), types.StringValue("b")}),
			StateValue:  types.SetValueMust(types.StringType, []attr.Value{types.StringValue("a"), types.StringValue("b")}),
			PlanValue:   types.SetValueMust(types.StringType, []attr.Value{types.StringValue("a"), types.StringValue("b")}),
		}

		ValuesAreRemovedFromState(t.Context(), req, resp)

		Expect(resp.RequiresReplace).To(BeFalse())
	})

	t.Run("Requires replace when state value is missing from plan", func(t *testing.T) {
		RegisterTestingT(t)
		resp := &setplanmodifier.RequiresReplaceIfFuncResponse{}
		req := planmodifier.SetRequest{
			ConfigValue: types.SetValueMust(types.StringType, []attr.Value{types.StringValue("a")}),
			StateValue:  types.SetValueMust(types.StringType, []attr.Value{types.StringValue("a"), types.StringValue("b")}),
			PlanValue:   types.SetValueMust(types.StringType, []attr.Value{types.StringValue("a")}),
		}

		ValuesAreRemovedFromState(t.Context(), req, resp)

		Expect(resp.RequiresReplace).To(BeTrue())
	})

	t.Run("Does not require replace when state is empty", func(t *testing.T) {
		RegisterTestingT(t)
		resp := &setplanmodifier.RequiresReplaceIfFuncResponse{}
		req := planmodifier.SetRequest{
			ConfigValue: types.SetValueMust(types.StringType, []attr.Value{types.StringValue("a")}),
			StateValue:  types.SetValueMust(types.StringType, []attr.Value{}),
			PlanValue:   types.SetValueMust(types.StringType, []attr.Value{types.StringValue("a")}),
		}

		ValuesAreRemovedFromState(t.Context(), req, resp)

		Expect(resp.RequiresReplace).To(BeFalse())
	})

	t.Run("Does not require replace when values are added to plan", func(t *testing.T) {
		RegisterTestingT(t)
		resp := &setplanmodifier.RequiresReplaceIfFuncResponse{}
		req := planmodifier.SetRequest{
			ConfigValue: types.SetValueMust(types.StringType, []attr.Value{types.StringValue("a"), types.StringValue("b"), types.StringValue("c")}),
			StateValue:  types.SetValueMust(types.StringType, []attr.Value{types.StringValue("a"), types.StringValue("b")}),
			PlanValue:   types.SetValueMust(types.StringType, []attr.Value{types.StringValue("a"), types.StringValue("b"), types.StringValue("c")}),
		}

		ValuesAreRemovedFromState(t.Context(), req, resp)

		Expect(resp.RequiresReplace).To(BeFalse())
	})
}
