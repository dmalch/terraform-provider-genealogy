package union

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestMissingEdges(t *testing.T) {
	t.Run("Returns only the planned ids not already present", func(t *testing.T) {
		RegisterTestingT(t)

		current := map[string]struct{}{"profile-1": {}, "profile-2": {}}
		planned := []string{"profile-1", "profile-2", "profile-3"}

		Expect(missingEdges(planned, current)).To(Equal([]string{"profile-3"}))
	})

	t.Run("Returns empty when every planned id is already present", func(t *testing.T) {
		RegisterTestingT(t)

		current := map[string]struct{}{"profile-1": {}, "profile-2": {}}
		planned := []string{"profile-1", "profile-2"}

		Expect(missingEdges(planned, current)).To(BeEmpty())
	})

	t.Run("Returns all planned ids when none are present", func(t *testing.T) {
		RegisterTestingT(t)

		planned := []string{"profile-1", "profile-2"}

		Expect(missingEdges(planned, map[string]struct{}{})).To(Equal([]string{"profile-1", "profile-2"}))
	})

	t.Run("Preserves the order of the planned ids", func(t *testing.T) {
		RegisterTestingT(t)

		current := map[string]struct{}{"profile-2": {}}
		planned := []string{"profile-3", "profile-2", "profile-1"}

		Expect(missingEdges(planned, current)).To(Equal([]string{"profile-3", "profile-1"}))
	})
}
