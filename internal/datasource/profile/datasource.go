package profile

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/dmalch/terraform-provider-genealogy/internal/config"
	"github.com/dmalch/terraform-provider-genealogy/internal/geni"
	"github.com/dmalch/terraform-provider-genealogy/internal/genibatch"
	"github.com/dmalch/terraform-provider-genealogy/internal/genicache"
)

var _ datasource.DataSource = &DataSource{}

type DataSource struct {
	datasource.DataSourceWithConfigure
	client                   *geni.Client
	batchClient              *genibatch.Client
	cacheClient              *genicache.Client
	useProfileCache          bool
	autoUpdateMergedProfiles bool
}

func NewDataSource() datasource.DataSource {
	return &DataSource{}
}

func (d *DataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "geni_profile"
}

func (d *DataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	cfg, ok := req.ProviderData.(*config.ClientData)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected DataSource Configure Type",
			fmt.Sprintf("Expected *config.ClientData, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.client = cfg.Client
	d.batchClient = cfg.BatchClient
	d.cacheClient = cfg.CacheClient
	d.useProfileCache = cfg.UseProfileCache
	d.autoUpdateMergedProfiles = cfg.AutoUpdateMergedProfiles
}
