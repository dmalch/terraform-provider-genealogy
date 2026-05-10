package profile

import (
	"context"
	"errors"
	"testing"

	"github.com/dmalch/terraform-provider-genealogy/internal/geni"
)

func TestResolveProfileImport_NotFoundProducesError(t *testing.T) {
	fetch := func(_ context.Context, _ string) (*geni.ProfileResponse, error) {
		return nil, geni.ErrResourceNotFound
	}

	_, _, diags := resolveProfileImport(context.Background(), "profile-missing", fetch)

	if !diags.HasError() {
		t.Fatal("expected error diagnostic when fetch returns ErrResourceNotFound, got none")
	}
}

func TestResolveProfileImport_TransportErrorSurfaced(t *testing.T) {
	fetch := func(_ context.Context, _ string) (*geni.ProfileResponse, error) {
		return nil, errors.New("network exploded")
	}

	_, _, diags := resolveProfileImport(context.Background(), "profile-x", fetch)

	if !diags.HasError() {
		t.Fatal("expected error diagnostic when fetch returns a non-not-found error, got none")
	}
}

func TestResolveProfileImport_HappyPathPopulatesIdentity(t *testing.T) {
	fetch := func(_ context.Context, id string) (*geni.ProfileResponse, error) {
		return &geni.ProfileResponse{Id: id}, nil
	}

	_, identity, diags := resolveProfileImport(context.Background(), "profile-42", fetch)

	if diags.HasError() {
		t.Fatalf("expected no diagnostics on happy path, got: %v", diags)
	}
	if identity.ID.ValueString() != "profile-42" {
		t.Fatalf("expected identity.ID=profile-42, got %q", identity.ID.ValueString())
	}
}
