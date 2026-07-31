package config

import (
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/dmalch/go-geni"
	"github.com/dmalch/terraform-provider-genealogy/internal/genibatch"
)

type GeniProviderConfig struct {
	AccessToken              types.String `tfsdk:"access_token"`
	ClientID                 types.String `tfsdk:"client_id"`
	ClientSecret             types.String `tfsdk:"client_secret"`
	UseSandboxEnv            types.Bool   `tfsdk:"use_sandbox_env"`
	AutoUpdateMergedProfiles types.Bool   `tfsdk:"auto_update_merged_profiles"`
}

type ClientData struct {
	Client                   *geni.Client
	BatchClient              *genibatch.Client
	AutoUpdateMergedProfiles bool
}
