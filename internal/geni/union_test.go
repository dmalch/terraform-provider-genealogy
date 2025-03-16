package geni

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestGetUnion1(t *testing.T) {
	t.Skip()
	RegisterTestingT(t)

	unionId := "union-1838"
	union, err := GetUnion(testAccessToken, unionId)

	Expect(err).ToNot(HaveOccurred())
	Expect(union.Id).To(BeEquivalentTo(unionId))
}
