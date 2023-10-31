package getVariable

import (
	"creatif/pkg/app/auth"
	"creatif/pkg/lib/logger"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Declaration variable tests", func() {
	ginkgo.It("should return a text variable with value with name", func() {
		projectId := testCreateProject("project")
		name := "variable"
		createdVariable := testCreateBasicDeclarationTextVariable(projectId, name, "modifiable")

		handler := New(NewModel(projectId, createdVariable.Name, "", "", "eng", []string{}), auth.NewNoopAuthentication(), logger.NewLogBuilder())
		variable, err := handler.Handle()
		gomega.Expect(err).Should(gomega.BeNil())

		gomega.Expect(variable).Should(gomega.HaveKey("id"))
		gomega.Expect(variable).Should(gomega.HaveKey("name"))
		gomega.Expect(variable).Should(gomega.HaveKey("behaviour"))
		gomega.Expect(variable).Should(gomega.HaveKey("metadata"))
		gomega.Expect(variable).Should(gomega.HaveKey("groups"))
		gomega.Expect(variable).Should(gomega.HaveKey("createdAt"))
		gomega.Expect(variable).Should(gomega.HaveKey("updatedAt"))
		gomega.Expect(variable).Should(gomega.HaveKey("projectID"))
		gomega.Expect(variable).Should(gomega.HaveKey("locale"))

		gomega.Expect(variable["id"]).ShouldNot(gomega.BeEmpty())
		gomega.Expect(variable["projectID"]).ShouldNot(gomega.BeEmpty())
		gomega.Expect(variable["name"]).Should(gomega.Equal(name))

		gomega.Expect(variable["value"]).ShouldNot(gomega.BeEmpty())
	})

	ginkgo.It("should return a text variable with value with id", func() {
		projectId := testCreateProject("project")
		name := "variable"
		createdVariable := testCreateBasicDeclarationTextVariable(projectId, name, "modifiable")

		handler := New(NewModel(projectId, "", createdVariable.ID, "", "eng", []string{}), auth.NewNoopAuthentication(), logger.NewLogBuilder())
		variable, err := handler.Handle()
		gomega.Expect(err).Should(gomega.BeNil())

		gomega.Expect(variable).Should(gomega.HaveKey("id"))
		gomega.Expect(variable).Should(gomega.HaveKey("name"))
		gomega.Expect(variable).Should(gomega.HaveKey("behaviour"))
		gomega.Expect(variable).Should(gomega.HaveKey("metadata"))
		gomega.Expect(variable).Should(gomega.HaveKey("groups"))
		gomega.Expect(variable).Should(gomega.HaveKey("createdAt"))
		gomega.Expect(variable).Should(gomega.HaveKey("updatedAt"))
		gomega.Expect(variable).Should(gomega.HaveKey("projectID"))
		gomega.Expect(variable).Should(gomega.HaveKey("locale"))

		gomega.Expect(variable["id"]).ShouldNot(gomega.BeEmpty())
		gomega.Expect(variable["projectID"]).ShouldNot(gomega.BeEmpty())
		gomega.Expect(variable["name"]).Should(gomega.Equal(name))

		gomega.Expect(variable["value"]).ShouldNot(gomega.BeEmpty())
	})

	ginkgo.It("should return a text variable with value with shortID", func() {
		projectId := testCreateProject("project")
		name := "variable"
		createdVariable := testCreateBasicDeclarationTextVariable(projectId, name, "modifiable")

		handler := New(NewModel(projectId, "", "", createdVariable.ShortID, "eng", []string{}), auth.NewNoopAuthentication(), logger.NewLogBuilder())
		variable, err := handler.Handle()
		gomega.Expect(err).Should(gomega.BeNil())

		gomega.Expect(variable).Should(gomega.HaveKey("id"))
		gomega.Expect(variable).Should(gomega.HaveKey("name"))
		gomega.Expect(variable).Should(gomega.HaveKey("behaviour"))
		gomega.Expect(variable).Should(gomega.HaveKey("metadata"))
		gomega.Expect(variable).Should(gomega.HaveKey("groups"))
		gomega.Expect(variable).Should(gomega.HaveKey("createdAt"))
		gomega.Expect(variable).Should(gomega.HaveKey("updatedAt"))
		gomega.Expect(variable).Should(gomega.HaveKey("projectID"))
		gomega.Expect(variable).Should(gomega.HaveKey("locale"))

		gomega.Expect(variable["id"]).ShouldNot(gomega.BeEmpty())
		gomega.Expect(variable["projectID"]).ShouldNot(gomega.BeEmpty())
		gomega.Expect(variable["name"]).Should(gomega.Equal(name))

		gomega.Expect(variable["value"]).ShouldNot(gomega.BeEmpty())
	})
})
