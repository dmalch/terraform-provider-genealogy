package profile

import (
	"context"
	"errors"
	"testing"

	. "github.com/onsi/gomega"

	"github.com/dmalch/terraform-provider-genealogy/internal/geni"
)

func TestResolveProfileImport(t *testing.T) {
	t.Run("Not-found from fetch produces an error diagnostic", func(t *testing.T) {
		RegisterTestingT(t)
		fetch := func(_ context.Context, _ string) (*geni.ProfileResponse, error) {
			return nil, geni.ErrResourceNotFound
		}

		_, _, diags := resolveProfileImport(t.Context(), "profile-missing", fetch)

		Expect(diags.HasError()).To(BeTrue())
	})

	t.Run("Transport error is surfaced as an error diagnostic", func(t *testing.T) {
		RegisterTestingT(t)
		fetch := func(_ context.Context, _ string) (*geni.ProfileResponse, error) {
			return nil, errors.New("network exploded")
		}

		_, _, diags := resolveProfileImport(t.Context(), "profile-x", fetch)

		Expect(diags.HasError()).To(BeTrue())
	})

	t.Run("Happy path populates identity from the fetched response", func(t *testing.T) {
		RegisterTestingT(t)
		fetch := func(_ context.Context, id string) (*geni.ProfileResponse, error) {
			return &geni.ProfileResponse{Id: id}, nil
		}

		_, identity, diags := resolveProfileImport(t.Context(), "profile-42", fetch)

		Expect(diags.HasError()).To(BeFalse())
		Expect(identity.ID.ValueString()).To(Equal("profile-42"))
	})
}
