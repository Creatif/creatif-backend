package getManyItems

import (
	"creatif/pkg/app/auth"
	"creatif/pkg/lib/logger"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Public API", func() {
	ginkgo.It("should get many items from an array of items (getManyItems)", ginkgo.Label("public_api"), func() {
		projectId := testCreateProject("project")
		listIds, mapIds, version := publishFullProject(projectId)

		part1 := listIds[0:50]
		part2 := mapIds[0:50]
		ids := make([]string, 0)
		ids = append(ids, part1...)
		ids = append(ids, part2...)

		handler := New(NewModel(version.Name, projectId, ids, Options{}), auth.NewTestingAuthentication(false, ""), logger.NewLogBuilder())
		m, err := handler.Handle()
		gomega.Expect(err).Should(gomega.BeNil())
		models := m.([]View)

		gomega.Expect(len(models)).Should(gomega.Equal(100))

		for _, model := range models {
			gomega.Expect(model.ProjectID).Should(gomega.Equal(projectId))
			gomega.Expect(model.Behaviour).ShouldNot(gomega.BeEmpty())
			gomega.Expect(model.Name).ShouldNot(gomega.BeEmpty())
			gomega.Expect(model.Groups).ShouldNot(gomega.BeEmpty())
			gomega.Expect(model.StructureShortID).ShouldNot(gomega.BeEmpty())
			gomega.Expect(model.StructureID).ShouldNot(gomega.BeEmpty())
			gomega.Expect(model.StructureName).ShouldNot(gomega.BeEmpty())
			gomega.Expect(model.ShortID).ShouldNot(gomega.BeEmpty())
		}
	})
})
