package geni

import (
	"net/http"
	"testing"

	. "github.com/onsi/gomega"
)

func TestWithModifier(t *testing.T) {
	t.Run("Sets relationship_modifier when value is non-empty", func(t *testing.T) {
		RegisterTestingT(t)
		req, _ := http.NewRequest(http.MethodPost, "https://example.com/api/profile-1/add-child", nil)

		WithModifier("foster")(req)

		Expect(req.URL.Query().Get("relationship_modifier")).To(Equal("foster"))
	})

	t.Run("Accepts adopt as a valid modifier", func(t *testing.T) {
		RegisterTestingT(t)
		req, _ := http.NewRequest(http.MethodPost, "https://example.com/api/profile-1/add-child", nil)

		WithModifier("adopt")(req)

		Expect(req.URL.Query().Get("relationship_modifier")).To(Equal("adopt"))
	})

	t.Run("Omits the query param when value is empty", func(t *testing.T) {
		RegisterTestingT(t)
		req, _ := http.NewRequest(http.MethodPost, "https://example.com/api/profile-1/add-child", nil)

		WithModifier("")(req)

		Expect(req.URL.Query().Has("relationship_modifier")).To(BeFalse())
	})

	t.Run("Preserves any pre-existing query params", func(t *testing.T) {
		RegisterTestingT(t)
		req, _ := http.NewRequest(http.MethodPost, "https://example.com/api/profile-1/add-child?fields=id", nil)

		WithModifier("foster")(req)

		Expect(req.URL.Query().Get("relationship_modifier")).To(Equal("foster"))
		Expect(req.URL.Query().Get("fields")).To(Equal("id"))
	})
}
