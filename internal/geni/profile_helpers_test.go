package geni

import (
	"testing"

	. "github.com/onsi/gomega"
	"golang.org/x/oauth2"
)

func TestFixResponse(t *testing.T) {
	t.Run("Strips production API URL prefix from union URLs", func(t *testing.T) {
		RegisterTestingT(t)
		client := NewClient(oauth2.StaticTokenSource(&oauth2.Token{AccessToken: "test"}), false)

		profile := &ProfileResponse{
			Unions: []string{
				"https://www.geni.com/api/union-123",
				"https://www.geni.com/api/union-456",
			},
		}

		client.fixResponse(profile)

		Expect(profile.Unions).To(Equal([]string{"union-123", "union-456"}))
	})

	t.Run("Strips sandbox API URL prefix from union URLs", func(t *testing.T) {
		RegisterTestingT(t)
		client := NewClient(oauth2.StaticTokenSource(&oauth2.Token{AccessToken: "test"}), true)

		profile := &ProfileResponse{
			Unions: []string{
				"https://api.sandbox.geni.com/union-789",
			},
		}

		client.fixResponse(profile)

		Expect(profile.Unions).To(Equal([]string{"union-789"}))
	})

	t.Run("Handles empty unions", func(t *testing.T) {
		RegisterTestingT(t)
		client := NewClient(oauth2.StaticTokenSource(&oauth2.Token{AccessToken: "test"}), false)

		profile := &ProfileResponse{
			Unions: []string{},
		}

		client.fixResponse(profile)

		Expect(profile.Unions).To(BeEmpty())
	})
}

func TestEscapeString(t *testing.T) {
	t.Run("Delegates to escapeStringToUTF", func(t *testing.T) {
		RegisterTestingT(t)
		Expect(escapeString("café")).To(Equal(escapeStringToUTF("café")))
	})

	t.Run("ASCII passthrough", func(t *testing.T) {
		RegisterTestingT(t)
		Expect(escapeString("Hello World")).To(Equal("Hello World"))
	})

	t.Run("Non-ASCII characters are escaped", func(t *testing.T) {
		RegisterTestingT(t)
		result := escapeString("Ф")
		Expect(result).To(Equal("\\u0424"))
	})
}
