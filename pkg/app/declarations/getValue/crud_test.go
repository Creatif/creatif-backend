package getValue

import (
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("GET value of declaration node", func() {
	ginkgo.It("should return a text node with value queried by ID and an empty value", func() {
		createdNode := testCreateBasicDeclarationTextNode("node", "modifiable")

		handler := New(NewModel(createdNode.Name))
		value, err := handler.Handle()
		gomega.Expect(err).Should(gomega.BeNil())

		gomega.Expect(value).Should(gomega.BeEmpty())
	})

	ginkgo.It("should return a text node with value queried by name and an empty value", func() {
		createdNode := testCreateBasicDeclarationTextNode("node", "modifiable")

		handler := New(NewModel(createdNode.Name))
		value, err := handler.Handle()
		gomega.Expect(err).Should(gomega.BeNil())

		gomega.Expect(value).Should(gomega.BeEmpty())
	})

	ginkgo.It("should return a text node with value queried by ID and a text value", func() {
		createdNode := testCreateBasicAssignmentTextNode("node")

		handler := New(NewModel(createdNode.Name))
		value, err := handler.Handle()
		testAssertErrNil(err)

		gomega.Expect(value).ShouldNot(gomega.BeEmpty())
	})
})
