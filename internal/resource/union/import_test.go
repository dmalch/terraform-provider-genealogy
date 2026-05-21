package union

import (
	"context"
	"errors"
	"testing"

	. "github.com/onsi/gomega"

	"github.com/dmalch/go-geni"
	geniunion "github.com/dmalch/go-geni/union"
)

func TestValidateUnionImportID(t *testing.T) {
	t.Run("Not-found from fetch produces an error diagnostic", func(t *testing.T) {
		RegisterTestingT(t)
		fetch := func(_ context.Context, _ string) (*geniunion.Union, error) {
			return nil, geni.ErrResourceNotFound
		}

		diags := validateUnionImportID(t.Context(), "union-missing", fetch)

		Expect(diags.HasError()).To(BeTrue())
	})

	t.Run("Transport error is surfaced as an error diagnostic", func(t *testing.T) {
		RegisterTestingT(t)
		fetch := func(_ context.Context, _ string) (*geniunion.Union, error) {
			return nil, errors.New("network exploded")
		}

		diags := validateUnionImportID(t.Context(), "union-x", fetch)

		Expect(diags.HasError()).To(BeTrue())
	})

	t.Run("Empty response Id is treated as not-found", func(t *testing.T) {
		RegisterTestingT(t)
		fetch := func(_ context.Context, _ string) (*geniunion.Union, error) {
			return &geniunion.Union{}, nil
		}

		diags := validateUnionImportID(t.Context(), "union-missing", fetch)

		Expect(diags.HasError()).To(BeTrue())
	})

	t.Run("Successful fetch yields no diagnostics", func(t *testing.T) {
		RegisterTestingT(t)
		fetch := func(_ context.Context, id string) (*geniunion.Union, error) {
			return &geniunion.Union{ID: id}, nil
		}

		diags := validateUnionImportID(t.Context(), "union-42", fetch)

		Expect(diags.HasError()).To(BeFalse())
	})
}
