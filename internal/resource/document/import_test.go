package document

import (
	"context"
	"errors"
	"testing"

	. "github.com/onsi/gomega"

	"github.com/dmalch/terraform-provider-genealogy/internal/geni"
)

func TestValidateDocumentImportID(t *testing.T) {
	t.Run("Not-found from fetch produces an error diagnostic", func(t *testing.T) {
		RegisterTestingT(t)
		fetch := func(_ context.Context, _ string) (*geni.DocumentResponse, error) {
			return nil, geni.ErrResourceNotFound
		}

		diags := validateDocumentImportID(t.Context(), "document-missing", fetch)

		Expect(diags.HasError()).To(BeTrue())
	})

	t.Run("Transport error is surfaced as an error diagnostic", func(t *testing.T) {
		RegisterTestingT(t)
		fetch := func(_ context.Context, _ string) (*geni.DocumentResponse, error) {
			return nil, errors.New("network exploded")
		}

		diags := validateDocumentImportID(t.Context(), "document-x", fetch)

		Expect(diags.HasError()).To(BeTrue())
	})

	t.Run("Empty response Id is treated as not-found", func(t *testing.T) {
		RegisterTestingT(t)
		fetch := func(_ context.Context, _ string) (*geni.DocumentResponse, error) {
			return &geni.DocumentResponse{}, nil
		}

		diags := validateDocumentImportID(t.Context(), "document-missing", fetch)

		Expect(diags.HasError()).To(BeTrue())
	})

	t.Run("Successful fetch yields no diagnostics", func(t *testing.T) {
		RegisterTestingT(t)
		fetch := func(_ context.Context, id string) (*geni.DocumentResponse, error) {
			return &geni.DocumentResponse{Id: id}, nil
		}

		diags := validateDocumentImportID(t.Context(), "document-42", fetch)

		Expect(diags.HasError()).To(BeFalse())
	})
}
