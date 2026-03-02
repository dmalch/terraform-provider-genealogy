package internal

import (
	"os"
	"path"
	"testing"

	. "github.com/onsi/gomega"
)

func TestTokenCacheFilePath(t *testing.T) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("failed to get home dir: %v", err)
	}

	t.Run("production environment", func(t *testing.T) {
		RegisterTestingT(t)

		result, err := tokenCacheFilePath(false)

		Expect(err).ToNot(HaveOccurred())
		Expect(result).To(Equal(path.Join(homeDir, ".genealogy", "geni_token.json")))
	})

	t.Run("sandbox environment", func(t *testing.T) {
		RegisterTestingT(t)

		result, err := tokenCacheFilePath(true)

		Expect(err).ToNot(HaveOccurred())
		Expect(result).To(Equal(path.Join(homeDir, ".genealogy", "geni_sandbox_token.json")))
	})
}

func TestClientId(t *testing.T) {
	t.Run("production", func(t *testing.T) {
		RegisterTestingT(t)

		Expect(clientId(false)).To(Equal("1855"))
	})

	t.Run("sandbox", func(t *testing.T) {
		RegisterTestingT(t)

		Expect(clientId(true)).To(Equal("8"))
	})
}
