package geni

import (
	"testing"

	. "github.com/onsi/gomega"
	"golang.org/x/oauth2"
)

func TestGetUnion1(t *testing.T) {
	t.Skip()
	RegisterTestingT(t)

	unionId := "union-1838"

	client := NewClient(oauth2.StaticTokenSource(&oauth2.Token{AccessToken: testAccessToken}), true)

	union, err := client.GetUnion(unionId)

	Expect(err).ToNot(HaveOccurred())
	Expect(union.Id).To(BeEquivalentTo(unionId))
}
