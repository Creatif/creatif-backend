package getVariable

import (
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Declaration variable tests", func() {
	ginkgo.It("should return a text variable with value", func() {
		name := "variable"
		createdVariable := testCreateBasicDeclarationTextVariable(name, "modifiable")

		handler := New(NewModel(createdVariable.Name, []string{}))
		variable, err := handler.Handle()
		gomega.Expect(err).Should(gomega.BeNil())

		gomega.Expect(variable).Should(gomega.HaveKey("id"))
		gomega.Expect(variable).Should(gomega.HaveKey("name"))
		gomega.Expect(variable).Should(gomega.HaveKey("behaviour"))
		gomega.Expect(variable).Should(gomega.HaveKey("metadata"))
		gomega.Expect(variable).Should(gomega.HaveKey("groups"))
		gomega.Expect(variable).Should(gomega.HaveKey("createdAt"))
		gomega.Expect(variable).Should(gomega.HaveKey("updatedAt"))

		gomega.Expect(variable["id"]).ShouldNot(gomega.BeEmpty())
		gomega.Expect(variable["name"]).Should(gomega.Equal(name))

		gomega.Expect(variable["value"]).ShouldNot(gomega.BeEmpty())
	})
})
