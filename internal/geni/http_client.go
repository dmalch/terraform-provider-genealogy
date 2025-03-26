package geni

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"golang.org/x/oauth2"
)

type errCode429WithRetry struct {
	statusCode        int
	secondsUntilRetry int
}

func (e errCode429WithRetry) Error() string {
	return fmt.Sprintf("received %d status, retry in %d seconds", e.statusCode, e.secondsUntilRetry)
}

func newErrWithRetry(statusCode int, secondsUntilRetry int) error {
	return errCode429WithRetry{
		statusCode:        statusCode,
		secondsUntilRetry: secondsUntilRetry,
	}
}

type Client struct {
	useSandboxEnv bool
	tokenSource   oauth2.TokenSource
}

func NewClient(tokenSource oauth2.TokenSource, useSandboxEnv bool) *Client {
	return &Client{
		useSandboxEnv: useSandboxEnv,
		tokenSource:   tokenSource,
	}
}

func BaseUrl(useSandboxEnv bool) string {
	if useSandboxEnv {
		return geniSandboxUrl
	}
	return geniProdUrl
}

func apiUrl(useSandboxEnv bool) string {
	if useSandboxEnv {
		return geniSandboxApiUrl
	}
	return geniProdApiUrl
}

func (c *Client) doRequest(ctx context.Context, req *http.Request) ([]byte, error) {
	client := &http.Client{}

	// Retry logic using retry-go
	return retry.DoWithData(
		func() ([]byte, error) {
			tflog.Debug(ctx, "Sending request", map[string]interface{}{"method": req.Method, "url": req.URL.String()})
			res, err := client.Do(req)
			if err != nil {
				slog.Error("Error sending request", "error", err)
				return nil, err
			}
			defer func(Body io.ReadCloser) {
				_ = Body.Close()
			}(res.Body)

			body, err := io.ReadAll(res.Body)
			if err != nil {
				tflog.Error(ctx, "Error reading response", map[string]interface{}{"error": err})
				return nil, err
			}

			if res.StatusCode != http.StatusOK {
				if res.StatusCode == http.StatusTooManyRequests {
					apiRateWindow := res.Header.Get("X-API-Rate-Window")
					secondsUntilRetry, err := strconv.Atoi(apiRateWindow)
					if err != nil {
						return nil, fmt.Errorf("invalid value for X-API-Rate-Window: %d", secondsUntilRetry)
					}

					tflog.Warn(ctx, "Received 429 Too Many Requests, retrying...", map[string]interface{}{"X-API-Rate-Window": secondsUntilRetry})
					return nil, newErrWithRetry(res.StatusCode, secondsUntilRetry)
				}

				if res.StatusCode == http.StatusUnauthorized {
					tflog.Warn(ctx, "Received 401 Unauthorized, retrying.")
					return nil, newErrWithRetry(res.StatusCode, 1)
				}

				if strings.Contains(string(body), "Request unsuccessful. Incapsula incident ID:") {
					tflog.Error(ctx, "Non-OK HTTP status", map[string]interface{}{"status": res.StatusCode, "body": string(body)})
					return nil, fmt.Errorf("non-OK HTTP status: %s", res.Status)
				}

				tflog.Error(ctx, "Non-OK HTTP status", map[string]interface{}{"status": res.StatusCode, "body": string(body)})
				return nil, fmt.Errorf("non-OK HTTP status: %s, body: %s", res.Status, string(body))
			}

			tflog.Debug(ctx, "Received response", map[string]interface{}{"status": res.StatusCode})
			tflog.Trace(ctx, "Received response", map[string]interface{}{"status": res.StatusCode, "body": string(body)})
			return body, nil
		},
		retry.RetryIf(func(err error) bool {
			var errCode429WithRetry errCode429WithRetry
			return errors.As(err, &errCode429WithRetry)
		}),
		retry.Attempts(5),
		retry.Delay(2*time.Second), // Wait 2 seconds between retries
		retry.DelayType(rateLimitingDelay),
		retry.OnRetry(func(n uint, err error) {
			tflog.Debug(ctx, "Retrying request", map[string]interface{}{"attempt": n + 1, "error": err})
		}),
	)
}

func rateLimitingDelay(n uint, err error, config *retry.Config) time.Duration {
	var retryErr errCode429WithRetry
	if errors.As(err, &retryErr) {
		return time.Duration(retryErr.secondsUntilRetry+1) * time.Second
	}
	return retry.FixedDelay(n, err, config)
}

func (c *Client) addStandardHeadersAndQueryParams(req *http.Request) error {
	query := req.URL.Query()

	token, err := c.tokenSource.Token()
	if err != nil {
		return fmt.Errorf("error getting token: %w", err)
	}

	query.Add("access_token", token.AccessToken)
	query.Add("api_version", apiVersion)
	// The returned data structures will contain urls to other objects by default,
	// unless the request includes 'only_ids=true.' Passing only_ids will force the
	// system to return ids only.
	query.Add("only_ids", "true")

	req.URL.RawQuery = query.Encode()
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	return nil
}
