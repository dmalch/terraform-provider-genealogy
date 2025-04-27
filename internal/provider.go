package internal

import (
	"context"
	"fmt"
	"os"
	"path"
	"sync"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"golang.org/x/oauth2"

	"github.com/dmalch/terraform-provider-genealogy/internal/authn"
	"github.com/dmalch/terraform-provider-genealogy/internal/config"
	"github.com/dmalch/terraform-provider-genealogy/internal/datasource/project"
	"github.com/dmalch/terraform-provider-genealogy/internal/geni"
	"github.com/dmalch/terraform-provider-genealogy/internal/genibatch"
	"github.com/dmalch/terraform-provider-genealogy/internal/genicache"
	"github.com/dmalch/terraform-provider-genealogy/internal/resource/document"
	"github.com/dmalch/terraform-provider-genealogy/internal/resource/profile"
	"github.com/dmalch/terraform-provider-genealogy/internal/resource/union"
)

type GeniProvider struct {
	initClientOnce sync.Once
}

var (
	client      *geni.Client
	batchClient *genibatch.Client
	cacheClient *genicache.Client
)

func New() provider.Provider {
	return &GeniProvider{
		initClientOnce: sync.Once{},
	}
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
				Description: "The Access Token for the Geni API, if not provided the provider will attempt to do a browser-based OAuth login flow",
			},
			"use_sandbox_env": schema.BoolAttribute{
				Optional:    true,
				Description: "Whether to use the Geni Sandbox environment",
			},
			"use_profile_cache": schema.BoolAttribute{
				Optional:    true,
				Description: "Whether to use the profile cache for faster lookups",
			},
			"use_document_cache": schema.BoolAttribute{
				Optional:    true,
				Description: "Whether to use the document cache for faster lookups",
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

	cacheFilePath, err := tokenCacheFilePath()
	if err != nil {
		resp.Diagnostics.AddError("error getting token cache file path", err.Error())
		return
	}

	var tokenSource = oauth2.ReuseTokenSource(nil,
		authn.NewCachingTokenSource(
			cacheFilePath,
			authn.NewAuthTokenSource(&oauth2.Config{
				ClientID: clientId(cfg.UseSandboxEnv.ValueBool()),
				Endpoint: oauth2.Endpoint{
					AuthURL: geni.BaseUrl(cfg.UseSandboxEnv.ValueBool()) + "platform/oauth/authorize",
				},
			})))

	if cfg.AccessToken.ValueString() != "" {
		tokenSource = oauth2.StaticTokenSource(&oauth2.Token{AccessToken: cfg.AccessToken.ValueString()})
	}

	p.initClientOnce.Do(func() {
		client = geni.NewClient(tokenSource, cfg.UseSandboxEnv.ValueBool())
		batchClient = genibatch.NewClient(client)
		cacheClient = genicache.NewClient(client, batchClient)
		go batchClient.UnionBulkProcessor(context.Background())
		go batchClient.ProfileBulkProcessor(context.Background())
	})

	resp.ResourceData = &config.ClientData{
		Client:           client,
		BatchClient:      batchClient,
		CacheClient:      cacheClient,
		UseProfileCache:  cfg.UseProfileCache.ValueBool(),
		UseDocumentCache: cfg.UseDocumentCache.ValueBool(),
	}

	resp.DataSourceData = &config.ClientData{
		Client: client,
	}
}

func tokenCacheFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("error getting user home directory: %w", err)
	}

	cacheFilePath := path.Join(homeDir, ".genealogy", "geni_token.json")
	return cacheFilePath, nil
}

func clientId(useSandboxEnv bool) string {
	if useSandboxEnv {
		// Sandbox client ID
		return "8"
	}

	// Production client ID
	return "1855"
}

func (p *GeniProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		profile.NewProfileResource,
		union.NewUnionResource,
		document.NewResource,
	}
}

func (p *GeniProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		project.NewDataSource,
	}
}
