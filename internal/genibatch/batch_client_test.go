package genibatch

import (
	"errors"
	"testing"

	. "github.com/onsi/gomega"

	"github.com/dmalch/go-geni"
)

func TestFulfillDocumentRequests(t *testing.T) {
	t.Run("Missing ID in bulk response surfaces ErrResourceNotFound", func(t *testing.T) {
		RegisterTestingT(t)
		missing := documentAsyncRequest{
			Id:       "document-missing",
			Response: make(chan *geni.DocumentResponse, 1),
			Error:    make(chan error, 1),
		}

		fulfillDocumentRequests([]documentAsyncRequest{missing}, nil)

		Expect(missing.Response).ToNot(Receive())
		var err error
		Expect(missing.Error).To(Receive(&err))
		Expect(errors.Is(err, geni.ErrResourceNotFound)).To(BeTrue())
	})

	t.Run("Mixed hits and misses route each request correctly", func(t *testing.T) {
		RegisterTestingT(t)
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

		fulfillDocumentRequests(
			[]documentAsyncRequest{hit, miss},
			[]geni.DocumentResponse{{Id: "document-hit", Title: "found"}},
		)

		var hitResp *geni.DocumentResponse
		Expect(hit.Response).To(Receive(&hitResp))
		Expect(hitResp.Id).To(Equal("document-hit"))
		Expect(hit.Error).ToNot(Receive())

		var missErr error
		Expect(miss.Error).To(Receive(&missErr))
		Expect(errors.Is(missErr, geni.ErrResourceNotFound)).To(BeTrue())
		Expect(miss.Response).ToNot(Receive())
	})
}

func TestFulfillUnionRequests(t *testing.T) {
	t.Run("Missing ID in bulk response surfaces ErrResourceNotFound", func(t *testing.T) {
		RegisterTestingT(t)
		missing := unionAsyncRequest{
			Id:       "union-missing",
			Response: make(chan *geni.UnionResponse, 1),
			Error:    make(chan error, 1),
		}

		fulfillUnionRequests([]unionAsyncRequest{missing}, nil)

		var err error
		Expect(missing.Error).To(Receive(&err))
		Expect(errors.Is(err, geni.ErrResourceNotFound)).To(BeTrue())
	})
}

func TestFulfillProfileRequests(t *testing.T) {
	t.Run("Missing ID in bulk response surfaces ErrResourceNotFound", func(t *testing.T) {
		RegisterTestingT(t)
		missing := profileAsyncRequest{
			Id:       "profile-missing",
			Response: make(chan *geni.ProfileResponse, 1),
			Error:    make(chan error, 1),
		}

		fulfillProfileRequests([]profileAsyncRequest{missing}, nil)

		var err error
		Expect(missing.Error).To(Receive(&err))
		Expect(errors.Is(err, geni.ErrResourceNotFound)).To(BeTrue())
	})
}
