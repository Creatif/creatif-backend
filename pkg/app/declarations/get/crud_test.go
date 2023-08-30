package get

import (
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Declaration node tests", func() {
	ginkgo.It("should return a text node with value queried by ID and an empty value", func() {
		name := "node"
		createdNode := testCreateBasicDeclarationTextNode("node", "modifiable")

		handler := New(NewGetNodeModel(createdNode.ID.String()))
		node, err := handler.Handle()
		gomega.Expect(err).Should(gomega.BeNil())
		gomega.Expect(node.ID).ShouldNot(gomega.BeEmpty())
		gomega.Expect(node.Name).Should(gomega.Equal(name))
		gomega.Expect(node.Value).Should(gomega.BeNil())
		gomega.Expect(node.Groups).Should(gomega.HaveLen(3))
	})

	ginkgo.It("should return a text node with value queried by name and an empty value", func() {
		name := "node"
		createdNode := testCreateBasicDeclarationTextNode("node", "modifiable")

		handler := New(NewGetNodeModel(createdNode.Name))
		view, err := handler.Handle()
		gomega.Expect(err).Should(gomega.BeNil())
		gomega.Expect(view.ID).ShouldNot(gomega.BeEmpty())
		gomega.Expect(view.Name).Should(gomega.Equal(name))
		gomega.Expect(view.Value).Should(gomega.BeNil())
	})

	ginkgo.It("should return a text node with value queried by ID and a text value", func() {
		name := "node"
		createdNode := testCreateBasicAssignmentTextNode("node")

		handler := New(NewGetNodeModel(createdNode.ID.String()))
		node, err := handler.Handle()
		gomega.Expect(err).Should(gomega.BeNil())
		gomega.Expect(node.ID).ShouldNot(gomega.BeEmpty())
		gomega.Expect(node.Name).Should(gomega.Equal(name))

		gomega.Expect(node.Value).Should(gomega.Equal("this is a text node"))
	})

	ginkgo.It("should return a text node with value queried by name and a text value", func() {
		name := "node"
		createdNode := testCreateBasicAssignmentTextNode("node")

		handler := New(NewGetNodeModel(createdNode.Name))
		node, err := handler.Handle()
		gomega.Expect(err).Should(gomega.BeNil())
		gomega.Expect(node.ID).ShouldNot(gomega.BeEmpty())
		gomega.Expect(node.Name).Should(gomega.Equal(name))

		gomega.Expect(node.Value).Should(gomega.Equal("this is a text node"))
	})

	ginkgo.It("should return a text node with value queried by ID and a boolean value", func() {
		name := "node"
		createdNode := testCreateBasicAssignmentBooleanNode("node", true)

		handler := New(NewGetNodeModel(createdNode.ID.String()))
		node, err := handler.Handle()
		gomega.Expect(err).Should(gomega.BeNil())
		gomega.Expect(node.ID).ShouldNot(gomega.BeEmpty())
		gomega.Expect(node.Name).Should(gomega.Equal(name))

		gomega.Expect(node.Value).Should(gomega.BeTrue())
	})

	ginkgo.It("should return a text node with value queried by name and a boolean value", func() {
		name := "node"
		createdNode := testCreateBasicAssignmentBooleanNode("node", true)

		handler := New(NewGetNodeModel(createdNode.Name))
		node, err := handler.Handle()
		gomega.Expect(err).Should(gomega.BeNil())
		gomega.Expect(node.ID).ShouldNot(gomega.BeEmpty())
		gomega.Expect(node.Name).Should(gomega.Equal(name))

		gomega.Expect(node.Value).Should(gomega.BeTrue())
	})
})
