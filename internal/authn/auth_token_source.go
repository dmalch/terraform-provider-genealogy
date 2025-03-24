package authn

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/skratchdot/open-golang/open"
	"golang.org/x/oauth2"
)

// authTokenSource implements oauth2.TokenSource.
type authTokenSource struct {
	config *oauth2.Config
}

func NewAuthTokenSource(config *oauth2.Config) *authTokenSource {
	return &authTokenSource{
		config: config,
	}
}

// Token retrieves a new token, performing the OAuth flow if necessary.
func (a *authTokenSource) Token() (*oauth2.Token, error) {
	sigIntCh := make(chan os.Signal, 1)
	signal.Notify(sigIntCh, os.Interrupt)
	defer signal.Stop(sigIntCh)

	// Start a local server to handle the callback
	callbackHandler := &callback{
		shutdownCh: make(chan error),
	}

	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.GET("/callback", callbackHandler.handle)
	go func() {
		err := e.Start(":8080")
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			callbackHandler.shutdownCh <- fmt.Errorf("failed to start Echo server: %w", err)
		}
	}()

	authURL := a.config.AuthCodeURL("", oauth2.AccessTypeOffline,
		oauth2.SetAuthURLParam("response_type", "token"),
		oauth2.SetAuthURLParam("display", "mobile"),
	)

	// Open the URL in the default browser
	err := open.Start(authURL)
	if err != nil {
		return nil, err
	}

	// Wait for the token to be set by the callback handler, for SIGINT to be
	// received, or for up to 5 minutes.
	select {
	case err := <-callbackHandler.shutdownCh:
		if err != nil {
			return nil, err
		}

		err = e.Shutdown(context.Background())
		if err != nil {
			return nil, fmt.Errorf("failed to shutdown the server: %w", err)
		}

		if callbackHandler.accessToken == "" {
			return nil, errors.New("no authentication access token was received")
		}

		expiresIn, err := strconv.ParseInt(callbackHandler.expiresIn, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse expires_in: %w", err)
		}

		return &oauth2.Token{
			AccessToken: callbackHandler.accessToken,
			ExpiresIn:   expiresIn,
			Expiry:      time.Now().Add(time.Duration(expiresIn) * time.Second),
		}, nil
	case <-sigIntCh:
		return nil, errors.New("interrupted")
	case <-time.After(5 * time.Minute):
		return nil, errors.New("timed out while waiting for a response")
	}
}

type callback struct {
	accessToken string
	expiresIn   string
	shutdownCh  chan error
}

func (handler *callback) handle(c echo.Context) error {
	accessToken := c.QueryParam("access_token")
	if accessToken != "" {
		handler.accessToken = accessToken
		handler.expiresIn = c.QueryParam("expires_in")
		_, _ = fmt.Fprintln(c.Response().Writer, "Login was successful. You can close the browser and return to the command line.")
	} else {
		_, _ = fmt.Fprintln(c.Response().Writer, "Login was not successful. You can close the browser and try again.")
	}
	handler.shutdownCh <- nil
	return nil
}
