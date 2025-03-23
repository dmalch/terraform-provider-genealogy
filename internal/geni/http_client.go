package geni

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/avast/retry-go/v4"
)

type errCode429WithRetry struct {
	statusCode        int
	secondsUntilRetry int
}

func (e errCode429WithRetry) Error() string {
	return fmt.Sprintf("received %d status, window is %d seconds", e.statusCode, e.secondsUntilRetry)
}

func newErrWithRetry(statusCode int, secondsUntilRetry int) error {
	return errCode429WithRetry{
		statusCode:        statusCode,
		secondsUntilRetry: secondsUntilRetry,
	}
}

type Client struct {
	accessToken   string
	useSandboxEnv bool
}

func NewClient(accessToken string, useSandboxEnv bool) *Client {
	return &Client{
		accessToken:   accessToken,
		useSandboxEnv: useSandboxEnv,
	}
}

func (c *Client) getBaseUrl() string {
	if c.useSandboxEnv {
		return geniSandboxUrl
	}
	return geniProdUrl
}

func (c *Client) getApiUrl() string {
	if c.useSandboxEnv {
		return geniSandboxApiUrl
	}
	return geniProdApiUrl
}

func (c *Client) doRequest(req *http.Request) ([]byte, error) {
	var body []byte
	var err error

	client := &http.Client{}

	// Retry logic using retry-go
	err = retry.Do(
		func() error {
			res, err := client.Do(req)
			if err != nil {
				slog.Error("Error sending request", "error", err)
				return err
			}
			defer func(Body io.ReadCloser) {
				_ = Body.Close()
			}(res.Body)

			body, err = io.ReadAll(res.Body)
			if err != nil {
				slog.Error("Error reading response", "error", err)
				return err
			}

			if res.StatusCode != http.StatusOK {
				if res.StatusCode == http.StatusTooManyRequests {
					slog.Warn("Received 429 Too Many Requests, retrying...")
					apiRateWindow := res.Header.Get("X-API-Rate-Window")
					secondsUntilRetry, err := strconv.Atoi(apiRateWindow)
					if err != nil {
						return fmt.Errorf("invalid value for X-API-Rate-Window: %d", secondsUntilRetry)
					}

					return newErrWithRetry(res.StatusCode, secondsUntilRetry)
				}

				if strings.Contains(string(body), "Request unsuccessful. Incapsula incident ID:") {
					slog.Warn("Non-OK HTTP status", "status", res.StatusCode, "body", string(body))
					return newErrWithRetry(res.StatusCode, 1)
				}

				return fmt.Errorf("non-OK HTTP status: %s", res.Status)
			}

			return nil
		},
		retry.RetryIf(func(err error) bool {
			var errCode429WithRetry errCode429WithRetry
			return errors.As(err, &errCode429WithRetry)
		}),
		retry.Attempts(5),
		retry.Delay(2*time.Second), // Wait 2 seconds between retries
		retry.DelayType(rateLimitingDelay),
		retry.OnRetry(func(n uint, err error) {
			slog.Info("Retrying request", "attempt", n+1, "error", err)
		}),
	)

	if err != nil {
		return nil, err
	}

	return body, nil
}

func rateLimitingDelay(n uint, err error, config *retry.Config) time.Duration {
	var retryErr errCode429WithRetry
	if errors.As(err, &retryErr) {
		return time.Duration(retryErr.secondsUntilRetry+1) * time.Second
	}
	return retry.FixedDelay(n, err, config)
}

func (c *Client) addStandardHeadersAndQueryParams(req *http.Request) {
	query := req.URL.Query()
	query.Add("access_token", c.accessToken)
	query.Add("api_version", apiVersion)
	// The returned data structures will contain urls to other objects by default,
	// unless the request includes 'only_ids=true.' Passing only_ids will force the
	// system to return ids only.
	query.Add("only_ids", "true")

	req.URL.RawQuery = query.Encode()
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
}
