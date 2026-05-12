package event

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// EventNeedsDatePreWipe returns true when applying a Geni profile/union Update
// requires a pre-wipe PATCH that nulls the event's date object before the
// regular PATCH writes the new state. Without this, Geni silently drops
// requests to clear individual numeric/boolean/string sub-fields inside
// `date` — the API deep-merges nested objects per-key (#94).
//
// A pre-wipe is needed iff state and plan both carry a non-null date AND at
// least one date sub-field is non-null in state but null in plan. Whole-date
// clears (plan.date == null) and pure adds work via the existing PATCH path.
func EventNeedsDatePreWipe(state, plan types.Object) bool {
	stateDate, ok := dateOf(state)
	if !ok {
		return false
	}
	planDate, ok := dateOf(plan)
	if !ok {
		return false
	}

	stateAttrs := stateDate.Attributes()
	planAttrs := planDate.Attributes()

	for name, stateVal := range stateAttrs {
		if stateVal.IsNull() {
			continue
		}
		if planVal, present := planAttrs[name]; present && planVal.IsNull() {
			return true
		}
	}
	return false
}

// dateOf returns the date sub-object of the supplied event, or false if the
// event is null/unknown or has no non-null date.
func dateOf(event types.Object) (types.Object, bool) {
	if event.IsNull() || event.IsUnknown() {
		return types.Object{}, false
	}
	dateVal, ok := event.Attributes()["date"]
	if !ok {
		return types.Object{}, false
	}
	dateObj, ok := dateVal.(types.Object)
	if !ok || dateObj.IsNull() || dateObj.IsUnknown() {
		return types.Object{}, false
	}
	return dateObj, true
}
