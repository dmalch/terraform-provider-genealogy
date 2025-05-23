package config

import (
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/dmalch/terraform-provider-genealogy/internal/geni"
	"github.com/dmalch/terraform-provider-genealogy/internal/genibatch"
	"github.com/dmalch/terraform-provider-genealogy/internal/genicache"
)

type GeniProviderConfig struct {
	AccessToken              types.String `tfsdk:"access_token"`
	UseSandboxEnv            types.Bool   `tfsdk:"use_sandbox_env"`
	UseProfileCache          types.Bool   `tfsdk:"use_profile_cache"`
	UseDocumentCache         types.Bool   `tfsdk:"use_document_cache"`
	AutoUpdateMergedProfiles types.Bool   `tfsdk:"auto_update_merged_profiles"`
}

type ClientData struct {
	Client                   *geni.Client
	BatchClient              *genibatch.Client
	CacheClient              *genicache.Client
	UseProfileCache          bool
	UseDocumentCache         bool
	AutoUpdateMergedProfiles bool
}
