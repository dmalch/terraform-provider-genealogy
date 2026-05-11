// Package listresource hosts the shared building blocks for Terraform 1.14
// List Resources backed by the Geni API. Each per-resource subpackage
// (profile/, document/) wires a list.ListResource to one of Geni's paginated
// list endpoints and reuses the managed-resource translators to populate
// query results.
package listresource

import (
	"context"
	"iter"

	"github.com/hashicorp/terraform-plugin-framework/list"
)

// paginate yields every element from a paginated API endpoint, fetching each
// page lazily and honoring consumer cancellation through the iter.Seq push
// contract. The fetchPage callback returns the elements for the given 1-based
// page plus the total count of elements across all pages; iteration stops as
// soon as the running count reaches the total or a page comes back empty.
//
// onError translates a fetchPage error into a single terminal ListResult that
// carries an error diagnostic; project translates one API element into a
// ListResult (returning false from project halts iteration silently — the
// caller is expected to have already attached any diagnostic to the result).
func Paginate[T any](
	ctx context.Context,
	fetchPage func(ctx context.Context, page int) ([]T, int, error),
	onError func(error) list.ListResult,
	project func(T) (list.ListResult, bool),
) iter.Seq[list.ListResult] {
	return func(push func(list.ListResult) bool) {
		seen := 0
		for page := 1; ; page++ {
			items, total, err := fetchPage(ctx, page)
			if err != nil {
				push(onError(err))
				return
			}
			if len(items) == 0 {
				return
			}
			for i := range items {
				result, ok := project(items[i])
				if !ok {
					return
				}
				if !push(result) {
					return
				}
			}
			seen += len(items)
			if seen >= total {
				return
			}
		}
	}
}
