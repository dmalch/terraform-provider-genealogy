package acceptance

import "os"

// testAccessToken is the access token used for acceptance tests.
// Can be set via the GENI_ACCESS_TOKEN environment variable or
// requested at https://sandbox.geni.com/platform/developer/api_explorer.
// If empty, the provider falls back to browser-based OAuth login.
var testAccessToken = os.Getenv("GENI_ACCESS_TOKEN")

func init() {
	os.Setenv("GENI_USE_SANDBOX", "true")
}
