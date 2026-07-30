package geniapp

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	. "github.com/onsi/gomega"
)

// isolate points HOME at an empty directory and clears every variable
// Resolve consults, so a developer's own exports cannot decide a test.
func isolate(t *testing.T) {
	t.Helper()

	t.Setenv("HOME", t.TempDir())
	for _, name := range []string{
		"GENI_CLIENT_ID", "GENI_CLIENT_SECRET",
		"GENI_SANDBOX_CLIENT_ID", "GENI_SANDBOX_CLIENT_SECRET",
	} {
		t.Setenv(name, "")
	}
}

// writeCLIConfig plants the file the geni CLI writes.
func writeCLIConfig(t *testing.T, body any) {
	t.Helper()

	path, err := ConfigFilePath()
	if err != nil {
		t.Fatalf("failed to build the config path: %v", err)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		t.Fatalf("failed to create the config directory: %v", err)
	}

	var encoded []byte
	switch v := body.(type) {
	case string:
		encoded = []byte(v)
	default:
		encoded, err = json.Marshal(v)
		if err != nil {
			t.Fatalf("failed to encode the config: %v", err)
		}
	}
	if err := os.WriteFile(path, encoded, 0o600); err != nil {
		t.Fatalf("failed to write the config: %v", err)
	}
}

func TestResolve(t *testing.T) {
	t.Run("Falls back to the built-in application with no secret", func(t *testing.T) {
		RegisterTestingT(t)
		isolate(t)

		app := Resolve(Explicit{}, false)
		Expect(app.ClientID).To(Equal("1855"))
		Expect(app.ClientSecret).To(BeEmpty())
		Expect(app.Refreshable()).To(BeFalse())
	})

	t.Run("Reads the environment first", func(t *testing.T) {
		RegisterTestingT(t)
		isolate(t)
		writeCLIConfig(t, cliConfig{Prod: cliOAuthApp{ClientID: "cfg", ClientSecret: "cfg-secret"}})
		t.Setenv("GENI_CLIENT_ID", "env")
		t.Setenv("GENI_CLIENT_SECRET", "env-secret")

		app := Resolve(Explicit{ClientID: "hcl", ClientSecret: "hcl-secret"}, false)
		Expect(app.ClientID).To(Equal("env"))
		Expect(app.ClientSecret).To(Equal("env-secret"))
	})

	t.Run("Prefers the provider block over the CLI config", func(t *testing.T) {
		RegisterTestingT(t)
		isolate(t)
		writeCLIConfig(t, cliConfig{Prod: cliOAuthApp{ClientID: "cfg", ClientSecret: "cfg-secret"}})

		app := Resolve(Explicit{ClientID: "hcl", ClientSecret: "hcl-secret"}, false)
		Expect(app.ClientID).To(Equal("hcl"))
		Expect(app.ClientSecret).To(Equal("hcl-secret"))
	})

	// The point of reading the CLI's file: "geni config client-secret"
	// configures both tools, the way they already share a token cache.
	t.Run("Falls back to the secret stored by the geni CLI", func(t *testing.T) {
		RegisterTestingT(t)
		isolate(t)
		writeCLIConfig(t, cliConfig{Prod: cliOAuthApp{ClientSecret: "cfg-secret"}})

		app := Resolve(Explicit{}, false)
		Expect(app.ClientID).To(Equal("1855"))
		Expect(app.ClientSecret).To(Equal("cfg-secret"))
		Expect(app.Refreshable()).To(BeTrue())
	})

	// A secret alone is enough, and keeps the built-in id: that is the
	// case for whoever owns the built-in application.
	t.Run("Keeps the built-in id when only a secret is supplied", func(t *testing.T) {
		RegisterTestingT(t)
		isolate(t)

		app := Resolve(Explicit{ClientSecret: "hcl-secret"}, false)
		Expect(app.ClientID).To(Equal("1855"))
		Expect(app.ClientSecret).To(Equal("hcl-secret"))
	})

	// The id travels with the secret: a level that sets an id must not
	// pick up a secret belonging to a different application.
	t.Run("Never mixes an id from one level with a secret from another", func(t *testing.T) {
		RegisterTestingT(t)
		isolate(t)
		writeCLIConfig(t, cliConfig{Prod: cliOAuthApp{ClientID: "cfg", ClientSecret: "cfg-secret"}})

		app := Resolve(Explicit{ClientID: "hcl"}, false)
		Expect(app.ClientID).To(Equal("hcl"))
		Expect(app.ClientSecret).To(BeEmpty())
	})

	t.Run("Keeps the environments apart", func(t *testing.T) {
		RegisterTestingT(t)
		isolate(t)
		writeCLIConfig(t, cliConfig{
			Prod:    cliOAuthApp{ClientSecret: "prod-secret"},
			Sandbox: cliOAuthApp{ClientSecret: "sandbox-secret"},
		})

		Expect(Resolve(Explicit{}, false).ClientSecret).To(Equal("prod-secret"))

		sandbox := Resolve(Explicit{}, true)
		Expect(sandbox.ClientID).To(Equal("8"))
		Expect(sandbox.ClientSecret).To(Equal("sandbox-secret"))
	})

	t.Run("Reads the sandbox environment variables under the sandbox", func(t *testing.T) {
		RegisterTestingT(t)
		isolate(t)
		t.Setenv("GENI_CLIENT_SECRET", "prod-secret")
		t.Setenv("GENI_SANDBOX_CLIENT_SECRET", "sandbox-secret")

		Expect(Resolve(Explicit{}, true).ClientSecret).To(Equal("sandbox-secret"))
		Expect(Resolve(Explicit{}, false).ClientSecret).To(Equal("prod-secret"))
	})

	// The config file belongs to the CLI. If it ever changes shape, the
	// worst this may cost is the behavior of every earlier release.
	t.Run("Falls back to an interactive login when the CLI config is unreadable", func(t *testing.T) {
		RegisterTestingT(t)
		isolate(t)
		writeCLIConfig(t, "{ not json")

		app := Resolve(Explicit{}, false)
		Expect(app.ClientID).To(Equal("1855"))
		Expect(app.Refreshable()).To(BeFalse())
	})
}

func TestDefaultClientID(t *testing.T) {
	RegisterTestingT(t)

	Expect(DefaultClientID(false)).To(Equal("1855"))
	Expect(DefaultClientID(true)).To(Equal("8"))
}
