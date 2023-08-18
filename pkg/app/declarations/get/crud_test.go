package create

import (
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Declaration node tests", func() {
	ginkgo.It("should return a text node with value queried by ID", func() {
		name := "node"
		createdNode := testCreateBasicDeclarationTextNode("node", "modifiable")

		handler := New(NewGetNodeModel(createdNode.ID))
		node, err := handler.Handle()
		gomega.Expect(err).Should(gomega.BeNil())
		gomega.Expect(node.ID).ShouldNot(gomega.BeEmpty())
		gomega.Expect(node.Name).Should(gomega.Equal(name))
	})

	ginkgo.It("should return a text node with value queried by name", func() {
		name := "node"
		createdNode := testCreateBasicDeclarationTextNode("node", "modifiable")

		handler := New(NewGetNodeModel(createdNode.Name))
		node, err := handler.Handle()
		gomega.Expect(err).Should(gomega.BeNil())
		gomega.Expect(node.ID).ShouldNot(gomega.BeEmpty())
		gomega.Expect(node.Name).Should(gomega.Equal(name))
	})
})
