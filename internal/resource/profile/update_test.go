package profile

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestAddedProjects(t *testing.T) {
	t.Run("an unchanged membership is not re-linked", func(t *testing.T) {
		RegisterTestingT(t)

		// The regression this exists for. AddProfile is a write only the
		// project's OWNER may perform, so re-sending a membership somebody
		// else created fails with "access denied" and fails the whole update
		// — after the profile itself already updated successfully.
		linked := []string{"project-28427", "project-4505748", "project-4518482"}

		Expect(addedProjects(linked, linked)).To(BeEmpty())
	})

	t.Run("only the newly added project is linked", func(t *testing.T) {
		RegisterTestingT(t)

		added := addedProjects(
			[]string{"project-4505748"},
			[]string{"project-4505748", "project-4518482"},
		)

		Expect(added).To(ConsistOf("project-4518482"))
	})

	t.Run("every project is linked when state has none", func(t *testing.T) {
		RegisterTestingT(t)

		added := addedProjects(nil, []string{"project-4505748", "project-4518482"})

		Expect(added).To(ConsistOf("project-4505748", "project-4518482"))
	})

	t.Run("a dropped project is not unlinked", func(t *testing.T) {
		RegisterTestingT(t)

		// Update has never unlinked, and this change deliberately does not
		// start: who owns the field is a separate decision.
		added := addedProjects(
			[]string{"project-4505748", "project-4518482"},
			[]string{"project-4505748"},
		)

		Expect(added).To(BeEmpty())
	})

	t.Run("duplicates in the plan are linked once", func(t *testing.T) {
		RegisterTestingT(t)

		added := addedProjects(nil, []string{"project-4518482", "project-4518482"})

		// tfset.Index collapses duplicates, but the plan side is a list — a
		// repeated id would otherwise be sent twice.
		Expect(added).To(HaveLen(1))
	})
}
