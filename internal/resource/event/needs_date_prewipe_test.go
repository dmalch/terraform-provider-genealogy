package event

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	. "github.com/onsi/gomega"
)

// TestEventNeedsDatePreWipe verifies the rule that decides whether the Geni
// API needs a pre-wipe PATCH before the main Update can succeed at clearing
// individual date sub-fields (#94). The rule:
//
//   - State and plan both carry a non-null date, AND
//   - At least one sub-field is non-null in state but null in plan.
//
// All other cases (no date in state, plan removing the whole date, plan
// matching state, plan adding new sub-fields) do not need a pre-wipe.
func TestEventNeedsDatePreWipe(t *testing.T) {
	t.Run("Returns false when state has no event", func(t *testing.T) {
		RegisterTestingT(t)
		state := types.ObjectNull(EventModelAttributeTypes())
		plan := eventWithDate(map[string]attr.Value{
			"year":  types.Int32Value(1900),
			"month": types.Int32Null(),
		})
		Expect(EventNeedsDatePreWipe(state, plan)).To(BeFalse())
	})

	t.Run("Returns false when plan has no event", func(t *testing.T) {
		RegisterTestingT(t)
		state := eventWithDate(map[string]attr.Value{
			"year":      types.Int32Value(1900),
			"end_month": types.Int32Value(7),
		})
		plan := types.ObjectNull(EventModelAttributeTypes())
		Expect(EventNeedsDatePreWipe(state, plan)).To(BeFalse())
	})

	t.Run("Returns false when plan clears the whole date object", func(t *testing.T) {
		RegisterTestingT(t)
		state := eventWithDate(map[string]attr.Value{
			"year":      types.Int32Value(1900),
			"end_month": types.Int32Value(7),
		})
		plan := eventWithNullDate()
		Expect(EventNeedsDatePreWipe(state, plan)).To(BeFalse())
	})

	t.Run("Returns false when state and plan dates are equal", func(t *testing.T) {
		RegisterTestingT(t)
		fields := map[string]attr.Value{
			"year":      types.Int32Value(1900),
			"end_month": types.Int32Value(7),
		}
		Expect(EventNeedsDatePreWipe(eventWithDate(fields), eventWithDate(fields))).To(BeFalse())
	})

	t.Run("Returns false when plan only adds new sub-fields", func(t *testing.T) {
		RegisterTestingT(t)
		state := eventWithDate(map[string]attr.Value{
			"year": types.Int32Value(1900),
		})
		plan := eventWithDate(map[string]attr.Value{
			"year":  types.Int32Value(1900),
			"month": types.Int32Value(7),
		})
		Expect(EventNeedsDatePreWipe(state, plan)).To(BeFalse())
	})

	t.Run("Returns true when plan clears a numeric sub-field", func(t *testing.T) {
		RegisterTestingT(t)
		state := eventWithDate(map[string]attr.Value{
			"year":      types.Int32Value(1752),
			"end_year":  types.Int32Value(1799),
			"end_month": types.Int32Value(7),
			"range":     types.StringValue("between"),
		})
		plan := eventWithDate(map[string]attr.Value{
			"year":     types.Int32Value(1752),
			"end_year": types.Int32Value(1799),
			"range":    types.StringValue("between"),
		})
		Expect(EventNeedsDatePreWipe(state, plan)).To(BeTrue())
	})

	t.Run("Returns true when plan clears a boolean sub-field", func(t *testing.T) {
		RegisterTestingT(t)
		state := eventWithDate(map[string]attr.Value{
			"year":      types.Int32Value(1900),
			"end_circa": types.BoolValue(true),
		})
		plan := eventWithDate(map[string]attr.Value{
			"year": types.Int32Value(1900),
		})
		Expect(EventNeedsDatePreWipe(state, plan)).To(BeTrue())
	})

	t.Run("Returns true when plan clears a string sub-field", func(t *testing.T) {
		RegisterTestingT(t)
		state := eventWithDate(map[string]attr.Value{
			"year":  types.Int32Value(1900),
			"range": types.StringValue("before"),
		})
		plan := eventWithDate(map[string]attr.Value{
			"year": types.Int32Value(1900),
		})
		Expect(EventNeedsDatePreWipe(state, plan)).To(BeTrue())
	})
}

// eventWithDate builds an event object whose date sub-object carries the
// supplied fields; any field not in `fields` is null.
func eventWithDate(fields map[string]attr.Value) types.Object {
	full := map[string]attr.Value{
		"range":     types.StringNull(),
		"circa":     types.BoolNull(),
		"day":       types.Int32Null(),
		"month":     types.Int32Null(),
		"year":      types.Int32Null(),
		"end_circa": types.BoolNull(),
		"end_day":   types.Int32Null(),
		"end_month": types.Int32Null(),
		"end_year":  types.Int32Null(),
	}
	for k, v := range fields {
		full[k] = v
	}
	dateObj := types.ObjectValueMust(DateRangeModelAttributeTypes(), full)
	return types.ObjectValueMust(EventModelAttributeTypes(), map[string]attr.Value{
		"name":        types.StringNull(),
		"description": types.StringNull(),
		"date":        dateObj,
		"location":    types.ObjectNull(LocationModelAttributeTypes()),
	})
}

func eventWithNullDate() types.Object {
	return types.ObjectValueMust(EventModelAttributeTypes(), map[string]attr.Value{
		"name":        types.StringNull(),
		"description": types.StringNull(),
		"date":        types.ObjectNull(DateRangeModelAttributeTypes()),
		"location":    types.ObjectNull(LocationModelAttributeTypes()),
	})
}
