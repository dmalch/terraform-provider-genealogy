package config

import "github.com/hashicorp/terraform-plugin-framework/types"

type GeniProviderConfig struct {
	ApiKey types.String `tfsdk:"api_key"`
}
