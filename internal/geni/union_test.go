package geni

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestGetUnion1(t *testing.T) {
	t.Skip()
	RegisterTestingT(t)

	unionId := "union-1838"

	client, err := NewClient(testAccessToken, true)
	Expect(err).ToNot(HaveOccurred())

	union, err := client.GetUnion(unionId)

	Expect(err).ToNot(HaveOccurred())
	Expect(union.Id).To(BeEquivalentTo(unionId))
}
