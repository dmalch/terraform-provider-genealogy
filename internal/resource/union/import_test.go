package union

import (
	"context"
	"errors"
	"testing"

	. "github.com/onsi/gomega"

	"github.com/dmalch/terraform-provider-genealogy/internal/geni"
)

func TestResolveUnionImport(t *testing.T) {
	t.Run("Not-found from fetch produces an error diagnostic", func(t *testing.T) {
		RegisterTestingT(t)
		fetch := func(_ context.Context, _ string) (*geni.UnionResponse, error) {
			return nil, geni.ErrResourceNotFound
		}

		_, _, diags := resolveUnionImport(t.Context(), "union-missing", fetch)

		Expect(diags.HasError()).To(BeTrue())
	})

	t.Run("Transport error is surfaced as an error diagnostic", func(t *testing.T) {
		RegisterTestingT(t)
		fetch := func(_ context.Context, _ string) (*geni.UnionResponse, error) {
			return nil, errors.New("network exploded")
		}

		_, _, diags := resolveUnionImport(t.Context(), "union-x", fetch)

		Expect(diags.HasError()).To(BeTrue())
	})

	t.Run("Happy path populates identity from the fetched response", func(t *testing.T) {
		RegisterTestingT(t)
		fetch := func(_ context.Context, id string) (*geni.UnionResponse, error) {
			return &geni.UnionResponse{Id: id}, nil
		}

		_, identity, diags := resolveUnionImport(t.Context(), "union-42", fetch)

		Expect(diags.HasError()).To(BeFalse())
		Expect(identity.ID.ValueString()).To(Equal("union-42"))
	})
}
