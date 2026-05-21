package acceptance

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccPhoto_createPhoto(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckPhotoDestroy,
		Steps: []resource.TestStep{
			{
				Config: `
					resource "geni_photo" "test" {
					  title     = "Acceptance test photo"
					  file      = filebase64("${path.module}/assets/cs-white-fff.png")
					  file_name = "cs-white-fff.png"
					}
				`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("geni_photo.test", tfjsonpath.New("id"),
						knownvalue.StringRegexp(regexp.MustCompile(`^photo-\d+$`))),
					statecheck.ExpectKnownValue("geni_photo.test", tfjsonpath.New("title"),
						knownvalue.StringExact("Acceptance test photo")),
				},
			},
		},
	})
}

func TestAccPhoto_createPhotoWithDescription(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckPhotoDestroy,
		Steps: []resource.TestStep{
			{
				Config: `
					resource "geni_photo" "test" {
					  title       = "Acceptance test photo with description"
					  description = "A photo created by the acceptance test suite."
					  file        = filebase64("${path.module}/assets/cs-white-fff.png")
					  file_name   = "cs-white-fff.png"
					}
				`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("geni_photo.test", tfjsonpath.New("description"),
						knownvalue.StringExact("A photo created by the acceptance test suite.")),
				},
			},
		},
	})
}
