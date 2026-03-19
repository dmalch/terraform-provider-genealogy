package acceptance

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"golang.org/x/oauth2"

	"github.com/dmalch/terraform-provider-genealogy/internal"
	"github.com/dmalch/terraform-provider-genealogy/internal/geni"
)

var testProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"geni": providerserver.NewProtocol6WithError(internal.New()),
}

func newTestClient() *geni.Client {
	var tokenSource oauth2.TokenSource
	if testAccessToken != "" {
		tokenSource = oauth2.StaticTokenSource(&oauth2.Token{AccessToken: testAccessToken})
	} else {
		tokenSource = oauth2.StaticTokenSource(&oauth2.Token{})
	}
	return geni.NewClient(tokenSource, true)
}

func testAccCheckProfileDestroy(s *terraform.State) error {
	if testAccessToken == "" {
		return nil
	}
	client := newTestClient()
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "geni_profile" {
			continue
		}
		profile, err := client.GetProfile(context.Background(), rs.Primary.ID)
		if errors.Is(err, geni.ErrResourceNotFound) {
			continue
		}
		if err != nil {
			return fmt.Errorf("error checking profile %s: %w", rs.Primary.ID, err)
		}
		if !profile.Deleted {
			return fmt.Errorf("profile %s still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckDocumentDestroy(s *terraform.State) error {
	if testAccessToken == "" {
		return nil
	}
	client := newTestClient()
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "geni_document" {
			continue
		}
		_, err := client.GetDocument(context.Background(), rs.Primary.ID)
		if errors.Is(err, geni.ErrResourceNotFound) {
			continue
		}
		if err != nil {
			return fmt.Errorf("error checking document %s: %w", rs.Primary.ID, err)
		}
		return fmt.Errorf("document %s still exists", rs.Primary.ID)
	}
	return nil
}
