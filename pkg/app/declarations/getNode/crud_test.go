package getNode

import (
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Declaration node tests", func() {
	ginkgo.It("should return a text node with value", func() {
		name := "node"
		createdNode := testCreateBasicDeclarationTextNode(name, "modifiable")

		handler := New(NewModel(createdNode.Name, []string{}))
		node, err := handler.Handle()
		gomega.Expect(err).Should(gomega.BeNil())

		gomega.Expect(node).Should(gomega.HaveKey("id"))
		gomega.Expect(node).Should(gomega.HaveKey("name"))
		gomega.Expect(node).Should(gomega.HaveKey("behaviour"))
		gomega.Expect(node).Should(gomega.HaveKey("metadata"))
		gomega.Expect(node).Should(gomega.HaveKey("groups"))
		gomega.Expect(node).Should(gomega.HaveKey("createdAt"))
		gomega.Expect(node).Should(gomega.HaveKey("updatedAt"))

		gomega.Expect(node["id"]).ShouldNot(gomega.BeEmpty())
		gomega.Expect(node["name"]).Should(gomega.Equal(name))

		gomega.Expect(node["value"]).ShouldNot(gomega.BeEmpty())
	})
})
