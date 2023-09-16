package getValue

import (
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("GET value of declaration variable", func() {
	ginkgo.It("should return a variable value", func() {
		createdVariable := testCreateDeclarationVariable("variable", "modifiable")

		handler := New(NewModel(createdVariable.Name))
		value, err := handler.Handle()
		gomega.Expect(err).Should(gomega.BeNil())

		gomega.Expect(value).ShouldNot(gomega.BeEmpty())
	})
})
