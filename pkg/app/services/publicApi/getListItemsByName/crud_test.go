package getListItemsByName

import (
	"creatif/pkg/app/auth"
	"creatif/pkg/lib/logger"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Public API", func() {
	ginkgo.It("should get public list items by name and default locale (getListItemsByName)", ginkgo.Label("public_api"), func() {
		projectId := testCreateProject("project")
		item, structure, _ := publishFullProject(projectId)

		handler := New(NewModel("", projectId, structure.Name, item.Name, "eng", Options{}), auth.NewTestingAuthentication(false, ""), logger.NewLogBuilder())
		m, err := handler.Handle()
		gomega.Expect(err).Should(gomega.BeNil())
		models := m.([]View)

		gomega.Expect(len(models)).Should(gomega.Equal(1))

		for _, model := range models {
			gomega.Expect(model.ProjectID).Should(gomega.Equal(projectId))
			gomega.Expect(model.Behaviour).ShouldNot(gomega.BeEmpty())
			gomega.Expect(model.Name).ShouldNot(gomega.BeEmpty())
			gomega.Expect(model.Groups).ShouldNot(gomega.BeEmpty())
			gomega.Expect(model.StructureShortID).ShouldNot(gomega.BeEmpty())
			gomega.Expect(model.StructureID).ShouldNot(gomega.BeEmpty())
			gomega.Expect(model.StructureName).ShouldNot(gomega.BeEmpty())
			gomega.Expect(model.ShortID).ShouldNot(gomega.BeEmpty())

			gomega.Expect(model.ID).Should(gomega.Equal(item.ID))
		}
	})
})
