package internal

import (
	"context"
	"fmt"
	"os"
	"path"
	"sync"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"golang.org/x/oauth2"

	"github.com/dmalch/go-geni"
	"github.com/dmalch/go-geni/auth"
	"github.com/dmalch/terraform-provider-genealogy/internal/config"
	profiledatasource "github.com/dmalch/terraform-provider-genealogy/internal/datasource/profile"
	"github.com/dmalch/terraform-provider-genealogy/internal/datasource/project"
	"github.com/dmalch/terraform-provider-genealogy/internal/geniapp"
	"github.com/dmalch/terraform-provider-genealogy/internal/genibatch"
	"github.com/dmalch/terraform-provider-genealogy/internal/resource/document"
	"github.com/dmalch/terraform-provider-genealogy/internal/resource/photo"
	"github.com/dmalch/terraform-provider-genealogy/internal/resource/profile"
	"github.com/dmalch/terraform-provider-genealogy/internal/resource/union"
)

var _ provider.ProviderWithListResources = (*GeniProvider)(nil)

// GeniProvider holds the configured API clients. State lives on the instance
// (not in package globals) so each provider is self-contained: every
// GeniProvider.New() is independent, and Configure is unit-testable.
type GeniProvider struct {
	// once guards one-time creation of the clients and batch-processor
	// goroutines; Configure may be invoked more than once on an instance.
	once        sync.Once
	client      *geni.Client
	batchClient *genibatch.Client
}

func New() provider.Provider {
	return &GeniProvider{}
}

func (p *GeniProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "geni"
}

func (p *GeniProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"access_token": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "The Access Token for the Geni API. Can also be set with the GENI_ACCESS_TOKEN environment variable. If not provided, the provider will attempt to do a browser-based OAuth login flow.",
			},
			"client_id": schema.StringAttribute{
				Optional:    true,
				Description: "The OAuth client id of your own registered Geni application, if you do not want to use the built-in one. Can also be set with the GENI_CLIENT_ID environment variable (GENI_SANDBOX_CLIENT_ID under the sandbox environment). Supply it together with `client_secret`.",
			},
			"client_secret": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "The OAuth client secret of the Geni application. Supplying it switches the browser login to Geni's server-side flow, which returns a refresh token, so the provider renews the token in the background instead of opening a browser every 24 hours. Can also be set with the GENI_CLIENT_SECRET environment variable (GENI_SANDBOX_CLIENT_SECRET under the sandbox environment), or with `geni config client-secret` in the geni CLI, which stores it in ~/.genealogy/config.json alongside the shared token cache.",
			},
			"use_sandbox_env": schema.BoolAttribute{
				Optional:    true,
				Description: "Whether to use the Geni Sandbox environment. Can also be set with the GENI_USE_SANDBOX environment variable.",
			},
			"auto_update_merged_profiles": schema.BoolAttribute{
				Optional:    true,
				Description: "Whether to automatically update merged profiles in the state",
			},
		},
		Description: "This provider enables managing data on Geni.com through Terraform. It exposes configuration and resources that help automate genealogical information. This application uses the Geni API but is not endorsed, operated, or sponsored by Geni.com.",
	}
}

func (p *GeniProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var cfg config.GeniProviderConfig

	resp.Diagnostics.Append(req.Config.Get(ctx, &cfg)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Resolve configuration with environment variable fallbacks.
	accessToken := cfg.AccessToken.ValueString()
	if accessToken == "" {
		accessToken = os.Getenv("GENI_ACCESS_TOKEN")
	}

	useSandboxEnv := cfg.UseSandboxEnv.ValueBool()
	if cfg.UseSandboxEnv.IsNull() {
		useSandboxEnv = os.Getenv("GENI_USE_SANDBOX") == "true"
	}

	cacheFilePath, err := tokenCacheFilePath(useSandboxEnv)
	if err != nil {
		resp.Diagnostics.AddError("error getting token cache file path", err.Error())
		return
	}

	app := geniapp.Resolve(geniapp.Explicit{
		ClientID:     cfg.ClientID.ValueString(),
		ClientSecret: cfg.ClientSecret.ValueString(),
	}, useSandboxEnv)

	tokenSource := browserTokenSource(app, auth.GeniEndpoint(geni.BaseURL(useSandboxEnv)), cacheFilePath)

	if accessToken != "" {
		tokenSource = oauth2.StaticTokenSource(&oauth2.Token{AccessToken: accessToken})
	}

	p.once.Do(func() {
		p.client = geni.NewClient(tokenSource, useSandboxEnv)
		p.batchClient = genibatch.NewClient(p.client)
		go p.batchClient.UnionBulkProcessor(context.Background())
		go p.batchClient.ProfileBulkProcessor(context.Background())
		go p.batchClient.DocumentBulkProcessor(context.Background())
		go p.batchClient.PhotoBulkProcessor(context.Background())
	})

	resp.ResourceData = &config.ClientData{
		Client:                   p.client,
		BatchClient:              p.batchClient,
		AutoUpdateMergedProfiles: cfg.AutoUpdateMergedProfiles.ValueBool(),
	}

	resp.DataSourceData = &config.ClientData{
		Client:                   p.client,
		BatchClient:              p.batchClient,
		AutoUpdateMergedProfiles: cfg.AutoUpdateMergedProfiles.ValueBool(),
	}

	resp.ListResourceData = &config.ClientData{
		Client: p.client,
	}
}

// browserTokenSource builds the token source behind an interactive
// login.
//
// With a client secret it runs Geni's server-side flow, which returns a
// refresh token, and renews from it instead of opening a browser once a
// day. Renewing happens inside the caching source, below
// oauth2.ReuseTokenSource: Geni rotates the refresh token on every
// renewal, so a refresher above the cache would renew into memory and
// lose the new token when Terraform exits — which, for a provider, is
// after every single command.
//
// Without a secret this is the client-side flow the provider has always
// used, and its 24-hour token.
// The endpoint is a parameter rather than derived here so a test can
// stand one up and exercise the refresh without a browser.
func browserTokenSource(app geniapp.Credentials, endpoint oauth2.Endpoint, cacheFilePath string) oauth2.TokenSource {
	oauthConfig := &oauth2.Config{
		ClientID:     app.ClientID,
		ClientSecret: app.ClientSecret,
		Endpoint:     endpoint,
	}

	if !app.Refreshable() {
		return oauth2.ReuseTokenSource(nil,
			auth.NewCachingTokenSource(cacheFilePath, auth.NewAuthTokenSource(oauthConfig)))
	}

	codeSource := auth.NewCodeTokenSource(oauthConfig)
	return oauth2.ReuseTokenSource(nil,
		auth.NewRefreshingCachingTokenSource(cacheFilePath, codeSource, codeSource))
}

func tokenCacheFilePath(useSandboxEnv bool) (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("error getting user home directory: %w", err)
	}

	cacheFilePath := path.Join(homeDir, ".genealogy", "geni_token.json")
	if useSandboxEnv {
		cacheFilePath = path.Join(homeDir, ".genealogy", "geni_sandbox_token.json")
	}

	return cacheFilePath, nil
}

func (p *GeniProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		profile.NewProfileResource,
		union.NewUnionResource,
		document.NewResource,
		photo.NewResource,
	}
}

func (p *GeniProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		project.NewDataSource,
		profiledatasource.NewDataSource,
	}
}

func (p *GeniProvider) ListResources(_ context.Context) []func() list.ListResource {
	return []func() list.ListResource{
		profile.NewListResource,
		document.NewListResource,
		photo.NewListResource,
	}
}
