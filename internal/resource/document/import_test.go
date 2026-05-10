package document

import (
	"context"
	"errors"
	"testing"

	. "github.com/onsi/gomega"

	"github.com/dmalch/terraform-provider-genealogy/internal/geni"
)

func TestResolveDocumentImport(t *testing.T) {
	t.Run("Not-found from fetch produces an error diagnostic", func(t *testing.T) {
		RegisterTestingT(t)
		fetch := func(_ context.Context, _ string) (*geni.DocumentResponse, error) {
			return nil, geni.ErrResourceNotFound
		}

		_, _, diags := resolveDocumentImport(t.Context(), "document-missing", fetch)

		Expect(diags.HasError()).To(BeTrue())
	})

	t.Run("Transport error is surfaced as an error diagnostic", func(t *testing.T) {
		RegisterTestingT(t)
		fetch := func(_ context.Context, _ string) (*geni.DocumentResponse, error) {
			return nil, errors.New("network exploded")
		}

		_, _, diags := resolveDocumentImport(t.Context(), "document-x", fetch)

		Expect(diags.HasError()).To(BeTrue())
	})

	t.Run("Happy path populates state and identity from the fetched response", func(t *testing.T) {
		RegisterTestingT(t)
		fetch := func(_ context.Context, id string) (*geni.DocumentResponse, error) {
			return &geni.DocumentResponse{Id: id, Title: "My Doc"}, nil
		}

		state, identity, diags := resolveDocumentImport(t.Context(), "document-42", fetch)

		Expect(diags.HasError()).To(BeFalse())
		Expect(identity.ID.ValueString()).To(Equal("document-42"))
		Expect(state.ID.ValueString()).To(Equal("document-42"))
		Expect(state.Title.ValueString()).To(Equal("My Doc"))
	})
}
