package getValue

import (
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("GET value of declaration variable", func() {
	ginkgo.It("should return a variable value", func() {
		projectId := testCreateProject("project")
		createdVariable := testCreateDeclarationVariable(projectId, "variable", "modifiable")

		handler := New(NewModel(projectId, createdVariable.Name, "eng"))
		value, err := handler.Handle()
		gomega.Expect(err).Should(gomega.BeNil())

		gomega.Expect(value).ShouldNot(gomega.BeEmpty())
	})
})
