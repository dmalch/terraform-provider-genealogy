package internal

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	. "github.com/onsi/gomega"
	"golang.org/x/oauth2"

	"github.com/dmalch/terraform-provider-genealogy/internal/config"
	"github.com/dmalch/terraform-provider-genealogy/internal/geniapp"
)

func TestTokenCacheFilePath(t *testing.T) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("failed to get home dir: %v", err)
	}

	t.Run("production environment", func(t *testing.T) {
		RegisterTestingT(t)

		result, err := tokenCacheFilePath(false)

		Expect(err).ToNot(HaveOccurred())
		Expect(result).To(Equal(path.Join(homeDir, ".genealogy", "geni_token.json")))
	})

	t.Run("sandbox environment", func(t *testing.T) {
		RegisterTestingT(t)

		result, err := tokenCacheFilePath(true)

		Expect(err).ToNot(HaveOccurred())
		Expect(result).To(Equal(path.Join(homeDir, ".genealogy", "geni_sandbox_token.json")))
	})
}

func TestConfigure(t *testing.T) {
	t.Run("populates client data when a static access token is provided", func(t *testing.T) {
		RegisterTestingT(t)
		p := newProvider(t)

		resp := configureProvider(t, p, "test-token")

		Expect(resp.Diagnostics.HasError()).To(BeFalse())
		Expect(p.client).ToNot(BeNil())
		Expect(p.batchClient).ToNot(BeNil())

		data, ok := resp.ResourceData.(*config.ClientData)
		Expect(ok).To(BeTrue())
		Expect(data.Client).To(BeIdenticalTo(p.client))
		Expect(data.BatchClient).To(BeIdenticalTo(p.batchClient))
	})

	t.Run("separate provider instances get independent clients", func(t *testing.T) {
		RegisterTestingT(t)
		p1 := newProvider(t)
		p2 := newProvider(t)

		Expect(configureProvider(t, p1, "token-1").Diagnostics.HasError()).To(BeFalse())
		Expect(configureProvider(t, p2, "token-2").Diagnostics.HasError()).To(BeFalse())

		Expect(p1.client).ToNot(BeNil())
		Expect(p2.client).ToNot(BeNil())
		Expect(p1.client).ToNot(BeIdenticalTo(p2.client))
	})
}

// newProvider returns a fresh *GeniProvider, failing the test if New ever
// returns a different concrete type.
func newProvider(t *testing.T) *GeniProvider {
	t.Helper()
	p, ok := New().(*GeniProvider)
	if !ok {
		t.Fatal("New() did not return a *GeniProvider")
	}
	return p
}

// configureProvider drives p.Configure with a provider config carrying only the
// given access token (every other attribute null), so the run stays offline:
// a static token skips the OAuth flow entirely.
func configureProvider(t *testing.T, p *GeniProvider, accessToken string) *provider.ConfigureResponse {
	t.Helper()
	ctx := t.Context()

	var schemaResp provider.SchemaResponse
	p.Schema(ctx, provider.SchemaRequest{}, &schemaResp)

	raw := tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
		"access_token":                tftypes.NewValue(tftypes.String, accessToken),
		"client_id":                   tftypes.NewValue(tftypes.String, nil),
		"client_secret":               tftypes.NewValue(tftypes.String, nil),
		"use_sandbox_env":             tftypes.NewValue(tftypes.Bool, nil),
		"auto_update_merged_profiles": tftypes.NewValue(tftypes.Bool, nil),
	})

	resp := &provider.ConfigureResponse{}
	p.Configure(ctx, provider.ConfigureRequest{
		Config: tfsdk.Config{Schema: schemaResp.Schema, Raw: raw},
	}, resp)
	return resp
}

