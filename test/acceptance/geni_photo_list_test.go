package acceptance

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/querycheck"
)

func TestAccPhoto_listResources(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckPhotoDestroy,
		Steps: []resource.TestStep{
			{
				// Apply step seeds the sandbox account with an uploaded photo
				// the next step's query is expected to surface.
				Config: `
					resource "geni_photo" "test" {
					  title     = "ListResource Test Photo"
					  file      = filebase64("${path.module}/assets/cs-white-fff.png")
					  file_name = "cs-white-fff.png"
					}
				`,
			},
			{
				// terraform query reads list blocks; the previously uploaded
				// photo stays on Geni (and in state across steps) so the query
				// surfaces at least one result.
				Config: `
					list "geni_photo" "all" {
					  provider         = geni
					  include_resource = true
					}
				`,
				Query: true,
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLengthAtLeast("list.geni_photo.all", 1),
				},
			},
		},
	})
}
