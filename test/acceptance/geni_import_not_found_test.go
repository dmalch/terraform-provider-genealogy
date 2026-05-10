package acceptance

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// These tests cover the import-time validation added in
// https://github.com/dmalch/terraform-provider-genealogy/issues/80.
// Before the fix, ImportState used a bare passthrough that wrote the
// user-supplied ID into state without contacting Geni, leaving a zombie
// row when the ID did not exist. The new ImportState round-trips the API
// up front and surfaces a clear "not found" diagnostic instead.

func TestAccDocument_importNonExistentIdFails(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					resource "geni_document" "missing" {
					  title = "placeholder"
					  text  = "placeholder"
					}
				`,
				ResourceName:  "geni_document.missing",
				ImportState:   true,
				ImportStateId: "document-99999999999",
				ExpectError:   regexp.MustCompile(`(?s)Document not found`),
			},
		},
	})
}

func TestAccProfile_importNonExistentIdFails(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					resource "geni_profile" "missing" {
					  names = {
						"en-US" = {
						  first_name = "Placeholder"
						  last_name  = "Placeholder"
						}
					  }
					  alive  = false
					  public = true
					}
				`,
				ResourceName:  "geni_profile.missing",
				ImportState:   true,
				ImportStateId: "profile-99999999999",
				ExpectError:   regexp.MustCompile(`(?s)Profile not found`),
			},
		},
	})
}

func TestAccUnion_importNonExistentIdFails(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					resource "geni_union" "missing" {
					  partners = ["profile-1", "profile-2"]
					}
				`,
				ResourceName:  "geni_union.missing",
				ImportState:   true,
				ImportStateId: "union-99999999999",
				ExpectError:   regexp.MustCompile(`(?s)Union not found`),
			},
		},
	})
}
