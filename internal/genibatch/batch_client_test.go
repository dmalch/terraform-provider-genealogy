package genibatch

import (
	"errors"
	"testing"

	. "github.com/onsi/gomega"

	"github.com/dmalch/go-geni"
	genidocument "github.com/dmalch/go-geni/document"
	geniphoto "github.com/dmalch/go-geni/photo"
	geniprofile "github.com/dmalch/go-geni/profile"
	geniunion "github.com/dmalch/go-geni/union"
)

func TestFulfillDocumentRequests(t *testing.T) {
	t.Run("Missing ID in bulk response surfaces ErrResourceNotFound", func(t *testing.T) {
		RegisterTestingT(t)
		missing := asyncRequest[genidocument.Document]{
			Id:       "document-missing",
			Response: make(chan *genidocument.Document, 1),
			Error:    make(chan error, 1),
		}

		fulfillDocumentRequests([]asyncRequest[genidocument.Document]{missing}, nil)

		Expect(missing.Response).ToNot(Receive())
		var err error
		Expect(missing.Error).To(Receive(&err))
		Expect(errors.Is(err, geni.ErrResourceNotFound)).To(BeTrue())
	})

	t.Run("Mixed hits and misses route each request correctly", func(t *testing.T) {
		RegisterTestingT(t)
		hit := asyncRequest[genidocument.Document]{
			Id:       "document-hit",
			Response: make(chan *genidocument.Document, 1),
			Error:    make(chan error, 1),
		}
		miss := asyncRequest[genidocument.Document]{
			Id:       "document-miss",
			Response: make(chan *genidocument.Document, 1),
			Error:    make(chan error, 1),
		}

		fulfillDocumentRequests(
			[]asyncRequest[genidocument.Document]{hit, miss},
			[]genidocument.Document{{ID: "document-hit", Title: "found"}},
		)

		var hitResp *genidocument.Document
		Expect(hit.Response).To(Receive(&hitResp))
		Expect(hitResp.ID).To(Equal("document-hit"))
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
		missing := asyncRequest[geniunion.Union]{
			Id:       "union-missing",
			Response: make(chan *geniunion.Union, 1),
			Error:    make(chan error, 1),
		}

		fulfillUnionRequests([]asyncRequest[geniunion.Union]{missing}, nil)

		var err error
		Expect(missing.Error).To(Receive(&err))
		Expect(errors.Is(err, geni.ErrResourceNotFound)).To(BeTrue())
	})
}

func TestFulfillProfileRequests(t *testing.T) {
	t.Run("Missing ID in bulk response surfaces ErrResourceNotFound", func(t *testing.T) {
		RegisterTestingT(t)
		missing := asyncRequest[geniprofile.Profile]{
			Id:       "profile-missing",
			Response: make(chan *geniprofile.Profile, 1),
			Error:    make(chan error, 1),
		}

		fulfillProfileRequests([]asyncRequest[geniprofile.Profile]{missing}, nil)

		var err error
		Expect(missing.Error).To(Receive(&err))
		Expect(errors.Is(err, geni.ErrResourceNotFound)).To(BeTrue())
	})
}

func TestFulfillPhotoRequests(t *testing.T) {
	t.Run("Missing ID in bulk response surfaces ErrResourceNotFound", func(t *testing.T) {
		RegisterTestingT(t)
		missing := asyncRequest[geniphoto.Photo]{
			Id:       "photo-missing",
			Response: make(chan *geniphoto.Photo, 1),
			Error:    make(chan error, 1),
		}

		fulfillPhotoRequests([]asyncRequest[geniphoto.Photo]{missing}, nil)

		var err error
		Expect(missing.Error).To(Receive(&err))
		Expect(errors.Is(err, geni.ErrResourceNotFound)).To(BeTrue())
	})
}

