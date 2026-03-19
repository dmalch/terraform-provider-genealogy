package authn

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/labstack/echo/v4"
	. "github.com/onsi/gomega"
)

func TestCallbackHandle(t *testing.T) {
	t.Run("Accepts token when state matches", func(t *testing.T) {
		RegisterTestingT(t)

		handler := &callback{
			expectedState: "valid-state",
			shutdownCh:    make(chan error, 1),
		}

		e := echo.New()
		q := make(url.Values)
		q.Set("state", "valid-state")
		q.Set("access_token", "my-token")
		q.Set("expires_in", "3600")
		req := httptest.NewRequest(http.MethodGet, "/callback?"+q.Encode(), nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.handle(c)
		Expect(err).ToNot(HaveOccurred())

		shutdownErr := <-handler.shutdownCh
		Expect(shutdownErr).ToNot(HaveOccurred())
		Expect(handler.accessToken).To(Equal("my-token"))
		Expect(handler.expiresIn).To(Equal("3600"))
		Expect(rec.Body.String()).To(ContainSubstring("Login was successful"))
	})

	t.Run("Rejects callback when state does not match", func(t *testing.T) {
		RegisterTestingT(t)

		handler := &callback{
			expectedState: "valid-state",
			shutdownCh:    make(chan error, 1),
		}

		e := echo.New()
		q := make(url.Values)
		q.Set("state", "wrong-state")
		q.Set("access_token", "my-token")
		req := httptest.NewRequest(http.MethodGet, "/callback?"+q.Encode(), nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.handle(c)
		Expect(err).ToNot(HaveOccurred())

		shutdownErr := <-handler.shutdownCh
		Expect(shutdownErr).To(HaveOccurred())
		Expect(shutdownErr.Error()).To(ContainSubstring("state parameter mismatch"))
		Expect(handler.accessToken).To(BeEmpty())
		Expect(rec.Body.String()).To(ContainSubstring("CSRF"))
	})

	t.Run("Rejects callback when state is missing", func(t *testing.T) {
		RegisterTestingT(t)

		handler := &callback{
			expectedState: "valid-state",
			shutdownCh:    make(chan error, 1),
		}

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/callback?access_token=my-token", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.handle(c)
		Expect(err).ToNot(HaveOccurred())

		shutdownErr := <-handler.shutdownCh
		Expect(shutdownErr).To(HaveOccurred())
		Expect(handler.accessToken).To(BeEmpty())
	})

	t.Run("Reports unsuccessful login when token is empty but state matches", func(t *testing.T) {
		RegisterTestingT(t)

		handler := &callback{
			expectedState: "valid-state",
			shutdownCh:    make(chan error, 1),
		}

		e := echo.New()
		q := make(url.Values)
		q.Set("state", "valid-state")
		req := httptest.NewRequest(http.MethodGet, "/callback?"+q.Encode(), nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.handle(c)
		Expect(err).ToNot(HaveOccurred())

		shutdownErr := <-handler.shutdownCh
		Expect(shutdownErr).ToNot(HaveOccurred())
		Expect(handler.accessToken).To(BeEmpty())
		Expect(rec.Body.String()).To(ContainSubstring("Login was not successful"))
	})
}

func TestNewAuthTokenSource(t *testing.T) {
	t.Run("Creates token source with config", func(t *testing.T) {
		RegisterTestingT(t)

		src := NewAuthTokenSource(nil)
		Expect(src).ToNot(BeNil())
		Expect(src.config).To(BeNil())
	})
}
