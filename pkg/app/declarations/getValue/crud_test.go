package getValue

import (
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("GET value of declaration node", func() {
	ginkgo.It("should return a node value", func() {
		createdNode := testCreateDeclarationNode("node", "modifiable")

		handler := New(NewModel(createdNode.Name))
		value, err := handler.Handle()
		gomega.Expect(err).Should(gomega.BeNil())

		gomega.Expect(value).ShouldNot(gomega.BeEmpty())
	})
})
