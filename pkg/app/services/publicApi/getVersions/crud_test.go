package getVersions

import (
	"creatif/pkg/app/auth"
	"creatif/pkg/lib/logger"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Public API", func() {
	ginkgo.It("should get all version of a published project", ginkgo.Label("public_api"), func() {
		projectId := testCreateProject("project")
		publishFullProject(projectId)

		handler := New(NewModel(projectId), auth.NewTestingAuthentication(false, ""), logger.NewLogBuilder())
		model, err := handler.Handle()
		gomega.Expect(err).Should(gomega.BeNil())
		gomega.Expect(len(model)).Should(gomega.Equal(1))

		for _, version := range model {
			gomega.Expect(version.ID).ShouldNot(gomega.BeEmpty())
			gomega.Expect(version.ProjectID).Should(gomega.Equal(projectId))
		}
	})
})
