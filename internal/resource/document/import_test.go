package document

import (
	"context"
	"errors"
	"testing"

	"github.com/dmalch/terraform-provider-genealogy/internal/geni"
)

func TestResolveDocumentImport_NotFoundProducesError(t *testing.T) {
	fetch := func(_ context.Context, _ string) (*geni.DocumentResponse, error) {
		return nil, geni.ErrResourceNotFound
	}

	_, _, diags := resolveDocumentImport(context.Background(), "document-missing", fetch)

	if !diags.HasError() {
		t.Fatal("expected error diagnostic when fetch returns ErrResourceNotFound, got none")
	}
}

func TestResolveDocumentImport_TransportErrorSurfaced(t *testing.T) {
	sentinel := errors.New("network exploded")
	fetch := func(_ context.Context, _ string) (*geni.DocumentResponse, error) {
		return nil, sentinel
	}

	_, _, diags := resolveDocumentImport(context.Background(), "document-x", fetch)

	if !diags.HasError() {
		t.Fatal("expected error diagnostic when fetch returns a non-not-found error, got none")
	}
}

func TestResolveDocumentImport_HappyPathPopulatesStateAndIdentity(t *testing.T) {
	fetch := func(_ context.Context, id string) (*geni.DocumentResponse, error) {
		return &geni.DocumentResponse{Id: id, Title: "My Doc"}, nil
	}

	state, identity, diags := resolveDocumentImport(context.Background(), "document-42", fetch)

	if diags.HasError() {
		t.Fatalf("expected no diagnostics on happy path, got: %v", diags)
	}
	if identity.ID.ValueString() != "document-42" {
		t.Fatalf("expected identity.ID=document-42, got %q", identity.ID.ValueString())
	}
	if state.ID.ValueString() != "document-42" {
		t.Fatalf("expected state.ID=document-42, got %q", state.ID.ValueString())
	}
	if state.Title.ValueString() != "My Doc" {
		t.Fatalf("expected state.Title=My Doc, got %q", state.Title.ValueString())
	}
}
