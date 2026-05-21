// Package tfset converts Terraform Plugin Framework set values into the plain
// Go collections used by the resource conversion code, and builds membership
// indexes over those values.
package tfset

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Strings returns the elements of a Terraform string set as a slice. A null or
// empty set yields a non-nil, empty slice and no diagnostics.
func Strings(ctx context.Context, set types.Set) ([]string, diag.Diagnostics) {
	if len(set.Elements()) == 0 {
		return []string{}, nil
	}

	slice := make([]string, 0, len(set.Elements()))
	diags := set.ElementsAs(ctx, &slice, false)
	return slice, diags
}

// Index builds a membership set over values for O(1) lookups. Duplicate values
// collapse into a single entry.
func Index(values []string) map[string]struct{} {
	index := make(map[string]struct{}, len(values))
	for _, v := range values {
		index[v] = struct{}{}
	}
	return index
}
