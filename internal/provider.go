package internal

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/dmalch/terraform-provider-genealogy/internal/config"
	"github.com/dmalch/terraform-provider-genealogy/internal/geni"
	"github.com/dmalch/terraform-provider-genealogy/internal/resource/profile"
	"github.com/dmalch/terraform-provider-genealogy/internal/resource/union"
)

type GeniProvider struct {
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
				Required:    true,
				Sensitive:   true,
				Description: "The Access Token for the Geni API",
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

	resp.ResourceData = &config.ClientData{
		Client: geni.NewClient(cfg.AccessToken.ValueString()),
	}
}

func (p *GeniProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		profile.NewProfileResource,
		union.NewUnionResource,
	}
}

func (p *GeniProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return nil
}
