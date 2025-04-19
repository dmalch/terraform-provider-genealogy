package geniplanmodifier

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
)

func ValuesAreRemovedFromState(_ context.Context, req planmodifier.SetRequest, resp *setplanmodifier.RequiresReplaceIfFuncResponse) {
	if req.ConfigValue.IsNull() {
		return
	}

	// If plan contains all values from state, no need to replace
	for _, v := range req.StateValue.Elements() {
		// Check if the value is in the plan
		if !contains(req.PlanValue.Elements(), v) {
			resp.RequiresReplace = true
			break
		}
	}
}

func contains(elements []attr.Value, v attr.Value) bool {
	for _, p := range elements {
		if v.Equal(p) {
			return true
		}
	}
	return false
}
