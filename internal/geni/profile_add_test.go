package geni

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"

	. "github.com/onsi/gomega"
	"golang.org/x/oauth2"
)

type captureTransport struct {
	lastRequest *http.Request
}

func (t *captureTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	t.lastRequest = req.Clone(req.Context())
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewBufferString(`{"id":"profile-tmp"}`)),
		Header:     make(http.Header),
	}, nil
}

func newCapturingClient() (*Client, *captureTransport) {
	ct := &captureTransport{}
	c := NewClient(oauth2.StaticTokenSource(&oauth2.Token{AccessToken: "test-token"}), true)
	c.client = &http.Client{Transport: ct}
	return c, ct
}

func TestAddChild_SendsRelationshipModifier(t *testing.T) {
	t.Run("With foster modifier sets query param", func(t *testing.T) {
		RegisterTestingT(t)
		c, ct := newCapturingClient()

		_, err := c.AddChild(context.Background(), "profile-1", WithModifier("foster"))

		Expect(err).ToNot(HaveOccurred())
		Expect(ct.lastRequest).ToNot(BeNil())
		Expect(ct.lastRequest.URL.Query().Get("relationship_modifier")).To(Equal("foster"))
	})

	t.Run("With adopt modifier sets query param", func(t *testing.T) {
		RegisterTestingT(t)
		c, ct := newCapturingClient()

		_, err := c.AddChild(context.Background(), "profile-1", WithModifier("adopt"))

		Expect(err).ToNot(HaveOccurred())
		Expect(ct.lastRequest.URL.Query().Get("relationship_modifier")).To(Equal("adopt"))
	})

	t.Run("Without options omits the query param", func(t *testing.T) {
		RegisterTestingT(t)
		c, ct := newCapturingClient()

		_, err := c.AddChild(context.Background(), "profile-1")

		Expect(err).ToNot(HaveOccurred())
		Expect(ct.lastRequest.URL.Query().Has("relationship_modifier")).To(BeFalse())
	})

	t.Run("Targets the /add-child path for the given id", func(t *testing.T) {
		RegisterTestingT(t)
		c, ct := newCapturingClient()

		_, err := c.AddChild(context.Background(), "union-42", WithModifier("foster"))

		Expect(err).ToNot(HaveOccurred())
		Expect(ct.lastRequest.URL.Path).To(HaveSuffix("/api/union-42/add-child"))
	})
}

func TestAddSibling_SendsRelationshipModifier(t *testing.T) {
	t.Run("With foster modifier sets query param", func(t *testing.T) {
		RegisterTestingT(t)
		c, ct := newCapturingClient()

		_, err := c.AddSibling(context.Background(), "profile-1", WithModifier("foster"))

		Expect(err).ToNot(HaveOccurred())
		Expect(ct.lastRequest.URL.Query().Get("relationship_modifier")).To(Equal("foster"))
	})

	t.Run("Without options omits the query param", func(t *testing.T) {
		RegisterTestingT(t)
		c, ct := newCapturingClient()

		_, err := c.AddSibling(context.Background(), "profile-1")

		Expect(err).ToNot(HaveOccurred())
		Expect(ct.lastRequest.URL.Query().Has("relationship_modifier")).To(BeFalse())
	})
}
