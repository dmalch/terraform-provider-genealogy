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
	"sync"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"golang.org/x/oauth2"
	"golang.org/x/time/rate"
)

var ErrResourceNotFound = fmt.Errorf("resource not found")
var ErrAccessDenied = fmt.Errorf("access denied")

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
	client        *http.Client
	limiter       *rate.Limiter
	urlMap        *sync.Map
}

func NewClient(tokenSource oauth2.TokenSource, useSandboxEnv bool) *Client {
	return &Client{
		useSandboxEnv: useSandboxEnv,
		tokenSource:   tokenSource,
		client:        &http.Client{},
		limiter:       rate.NewLimiter(rate.Every(1*time.Second), 1),
		urlMap:        &sync.Map{},
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

type opt struct {
	getRequestKey          func() string
	prepareBulkRequestFrom func(*http.Request, *sync.Map)
	parseBulkResponse      func(*http.Request, []byte, *sync.Map) ([]byte, error)
}

func WithRequestKey(fn func() string) func(*opt) {
	return func(o *opt) {
		o.getRequestKey = fn
	}
}

func WithPrepareBulkRequest(fn func(*http.Request, *sync.Map)) func(*opt) {
	return func(o *opt) {
		o.prepareBulkRequestFrom = fn
	}
}

func WithParseBulkResponse(fn func(*http.Request, []byte, *sync.Map) ([]byte, error)) func(*opt) {
	return func(o *opt) {
		o.parseBulkResponse = fn
	}
}

func (c *Client) doRequest(ctx context.Context, req *http.Request, opts ...func(*opt)) ([]byte, error) {
	// Initialize the opt struct with default no-op functions
	options := opt{
		prepareBulkRequestFrom: func(*http.Request, *sync.Map) {},
	}

	// Apply the provided opts to the options struct
	for _, o := range opts {
		o(&options)
	}

	if err := c.addStandardHeadersAndQueryParams(req); err != nil {
		return nil, err
	}

	// Retry logic using retry-go
	return retry.DoWithData(
		func() ([]byte, error) {
			limiterCtx, limiterCtxCancelFunc := context.WithCancel(ctx)
			defer limiterCtxCancelFunc()

			if options.getRequestKey != nil {
				// Store the key in the map
				c.urlMap.Store(options.getRequestKey(), limiterCtxCancelFunc)
			}

			if err := c.limiter.Wait(limiterCtx); err != nil {
				// If the context is canceled, we should not return an error
				if !errors.Is(err, context.Canceled) {
					tflog.Error(ctx, "Error waiting for rate limiter", map[string]interface{}{"error": err})
					return nil, err
				}
			}

			// Check if the response is already cached
			if options.getRequestKey != nil {
				if cachedRes, ok := c.urlMap.LoadAndDelete(options.getRequestKey()); ok && cachedRes != nil {
					if res, ok := cachedRes.([]byte); ok {
						tflog.Debug(ctx, "Using cached response")
						return res, nil
					}
				}

				options.prepareBulkRequestFrom(req, c.urlMap)
			}

			tflog.Debug(ctx, "Sending request", map[string]interface{}{"method": req.Method, "url": req.URL.String()})
			res, err := c.client.Do(req)
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

			apiRateWindow := res.Header.Get("X-API-Rate-Window")
			apiRateLimit := res.Header.Get("X-API-Rate-Limit")

			tflog.Debug(ctx, "Received response", map[string]interface{}{"status": res.StatusCode, "X-API-Rate-Window": apiRateWindow, "X-API-Rate-Limit": apiRateLimit})
			tflog.Trace(ctx, "Received response", map[string]interface{}{"status": res.StatusCode, "body": string(body), "X-API-Rate-Window": apiRateWindow, "X-API-Rate-Limit": apiRateLimit})

			secondsUntilRetry, err := strconv.Atoi(apiRateWindow)
			if err == nil {
				if apiRateLimitNumber, err := strconv.Atoi(apiRateLimit); err == nil {
					newLimit := rate.Every(time.Duration(secondsUntilRetry) * time.Second / time.Duration(apiRateLimitNumber))
					if c.limiter.Limit() != newLimit {
						tflog.Debug(ctx, "Setting rate limit", map[string]interface{}{"limit": newLimit, "seconds_until_retry": secondsUntilRetry, "api_rate_limit": apiRateLimit})
						c.limiter.SetLimit(newLimit)
					}
					if c.limiter.Burst() != apiRateLimitNumber {
						tflog.Debug(ctx, "Setting rate burst", map[string]interface{}{"burst": apiRateLimitNumber})
						c.limiter.SetBurst(apiRateLimitNumber)
					}
				}
			}

			if res.StatusCode != http.StatusOK {
				if res.StatusCode == http.StatusTooManyRequests {
					tflog.Warn(ctx, "Received 429 Too Many Requests, retrying...", map[string]interface{}{"X-API-Rate-Window": secondsUntilRetry})
					return nil, newErrWithRetry(res.StatusCode, secondsUntilRetry)
				}

				if res.StatusCode == http.StatusUnauthorized {
					tflog.Warn(ctx, "Received 401 Unauthorized, retrying...")
					return nil, newErrWithRetry(res.StatusCode, 1)
				}

				if res.StatusCode == http.StatusForbidden {
					tflog.Warn(ctx, "Received 403 Forbidden.")
					return nil, ErrAccessDenied
				}

				if res.StatusCode == http.StatusNotFound {
					tflog.Warn(ctx, "Received 404 Not Found.")
					return nil, ErrResourceNotFound
				}

				if strings.Contains(string(body), "Request unsuccessful. Incapsula incident ID:") {
					// Incapsula is a DDoS protection service that Geni uses. If we get a response
					// with this message, it means that the request was blocked by Incapsula.
					tflog.Warn(ctx, "Incapsula blocked request.")
					return nil, fmt.Errorf("incapsula blocked request")
				}

				tflog.Error(ctx, "Non-OK HTTP status", map[string]interface{}{"status": res.StatusCode, "body": string(body)})
				return nil, fmt.Errorf("non-OK HTTP status: %s, body: %s", res.Status, string(body))
			}

			if options.parseBulkResponse != nil {
				return options.parseBulkResponse(req, body, c.urlMap)
			}

			return body, nil
		},
		retry.RetryIf(func(err error) bool {
			var errCode429WithRetry errCode429WithRetry
			return errors.As(err, &errCode429WithRetry)
		}),
		retry.Context(ctx),
		retry.Attempts(4),
		retry.Delay(2*time.Second),     // Wait 2 seconds between retries
		retry.MaxJitter(2*time.Second), // Add up to 2 seconds of jitter to each retry
		retry.DelayType(retry.CombineDelay(retry.FixedDelay, retry.RandomDelay)),
		retry.OnRetry(func(n uint, err error) {
			tflog.Debug(ctx, "Retrying request", map[string]interface{}{"attempt": n + 1, "error": err})
		}),
	)
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
	req.Header.Add("User-Agent", "terraform-provider-genealogy/0.1")

	return nil
}
