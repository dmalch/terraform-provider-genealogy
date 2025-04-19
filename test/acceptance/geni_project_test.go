package acceptance

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/dmalch/terraform-provider-genealogy/internal"
)

func TestAccProject_getProject(t *testing.T) {
	resource.Test(t, resource.TestCase{
		//IsUnitTest: true,
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"geni": providerserver.NewProtocol6WithError(internal.New()),
		},
		Steps: []resource.TestStep{
			{
				Config: `
					provider "geni" {
					  access_token = "` + testAccessToken + `"
					  use_sandbox_env = true
					}
			
					data "geni_project" "test" {
					  id = "project-6"
					}
					`,
				PlanOnly: true,
			},
		},
	})
}
