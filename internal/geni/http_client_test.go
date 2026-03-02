package geni

import (
	"net/url"
	"testing"

	. "github.com/onsi/gomega"
)

func TestBaseUrl(t *testing.T) {
	t.Run("Returns production URL when not using sandbox", func(t *testing.T) {
		RegisterTestingT(t)
		Expect(BaseUrl(false)).To(Equal("https://www.geni.com/"))
	})

	t.Run("Returns sandbox URL when using sandbox", func(t *testing.T) {
		RegisterTestingT(t)
		Expect(BaseUrl(true)).To(Equal("https://sandbox.geni.com/"))
	})
}

func TestApiUrl(t *testing.T) {
	t.Run("Returns production API URL when not using sandbox", func(t *testing.T) {
		RegisterTestingT(t)
		Expect(apiUrl(false)).To(Equal("https://www.geni.com/api/"))
	})

	t.Run("Returns sandbox API URL when using sandbox", func(t *testing.T) {
		RegisterTestingT(t)
		Expect(apiUrl(true)).To(Equal("https://api.sandbox.geni.com/"))
	})
}

func TestRedactURL(t *testing.T) {
	t.Run("Redacts access_token from URL", func(t *testing.T) {
		RegisterTestingT(t)
		u, _ := url.Parse("https://www.geni.com/api/profile?access_token=secret123&fields=id")

		result := redactURL(u)

		Expect(result).To(ContainSubstring("access_token=REDACTED"))
		Expect(result).ToNot(ContainSubstring("secret123"))
		Expect(result).To(ContainSubstring("fields=id"))
	})

	t.Run("Returns URL unchanged when no access_token", func(t *testing.T) {
		RegisterTestingT(t)
		u, _ := url.Parse("https://www.geni.com/api/profile?fields=id")

		result := redactURL(u)

		Expect(result).To(Equal("https://www.geni.com/api/profile?fields=id"))
	})

	t.Run("Preserves other query params", func(t *testing.T) {
		RegisterTestingT(t)
		u, _ := url.Parse("https://www.geni.com/api/profile?access_token=secret&api_version=1&only_ids=true")

		result := redactURL(u)

		Expect(result).To(ContainSubstring("access_token=REDACTED"))
		Expect(result).To(ContainSubstring("api_version=1"))
		Expect(result).To(ContainSubstring("only_ids=true"))
	})
}

func TestErrCode429WithRetry(t *testing.T) {
	t.Run("Formats error message correctly", func(t *testing.T) {
		RegisterTestingT(t)
		err := errCode429WithRetry{
			statusCode:        429,
			secondsUntilRetry: 30,
		}

		Expect(err.Error()).To(Equal("received 429 status, retry in 30 seconds"))
	})

	t.Run("Works with different status codes", func(t *testing.T) {
		RegisterTestingT(t)
		err := errCode429WithRetry{
			statusCode:        401,
			secondsUntilRetry: 1,
		}

		Expect(err.Error()).To(Equal("received 401 status, retry in 1 seconds"))
	})
}

func TestEscapeStringToUTF(t *testing.T) {
	t.Run("ASCII characters pass through unchanged", func(t *testing.T) {
		RegisterTestingT(t)
		Expect(escapeStringToUTF("Hello World")).To(Equal("Hello World"))
	})

	t.Run("Non-ASCII characters are escaped", func(t *testing.T) {
		RegisterTestingT(t)
		Expect(escapeStringToUTF("café")).To(Equal("caf\\u00e9"))
	})

	t.Run("Empty string returns empty", func(t *testing.T) {
		RegisterTestingT(t)
		Expect(escapeStringToUTF("")).To(Equal(""))
	})

	t.Run("Mixed ASCII and non-ASCII content", func(t *testing.T) {
		RegisterTestingT(t)
		result := escapeStringToUTF("Hello Wörld")
		Expect(result).To(Equal("Hello W\\u00f6rld"))
	})

	t.Run("Cyrillic characters are escaped", func(t *testing.T) {
		RegisterTestingT(t)
		result := escapeStringToUTF("Привет")
		Expect(result).To(Equal("\\u041f\\u0440\\u0438\\u0432\\u0435\\u0442"))
	})
}
