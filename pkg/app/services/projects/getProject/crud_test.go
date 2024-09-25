package getProject

import (
	"creatif/pkg/app/auth"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Get project tests", func() {
	ginkgo.It("should get a project by ID", func() {
		projectId := testCreateProject("project")
		handler := New(NewModel(projectId), auth.NewTestingAuthentication(true, ""))

		model, err := handler.Handle()
		testAssertErrNil(err)
		testAssertIDValid(model.ID)

		gomega.Expect(model.Name).Should(gomega.Equal("project"))
		gomega.Expect(model.ID).ShouldNot(gomega.BeEmpty())
	})
})
