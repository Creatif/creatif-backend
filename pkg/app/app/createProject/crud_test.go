package createProject

import (
	"creatif/pkg/app/auth"
	"creatif/pkg/lib/logger"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Create project tests", func() {
	ginkgo.It("should create a new project", func() {
		handler := New(NewModel("project name"), auth.NewNoopAuthentication(), logger.NewLogBuilder())

		model, err := handler.Handle()
		testAssertErrNil(err)
		testAssertIDValid(model.ID)

		gomega.Expect(model.Name).Should(gomega.Equal("project name"))
	})

	ginkgo.It("should fail if project already exists", func() {
		testCreateProject("project name")
		handler := New(NewModel("project name"), auth.NewNoopAuthentication(), logger.NewLogBuilder())

		_, err := handler.Handle()

		gomega.Expect(err).ShouldNot(gomega.BeNil())
	})
})
