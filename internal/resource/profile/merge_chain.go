package profile

import (
	"context"

	"github.com/dmalch/go-geni"
)

// FollowMergedInto walks the merged_into chain starting from initial, fetching
// each successor through fetch, while the current profile is marked as deleted
// and points at a non-empty successor. It stops at the first live profile, when
// merged_into is empty, when maxHops fetches have happened, or when fetch
// returns an error.
//
// The returned profile may still be Deleted=true if the chain ran out before a
// live profile was found; callers decide how to surface that condition. When
// fetch errors, the last successfully resolved profile is returned alongside
// the error so callers can include identity information in diagnostics.
func FollowMergedInto(
	ctx context.Context,
	initial *geni.ProfileResponse,
	fetch func(context.Context, string) (*geni.ProfileResponse, error),
	maxHops int,
) (*geni.ProfileResponse, error) {
	current := initial
	for i := 0; i < maxHops && current.Deleted && current.MergedInto != ""; i++ {
		next, err := fetch(ctx, current.MergedInto)
		if err != nil {
			return current, err
		}
		current = next
	}
	return current, nil
}
