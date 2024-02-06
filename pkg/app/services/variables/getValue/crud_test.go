package getValue

import (
	"creatif/pkg/app/auth"
	"creatif/pkg/lib/logger"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("GET value of declaration variable", func() {
	ginkgo.It("should return a variable value by name", func() {
		projectId := testCreateProject("project")
		testCreateGroups(projectId, []string{"one", "two", "three"})
		createdVariable := testCreateDeclarationVariable(projectId, "variable", "modifiable")

		handler := New(NewModel(projectId, "", "", createdVariable.Name, "eng"), auth.NewTestingAuthentication(false, ""), logger.NewLogBuilder())
		value, err := handler.Handle()
		gomega.Expect(err).Should(gomega.BeNil())

		gomega.Expect(value).ShouldNot(gomega.BeEmpty())
	})

	ginkgo.It("should return a variable value by id", func() {
		projectId := testCreateProject("project")
		testCreateGroups(projectId, []string{"one", "two", "three"})
		createdVariable := testCreateDeclarationVariable(projectId, "variable", "modifiable")

		handler := New(NewModel(projectId, createdVariable.ID, "", "", "eng"), auth.NewTestingAuthentication(false, ""), logger.NewLogBuilder())
		value, err := handler.Handle()
		gomega.Expect(err).Should(gomega.BeNil())

		gomega.Expect(value).ShouldNot(gomega.BeEmpty())
	})

	ginkgo.It("should return a variable value by shortId", func() {
		projectId := testCreateProject("project")
		testCreateGroups(projectId, []string{"one", "two", "three"})
		createdVariable := testCreateDeclarationVariable(projectId, "variable", "modifiable")

		handler := New(NewModel(projectId, "", createdVariable.ShortID, "", "eng"), auth.NewTestingAuthentication(false, ""), logger.NewLogBuilder())
		value, err := handler.Handle()
		gomega.Expect(err).Should(gomega.BeNil())

		gomega.Expect(value).ShouldNot(gomega.BeEmpty())
	})
})
