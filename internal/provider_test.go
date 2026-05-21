package internal

import (
	"os"
	"path"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	. "github.com/onsi/gomega"

	"github.com/dmalch/terraform-provider-genealogy/internal/config"
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

func TestClientId(t *testing.T) {
	t.Run("production", func(t *testing.T) {
		RegisterTestingT(t)

		Expect(clientId(false)).To(Equal("1855"))
	})

	t.Run("sandbox", func(t *testing.T) {
		RegisterTestingT(t)

		Expect(clientId(true)).To(Equal("8"))
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
		"use_sandbox_env":             tftypes.NewValue(tftypes.Bool, nil),
		"use_profile_cache":           tftypes.NewValue(tftypes.Bool, nil),
		"use_document_cache":          tftypes.NewValue(tftypes.Bool, nil),
		"auto_update_merged_profiles": tftypes.NewValue(tftypes.Bool, nil),
	})

	resp := &provider.ConfigureResponse{}
	p.Configure(ctx, provider.ConfigureRequest{
		Config: tfsdk.Config{Schema: schemaResp.Schema, Raw: raw},
	}, resp)
	return resp
}
