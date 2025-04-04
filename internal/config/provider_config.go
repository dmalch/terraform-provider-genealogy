package config

import (
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/dmalch/terraform-provider-genealogy/internal/geni"
)

type GeniProviderConfig struct {
	AccessToken   types.String `tfsdk:"access_token"`
	UseSandboxEnv types.Bool   `tfsdk:"use_sandbox_env"`
}

type ClientData struct {
	Client *geni.Client
}