// TestBrowserTokenSource covers the reason this provider learned about
// client secrets: a token that renews itself instead of opening a
// browser every 24 hours.
func TestBrowserTokenSource(t *testing.T) {
	// seedCache writes a token into a fresh cache file and returns its
	// path.
	seedCache := func(t *testing.T, token *oauth2.Token) string {
		t.Helper()

		cachePath := path.Join(t.TempDir(), "geni_token.json")
		body, err := json.Marshal(token)
		if err != nil {
			t.Fatalf("failed to encode the token: %v", err)
		}
		if err := os.WriteFile(cachePath, body, 0o600); err != nil {
			t.Fatalf("failed to seed the token cache: %v", err)
		}
		return cachePath
	}

	// tokenEndpoint stands in for /platform/oauth/request_token.
	tokenEndpoint := func(t *testing.T, body string) (oauth2.Endpoint, *url.Values) {
		t.Helper()

		var seen url.Values
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if err := r.ParseForm(); err != nil {
				t.Errorf("failed to parse the token request: %v", err)
			}
			seen = r.PostForm
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(body))
		}))
		t.Cleanup(server.Close)

		return oauth2.Endpoint{TokenURL: server.URL, AuthStyle: oauth2.AuthStyleInParams}, &seen
	}

	// The whole feature in one assertion: an expired token renews from
	// disk, with no browser involved.
	t.Run("Renews an expired token from the cached refresh token", func(t *testing.T) {
		RegisterTestingT(t)

		cachePath := seedCache(t, &oauth2.Token{
			AccessToken:  "stale",
			RefreshToken: "stored-rt",
			Expiry:       time.Now().Add(-time.Hour),
		})
		endpoint, seen := tokenEndpoint(t,
			`{"access_token":"renewed","refresh_token":"rotated-rt","expires_in":86400}`)

		app := geniapp.Credentials{ClientID: "1855", ClientSecret: "app-secret"}
		token, err := browserTokenSource(app, endpoint, cachePath).Token()

		Expect(err).ToNot(HaveOccurred())
		Expect(token.AccessToken).To(Equal("renewed"))
		Expect(seen.Get("grant_type")).To(Equal("refresh_token"))
		Expect(seen.Get("refresh_token")).To(Equal("stored-rt"))
		Expect(seen.Get("client_secret")).To(Equal("app-secret"))
	})

	// Geni rotates the refresh token on every renewal, so the new one has
	// to reach disk — Terraform exits after every command.
	t.Run("Writes the rotated refresh token back to the cache", func(t *testing.T) {
		RegisterTestingT(t)

		cachePath := seedCache(t, &oauth2.Token{
			AccessToken:  "stale",
			RefreshToken: "stored-rt",
			Expiry:       time.Now().Add(-time.Hour),
		})
		endpoint, _ := tokenEndpoint(t,
			`{"access_token":"renewed","refresh_token":"rotated-rt","expires_in":86400}`)

		app := geniapp.Credentials{ClientID: "1855", ClientSecret: "app-secret"}
		_, err := browserTokenSource(app, endpoint, cachePath).Token()
		Expect(err).ToNot(HaveOccurred())

		body, err := os.ReadFile(cachePath)
		Expect(err).ToNot(HaveOccurred())

		var stored oauth2.Token
		Expect(json.Unmarshal(body, &stored)).To(Succeed())
		Expect(stored.AccessToken).To(Equal("renewed"))
		Expect(stored.RefreshToken).To(Equal("rotated-rt"))
	})

	// Without a secret the provider keeps the flow it has always used,
	// and a cached token is still served straight from disk.
	t.Run("Serves a valid cached token without a client secret", func(t *testing.T) {
		RegisterTestingT(t)

		cachePath := seedCache(t, &oauth2.Token{
			AccessToken: "cached",
			Expiry:      time.Now().Add(time.Hour),
		})

		app := geniapp.Credentials{ClientID: "1855"}
		token, err := browserTokenSource(app, oauth2.Endpoint{}, cachePath).Token()

		Expect(err).ToNot(HaveOccurred())
		Expect(token.AccessToken).To(Equal("cached"))
	})
}
