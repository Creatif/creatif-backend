package getListItemById

import (
	"creatif/pkg/app/auth"
	"creatif/pkg/lib/logger"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Public API", func() {
	ginkgo.It("should get public list item by id (getListItemById)", ginkgo.Label("public_api"), func() {
		projectId := testCreateProject("project")
		mapItem, _ := publishFullProject(projectId)

		handler := New(NewModel("", projectId, mapItem.ID, Options{}), auth.NewTestingAuthentication(false, ""), logger.NewLogBuilder())
		model, err := handler.Handle()
		gomega.Expect(err).Should(gomega.BeNil())
		item := model.(View)

		gomega.Expect(item.ProjectID).Should(gomega.Equal(projectId))
		gomega.Expect(item.Behaviour).ShouldNot(gomega.BeEmpty())
		gomega.Expect(item.Name).ShouldNot(gomega.BeEmpty())
		gomega.Expect(item.Groups).ShouldNot(gomega.BeEmpty())
		gomega.Expect(item.StructureShortID).ShouldNot(gomega.BeEmpty())
		gomega.Expect(item.StructureID).ShouldNot(gomega.BeEmpty())
		gomega.Expect(item.StructureName).ShouldNot(gomega.BeEmpty())
		gomega.Expect(item.ShortID).ShouldNot(gomega.BeEmpty())

		gomega.Expect(item.ID).Should(gomega.Equal(mapItem.ID))
	})
})
