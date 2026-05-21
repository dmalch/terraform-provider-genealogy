package profile

import (
	"context"
	"errors"
	"testing"

	. "github.com/onsi/gomega"

	geniprofile "github.com/dmalch/go-geni/profile"
)

func TestFollowMergedInto(t *testing.T) {
	t.Run("Returns the input unchanged when the profile is live", func(t *testing.T) {
		RegisterTestingT(t)

		initial := &geniprofile.Profile{ID: "profile-1", Deleted: false}
		fetchCalls := 0
		fetch := func(_ context.Context, _ string) (*geniprofile.Profile, error) {
			fetchCalls++
			return nil, errors.New("unexpected fetch")
		}

		result, err := FollowMergedInto(t.Context(), initial, fetch, 10)

		Expect(err).NotTo(HaveOccurred())
		Expect(result).To(BeIdenticalTo(initial))
		Expect(fetchCalls).To(Equal(0))
	})

	t.Run("Returns the input when deleted but merged_into is empty", func(t *testing.T) {
		RegisterTestingT(t)

		initial := &geniprofile.Profile{ID: "profile-1", Deleted: true, MergedInto: ""}
		fetchCalls := 0
		fetch := func(_ context.Context, _ string) (*geniprofile.Profile, error) {
			fetchCalls++
			return nil, errors.New("unexpected fetch")
		}

		result, err := FollowMergedInto(t.Context(), initial, fetch, 10)

		Expect(err).NotTo(HaveOccurred())
		Expect(result).To(BeIdenticalTo(initial))
		Expect(fetchCalls).To(Equal(0))
	})

	t.Run("Follows a single hop to the live successor", func(t *testing.T) {
		RegisterTestingT(t)

		initial := &geniprofile.Profile{ID: "profile-1", Deleted: true, MergedInto: "profile-2"}
		successor := &geniprofile.Profile{ID: "profile-2", Deleted: false}
		fetchCalls := 0
		fetch := func(_ context.Context, id string) (*geniprofile.Profile, error) {
			fetchCalls++
			Expect(id).To(Equal("profile-2"))
			return successor, nil
		}

		result, err := FollowMergedInto(t.Context(), initial, fetch, 10)

		Expect(err).NotTo(HaveOccurred())
		Expect(result).To(BeIdenticalTo(successor))
		Expect(fetchCalls).To(Equal(1))
	})

	t.Run("Walks a multi-hop chain until reaching a live profile", func(t *testing.T) {
		RegisterTestingT(t)

		initial := &geniprofile.Profile{ID: "profile-1", Deleted: true, MergedInto: "profile-2"}
		chain := map[string]*geniprofile.Profile{
			"profile-2": {ID: "profile-2", Deleted: true, MergedInto: "profile-3"},
			"profile-3": {ID: "profile-3", Deleted: true, MergedInto: "profile-4"},
			"profile-4": {ID: "profile-4", Deleted: false},
		}
		var visited []string
		fetch := func(_ context.Context, id string) (*geniprofile.Profile, error) {
			visited = append(visited, id)
			p, ok := chain[id]
			if !ok {
				return nil, errors.New("not found")
			}
			return p, nil
		}

		result, err := FollowMergedInto(t.Context(), initial, fetch, 10)

		Expect(err).NotTo(HaveOccurred())
		Expect(result.ID).To(Equal("profile-4"))
		Expect(result.Deleted).To(BeFalse())
		Expect(visited).To(Equal([]string{"profile-2", "profile-3", "profile-4"}))
	})

	t.Run("Stops at maxHops even if the chain is still deleted", func(t *testing.T) {
		RegisterTestingT(t)

		initial := &geniprofile.Profile{ID: "profile-1", Deleted: true, MergedInto: "profile-2"}
		chain := map[string]*geniprofile.Profile{
			"profile-2": {ID: "profile-2", Deleted: true, MergedInto: "profile-3"},
			"profile-3": {ID: "profile-3", Deleted: true, MergedInto: "profile-4"},
			"profile-4": {ID: "profile-4", Deleted: true, MergedInto: "profile-5"},
		}
		fetch := func(_ context.Context, id string) (*geniprofile.Profile, error) {
			return chain[id], nil
		}

		result, err := FollowMergedInto(t.Context(), initial, fetch, 2)

		Expect(err).NotTo(HaveOccurred())
		Expect(result.ID).To(Equal("profile-3"))
		Expect(result.Deleted).To(BeTrue())
	})

	t.Run("Returns the fetch error and the last successfully resolved profile", func(t *testing.T) {
		RegisterTestingT(t)

		initial := &geniprofile.Profile{ID: "profile-1", Deleted: true, MergedInto: "profile-2"}
		boom := errors.New("transport blew up")
		fetch := func(_ context.Context, _ string) (*geniprofile.Profile, error) {
			return nil, boom
		}

		result, err := FollowMergedInto(t.Context(), initial, fetch, 10)

		Expect(err).To(MatchError(boom))
		Expect(result).To(BeIdenticalTo(initial))
	})
}
