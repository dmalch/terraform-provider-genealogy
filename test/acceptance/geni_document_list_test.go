package acceptance

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/querycheck"
)

func TestAccDocument_listResources(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckDocumentDestroy,
		Steps: []resource.TestStep{
			{
				// Apply step seeds the sandbox account with an uploaded
				// document the next step's query is expected to surface.
				Config: `
					resource "geni_document" "test" {
					  title      = "ListResource Test Document"
					  source_url = "https://example.com/list-test"
					}
				`,
			},
			{
				// terraform query reads a *.tfquery.hcl file containing only
				// list blocks. The previously uploaded document remains on
				// Geni (and stays in state across steps) so the query
				// surfaces at least one result.
				Config: `
					list "geni_document" "all" {
					  provider         = geni
					  include_resource = true
					}
				`,
				Query: true,
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLengthAtLeast("list.geni_document.all", 1),
				},
			},
		},
	})
}
