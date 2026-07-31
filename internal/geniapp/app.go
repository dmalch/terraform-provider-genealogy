// Package geniapp resolves the registered Geni OAuth application the
// provider authenticates as.
//
// A client secret is what unlocks Geni's server-side flow, which returns a
// refresh token; without one the provider falls back to the client-side
// flow and its browser round trip every 24 hours.
package geniapp

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Built-in application ids, used when the caller does not bring their own.
const (
	prodClientID    = "1855"
	sandboxClientID = "8"
)

// Credentials identify a registered Geni application.
type Credentials struct {
	ClientID     string
	ClientSecret string
}

// Refreshable reports whether these credentials can drive the
// server-side flow.
func (c Credentials) Refreshable() bool {
	return c.ClientSecret != ""
}

// Explicit is what the practitioner wrote in the provider block. Empty
// strings mean "not set".
type Explicit struct {
	ClientID     string
	ClientSecret string
}

// Resolve picks the application to authenticate as, in descending
// precedence:
//
//  1. the GENI_CLIENT_ID / GENI_CLIENT_SECRET environment variables
//     (GENI_SANDBOX_* under the sandbox environment);
//  2. the provider configuration block;
//  3. ~/.genealogy/config.json, which the geni CLI writes — so
//     "geni config client-secret" configures both tools at once, the way
//     they already share a token cache.
//
// The id travels with the secret at each level: the secret for the
// built-in application belongs to its owner, so anyone else registers
// their own application and supplies both halves. A level that carries
// only a secret keeps the built-in id.
//
// Resolving never fails on an unreadable config file. The worst case is
// the behavior of every release before this one — an interactive login.
func Resolve(explicit Explicit, useSandboxEnv bool) Credentials {
	for _, candidate := range []Credentials{
		fromEnvironment(useSandboxEnv),
		{ClientID: explicit.ClientID, ClientSecret: explicit.ClientSecret},
		fromCLIConfig(useSandboxEnv),
	} {
		if candidate.ClientID == "" && candidate.ClientSecret == "" {
			continue
		}
		if candidate.ClientID == "" {
			candidate.ClientID = DefaultClientID(useSandboxEnv)
		}
		return candidate
	}
	return Credentials{ClientID: DefaultClientID(useSandboxEnv)}
}

// DefaultClientID returns the built-in application id for the
// environment.
func DefaultClientID(useSandboxEnv bool) string {
	if useSandboxEnv {
		return sandboxClientID
	}
	return prodClientID
}

func fromEnvironment(useSandboxEnv bool) Credentials {
	idVar, secretVar := "GENI_CLIENT_ID", "GENI_CLIENT_SECRET"
	if useSandboxEnv {
		idVar, secretVar = "GENI_SANDBOX_CLIENT_ID", "GENI_SANDBOX_CLIENT_SECRET"
	}
	return Credentials{
		ClientID:     os.Getenv(idVar),
		ClientSecret: os.Getenv(secretVar),
	}
}

// cliConfig is the part of the geni CLI's ~/.genealogy/config.json this
// provider cares about. The file belongs to the CLI; unknown fields are
// ignored, and a format change there costs nothing worse here than
// falling back to an interactive login.
type cliConfig struct {
	Prod    cliOAuthApp `json:"prod"`
	Sandbox cliOAuthApp `json:"sandbox"`
}

type cliOAuthApp struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

func fromCLIConfig(useSandboxEnv bool) Credentials {
	path, err := ConfigFilePath()
	if err != nil {
		return Credentials{}
	}

	// A missing or unreadable file is not an error worth a diagnostic:
	// the login still works, it just asks for a browser.
	body, err := os.ReadFile(path)
	if err != nil {
		return Credentials{}
	}

	var c cliConfig
	if err := json.Unmarshal(body, &c); err != nil {
		return Credentials{}
	}

	app := c.Prod
	if useSandboxEnv {
		app = c.Sandbox
	}
	return Credentials(app)
}

// ConfigFilePath returns the geni CLI's configuration file, which lives
// beside the token cache the two tools share.
func ConfigFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".genealogy", "config.json"), nil
}
