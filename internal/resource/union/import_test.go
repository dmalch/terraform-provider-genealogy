package union

import (
	"context"
	"errors"
	"testing"

	"github.com/dmalch/terraform-provider-genealogy/internal/geni"
)

func TestResolveUnionImport_NotFoundProducesError(t *testing.T) {
	fetch := func(_ context.Context, _ string) (*geni.UnionResponse, error) {
		return nil, geni.ErrResourceNotFound
	}

	_, _, diags := resolveUnionImport(context.Background(), "union-missing", fetch)

	if !diags.HasError() {
		t.Fatal("expected error diagnostic when fetch returns ErrResourceNotFound, got none")
	}
}

func TestResolveUnionImport_TransportErrorSurfaced(t *testing.T) {
	fetch := func(_ context.Context, _ string) (*geni.UnionResponse, error) {
		return nil, errors.New("network exploded")
	}

	_, _, diags := resolveUnionImport(context.Background(), "union-x", fetch)

	if !diags.HasError() {
		t.Fatal("expected error diagnostic when fetch returns a non-not-found error, got none")
	}
}

func TestResolveUnionImport_HappyPathPopulatesIdentity(t *testing.T) {
	fetch := func(_ context.Context, id string) (*geni.UnionResponse, error) {
		return &geni.UnionResponse{Id: id}, nil
	}

	_, identity, diags := resolveUnionImport(context.Background(), "union-42", fetch)

	if diags.HasError() {
		t.Fatalf("expected no diagnostics on happy path, got: %v", diags)
	}
	if identity.ID.ValueString() != "union-42" {
		t.Fatalf("expected identity.ID=union-42, got %q", identity.ID.ValueString())
	}
}
