package acceptance

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/querycheck"
)

func TestAccProfile_listResources(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckProfileDestroy,
		Steps: []resource.TestStep{
			{
				// Apply step seeds the sandbox account with a managed profile
				// the next step's query is expected to surface.
				Config: `
					resource "geni_profile" "test" {
					  names = {
						"en-US" = {
							first_name = "ListResource"
							last_name  = "TestProfile"
						}
					  }
					  alive  = false
					  public = true
					}
				`,
			},
			{
				// terraform query reads a *.tfquery.hcl file containing only
				// list blocks. The framework writes step.Config to that file
				// when Query is true; the previously created profile is still
				// on Geni (and stays in state across steps) so the query
				// surfaces at least one result.
				Config: `
					list "geni_profile" "all" {
					  provider         = geni
					  include_resource = true
					}
				`,
				Query: true,
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLengthAtLeast("list.geni_profile.all", 1),
				},
			},
		},
	})
}
