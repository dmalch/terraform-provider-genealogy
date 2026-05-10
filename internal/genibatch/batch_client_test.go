package genibatch

import (
	"errors"
	"testing"

	"github.com/dmalch/terraform-provider-genealogy/internal/geni"
)

func TestFulfillDocumentRequests_MissingIDReturnsResourceNotFound(t *testing.T) {
	missing := documentAsyncRequest{
		Id:       "document-missing",
		Response: make(chan *geni.DocumentResponse, 1),
		Error:    make(chan error, 1),
	}

	fulfillDocumentRequests([]documentAsyncRequest{missing}, nil)

	select {
	case err := <-missing.Error:
		if !errors.Is(err, geni.ErrResourceNotFound) {
			t.Fatalf("expected error to wrap geni.ErrResourceNotFound, got: %v", err)
		}
	case <-missing.Response:
		t.Fatal("expected error on missing ID, got a response")
	}
}

func TestFulfillDocumentRequests_MixedHitsAndMisses(t *testing.T) {
	hit := documentAsyncRequest{
		Id:       "document-hit",
		Response: make(chan *geni.DocumentResponse, 1),
		Error:    make(chan error, 1),
	}
	miss := documentAsyncRequest{
		Id:       "document-miss",
		Response: make(chan *geni.DocumentResponse, 1),
		Error:    make(chan error, 1),
	}

	results := []geni.DocumentResponse{{Id: "document-hit", Title: "found"}}
	fulfillDocumentRequests([]documentAsyncRequest{hit, miss}, results)

	select {
	case res := <-hit.Response:
		if res == nil || res.Id != "document-hit" {
			t.Fatalf("expected hit response with Id=document-hit, got: %#v", res)
		}
	case err := <-hit.Error:
		t.Fatalf("expected response for hit, got error: %v", err)
	}

	select {
	case err := <-miss.Error:
		if !errors.Is(err, geni.ErrResourceNotFound) {
			t.Fatalf("expected miss error to wrap geni.ErrResourceNotFound, got: %v", err)
		}
	case <-miss.Response:
		t.Fatal("expected error for miss, got a response")
	}
}

func TestFulfillUnionRequests_MissingIDReturnsResourceNotFound(t *testing.T) {
	missing := unionAsyncRequest{
		Id:       "union-missing",
		Response: make(chan *geni.UnionResponse, 1),
		Error:    make(chan error, 1),
	}

	fulfillUnionRequests([]unionAsyncRequest{missing}, nil)

	select {
	case err := <-missing.Error:
		if !errors.Is(err, geni.ErrResourceNotFound) {
			t.Fatalf("expected error to wrap geni.ErrResourceNotFound, got: %v", err)
		}
	case <-missing.Response:
		t.Fatal("expected error on missing ID, got a response")
	}
}

func TestFulfillProfileRequests_MissingIDReturnsResourceNotFound(t *testing.T) {
	missing := profileAsyncRequest{
		Id:       "profile-missing",
		Response: make(chan *geni.ProfileResponse, 1),
		Error:    make(chan error, 1),
	}

	fulfillProfileRequests([]profileAsyncRequest{missing}, nil)

	select {
	case err := <-missing.Error:
		if !errors.Is(err, geni.ErrResourceNotFound) {
			t.Fatalf("expected error to wrap geni.ErrResourceNotFound, got: %v", err)
		}
	case <-missing.Response:
		t.Fatal("expected error on missing ID, got a response")
	}
}
