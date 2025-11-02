package profile

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/dmalch/terraform-provider-genealogy/internal/config"
	"github.com/dmalch/terraform-provider-genealogy/internal/geni"
	"github.com/dmalch/terraform-provider-genealogy/internal/genibatch"
	"github.com/dmalch/terraform-provider-genealogy/internal/genicache"
)

var _ resource.Resource = &Resource{}
var _ resource.ResourceWithUpgradeState = &Resource{}

type Resource struct {
	resource.ResourceWithConfigure
	client                   *geni.Client
	batchClient              *genibatch.Client
	cacheClient              *genicache.Client
	useProfileCache          bool
	autoUpdateMergedProfiles bool
}

func NewProfileResource() resource.Resource {
	return &Resource{}
}

// Metadata provides the resource type name.
func (r *Resource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "geni_profile"
}

func (r *Resource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Always perform a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	cfg, ok := req.ProviderData.(*config.ClientData)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *config.ClientData, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = cfg.Client
	r.batchClient = cfg.BatchClient
	r.cacheClient = cfg.CacheClient
	r.useProfileCache = cfg.UseProfileCache
	r.autoUpdateMergedProfiles = cfg.AutoUpdateMergedProfiles
}
