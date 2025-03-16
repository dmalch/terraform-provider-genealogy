package geni

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/avast/retry-go/v4"
)

var errCode429 = errors.New("received 429 status")

func doRequest(req *http.Request) ([]byte, error) {
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
			defer res.Body.Close()

			body, err = io.ReadAll(res.Body)
			if err != nil {
				slog.Error("Error reading response", "error", err)
				return err
			}

			if res.StatusCode == http.StatusTooManyRequests {
				slog.Warn("Received 429 Too Many Requests, retrying...")
				return errCode429
			}

			if res.StatusCode != http.StatusOK {
				slog.Error("Non-OK HTTP status", "status", res.StatusCode, "body", string(body))
				return fmt.Errorf("non-OK HTTP status: %s", res.Status)
			}

			return nil
		},
		retry.RetryIf(func(err error) bool {
			if errors.Is(err, errCode429) {
				return true
			}
			return false
		}),
		retry.Attempts(3),
		retry.Delay(2*time.Second),        // Wait 2 seconds between retries
		retry.DelayType(retry.FixedDelay), // Use a fixed delay between retries
		retry.OnRetry(func(n uint, err error) {
			slog.Info("Retrying request", "attempt", n+1, "error", err)
		}),
	)

	if err != nil {
		return nil, err
	}

	return body, nil
}

func addStandardHeadersAndQueryParams(req *http.Request, accessToken string) {
	query := req.URL.Query()
	query.Add("access_token", accessToken)
	query.Add("api_version", apiVersion)
	// The returned data structures will contain urls to other objects by default,
	// unless the request includes 'only_ids=true.' Passing only_ids will force the
	// system to return ids only.
	query.Add("only_ids", "true")

	req.URL.RawQuery = query.Encode()
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
}
