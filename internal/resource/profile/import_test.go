package profile

import (
	"context"
	"errors"
	"testing"

	. "github.com/onsi/gomega"

	"github.com/dmalch/go-geni"
)

func TestValidateProfileImportID(t *testing.T) {
	t.Run("Not-found from fetch produces an error diagnostic", func(t *testing.T) {
		RegisterTestingT(t)
		fetch := func(_ context.Context, _ string) (*geni.ProfileResponse, error) {
			return nil, geni.ErrResourceNotFound
		}

		diags := validateProfileImportID(t.Context(), "profile-missing", fetch)

		Expect(diags.HasError()).To(BeTrue())
	})

	t.Run("Transport error is surfaced as an error diagnostic", func(t *testing.T) {
		RegisterTestingT(t)
		fetch := func(_ context.Context, _ string) (*geni.ProfileResponse, error) {
			return nil, errors.New("network exploded")
		}

		diags := validateProfileImportID(t.Context(), "profile-x", fetch)

		Expect(diags.HasError()).To(BeTrue())
	})

	t.Run("Empty response Id is treated as not-found", func(t *testing.T) {
		RegisterTestingT(t)
		fetch := func(_ context.Context, _ string) (*geni.ProfileResponse, error) {
			return &geni.ProfileResponse{}, nil
		}

		diags := validateProfileImportID(t.Context(), "profile-missing", fetch)

		Expect(diags.HasError()).To(BeTrue())
	})

	t.Run("Successful fetch yields no diagnostics", func(t *testing.T) {
		RegisterTestingT(t)
		fetch := func(_ context.Context, id string) (*geni.ProfileResponse, error) {
			return &geni.ProfileResponse{Id: id}, nil
		}

		diags := validateProfileImportID(t.Context(), "profile-42", fetch)

		Expect(diags.HasError()).To(BeFalse())
	})
}