// TestProcessBatchPanicRecovery verifies that a panic inside a batch worker
// goroutine is recovered and broadcast as an error to every request in the
// batch, rather than crashing the provider process and stranding every caller.
// A Client with a nil geni client makes the first API call panic, standing in
// for any unexpected panic inside a worker.
func TestProcessBatchPanicRecovery(t *testing.T) {
	t.Run("A panic in the union batch worker fails every request", func(t *testing.T) {
		RegisterTestingT(t)
		c := &Client{}
		req1 := asyncRequest[geniunion.Union]{
			Id:       "union-1",
			Response: make(chan *geniunion.Union, 1),
			Error:    make(chan error, 1),
		}
		req2 := asyncRequest[geniunion.Union]{
			Id:       "union-2",
			Response: make(chan *geniunion.Union, 1),
			Error:    make(chan error, 1),
		}

		c.processBatchOfUnions(t.Context(), []asyncRequest[geniunion.Union]{req1, req2})

		for _, req := range []asyncRequest[geniunion.Union]{req1, req2} {
			var err error
			Expect(req.Error).To(Receive(&err))
			Expect(err).To(MatchError(ContainSubstring("panic")))
			Expect(req.Response).ToNot(Receive())
		}
	})

	t.Run("A panic in the profile batch worker fails every request", func(t *testing.T) {
		RegisterTestingT(t)
		c := &Client{}
		req1 := asyncRequest[geniprofile.Profile]{
			Id:       "profile-1",
			Response: make(chan *geniprofile.Profile, 1),
			Error:    make(chan error, 1),
		}
		req2 := asyncRequest[geniprofile.Profile]{
			Id:       "profile-2",
			Response: make(chan *geniprofile.Profile, 1),
			Error:    make(chan error, 1),
		}

		c.processBatchOfProfiles(t.Context(), []asyncRequest[geniprofile.Profile]{req1, req2})

		for _, req := range []asyncRequest[geniprofile.Profile]{req1, req2} {
			var err error
			Expect(req.Error).To(Receive(&err))
			Expect(err).To(MatchError(ContainSubstring("panic")))
			Expect(req.Response).ToNot(Receive())
		}
	})

	t.Run("A panic in the document batch worker fails every request", func(t *testing.T) {
		RegisterTestingT(t)
		c := &Client{}
		req1 := asyncRequest[genidocument.Document]{
			Id:       "document-1",
			Response: make(chan *genidocument.Document, 1),
			Error:    make(chan error, 1),
		}
		req2 := asyncRequest[genidocument.Document]{
			Id:       "document-2",
			Response: make(chan *genidocument.Document, 1),
			Error:    make(chan error, 1),
		}

		c.processBatchOfDocuments(t.Context(), []asyncRequest[genidocument.Document]{req1, req2})

		for _, req := range []asyncRequest[genidocument.Document]{req1, req2} {
			var err error
			Expect(req.Error).To(Receive(&err))
			Expect(err).To(MatchError(ContainSubstring("panic")))
			Expect(req.Response).ToNot(Receive())
		}
	})

	t.Run("A panic in the photo batch worker fails every request", func(t *testing.T) {
		RegisterTestingT(t)
		c := &Client{}
		req1 := asyncRequest[geniphoto.Photo]{
			Id:       "photo-1",
			Response: make(chan *geniphoto.Photo, 1),
			Error:    make(chan error, 1),
		}
		req2 := asyncRequest[geniphoto.Photo]{
			Id:       "photo-2",
			Response: make(chan *geniphoto.Photo, 1),
			Error:    make(chan error, 1),
		}

		c.processBatchOfPhotos(t.Context(), []asyncRequest[geniphoto.Photo]{req1, req2})

		for _, req := range []asyncRequest[geniphoto.Photo]{req1, req2} {
			var err error
			Expect(req.Error).To(Receive(&err))
			Expect(err).To(MatchError(ContainSubstring("panic")))
			Expect(req.Response).ToNot(Receive())
		}
	})
}
