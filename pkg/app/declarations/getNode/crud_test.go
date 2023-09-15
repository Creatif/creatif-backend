package getNode

import (
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Declaration node tests", func() {
	ginkgo.It("should return a text node with value queried by ID and an empty value", func() {
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

		gomega.Expect(node["value"]).Should(gomega.BeEmpty())
	})

	ginkgo.It("should return a text node with value queried by name and an empty value", func() {
		name := "node"
		createdNode := testCreateBasicDeclarationTextNode("node", "modifiable")

		handler := New(NewModel(createdNode.Name, []string{"value"}))
		view, err := handler.Handle()
		gomega.Expect(err).Should(gomega.BeNil())

		gomega.Expect(view).Should(gomega.HaveKey("id"))
		gomega.Expect(view).Should(gomega.HaveKey("name"))
		gomega.Expect(view).Should(gomega.HaveKey("value"))

		gomega.Expect(view["id"]).ShouldNot(gomega.BeEmpty())
		gomega.Expect(view["name"]).Should(gomega.Equal(name))

		gomega.Expect(view["value"]).Should(gomega.BeEmpty())
	})

	ginkgo.It("should return a text node with value queried by ID and a text value", func() {
		name := "node"
		createdNode := testCreateBasicAssignmentTextNode("node")

		handler := New(NewModel(createdNode.Name, []string{"value", "behaviour"}))
		view, err := handler.Handle()
		testAssertErrNil(err)

		gomega.Expect(view).Should(gomega.HaveKey("id"))
		gomega.Expect(view).Should(gomega.HaveKey("name"))
		gomega.Expect(view).Should(gomega.HaveKey("value"))
		gomega.Expect(view).Should(gomega.HaveKey("behaviour"))

		gomega.Expect(view).ShouldNot(gomega.HaveKey("groups"))

		gomega.Expect(view["id"]).ShouldNot(gomega.BeEmpty())
		gomega.Expect(view["name"]).Should(gomega.Equal(name))

		gomega.Expect(view["value"]).ShouldNot(gomega.BeEmpty())
	})
})
