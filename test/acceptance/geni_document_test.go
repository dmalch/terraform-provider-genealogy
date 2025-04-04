package acceptance

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	"github.com/dmalch/terraform-provider-genealogy/internal"
)

func TestAccProfile_createDocument(t *testing.T) {
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
			
					resource "geni_document" "test" {
					  title = "Test Document"
					  text = "This is a test document."
					}
					`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("geni_document.test", tfjsonpath.New("title"), knownvalue.StringExact("Test Document")),
					statecheck.ExpectKnownValue("geni_document.test", tfjsonpath.New("content_type"), knownvalue.StringExact("text/plain")),
					statecheck.ExpectKnownValue("geni_document.test", tfjsonpath.New("text"), knownvalue.StringExact("This is a test document.")),
				},
			},
		},
	})
}
