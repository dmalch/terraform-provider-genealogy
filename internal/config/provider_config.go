package config

import "github.com/hashicorp/terraform-plugin-framework/types"

type GeniProviderConfig struct {
	AccessToken types.String `tfsdk:"access_token"`
}
