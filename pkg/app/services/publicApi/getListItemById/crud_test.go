package getListItemById

import (
	"creatif/pkg/app/auth"
	"creatif/pkg/lib/logger"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Public API", func() {
	ginkgo.It("should get public list item by id (getListItemById)", func() {
		projectId := testCreateProject("project")
		mapItem, version := publishFullProject(projectId)

		handler := New(NewModel(projectId, mapItem.ID, version.Name), auth.NewTestingAuthentication(false, ""), logger.NewLogBuilder())
		model, err := handler.Handle()
		gomega.Expect(err).Should(gomega.BeNil())

		gomega.Expect(model.ProjectID).Should(gomega.Equal(projectId))
		gomega.Expect(model.Behaviour).ShouldNot(gomega.BeEmpty())
		gomega.Expect(model.ItemName).ShouldNot(gomega.BeEmpty())
		gomega.Expect(model.Groups).ShouldNot(gomega.BeEmpty())
		gomega.Expect(model.StructureShortID).ShouldNot(gomega.BeEmpty())
		gomega.Expect(model.StructureID).ShouldNot(gomega.BeEmpty())
		gomega.Expect(model.StructureName).ShouldNot(gomega.BeEmpty())
		gomega.Expect(model.ItemShortID).ShouldNot(gomega.BeEmpty())

		gomega.Expect(model.ItemID).Should(gomega.Equal(mapItem.ID))
		gomega.Expect(len(model.Connections)).Should(gomega.Equal(4))
	})
})
