package getMapItemByName

import (
	"creatif/pkg/app/auth"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Public API", func() {
	ginkgo.It("should get public map item by name and default locale (getMapItemByName)", ginkgo.Label("public_api"), func() {
		projectId := testCreateProject("project")
		mapItem, structure, publishView := publishFullProject(projectId)

		handler := New(NewModel(publishView.Name, projectId, structure.Name, mapItem.Name, "eng", Options{}), auth.NewTestingAuthentication(false, ""))
		m, err := handler.Handle()
		gomega.Expect(err).Should(gomega.BeNil())
		model := m.(View)

		gomega.Expect(model.ProjectID).Should(gomega.Equal(projectId))
		gomega.Expect(model.Behaviour).ShouldNot(gomega.BeEmpty())
		gomega.Expect(model.Name).ShouldNot(gomega.BeEmpty())
		gomega.Expect(model.Groups).ShouldNot(gomega.BeEmpty())
		gomega.Expect(model.StructureShortID).ShouldNot(gomega.BeEmpty())
		gomega.Expect(model.StructureID).ShouldNot(gomega.BeEmpty())
		gomega.Expect(model.StructureName).ShouldNot(gomega.BeEmpty())
		gomega.Expect(model.ShortID).ShouldNot(gomega.BeEmpty())

		gomega.Expect(model.ID).Should(gomega.Equal(mapItem.ID))
	})

	ginkgo.It("should get public map item by name and eng locale (getMapItemByName)", ginkgo.Label("public_api"), func() {
		projectId := testCreateProject("project")
		mapItem, structure, publishView := publishFullProject(projectId)

		handler := New(NewModel(publishView.Name, projectId, structure.Name, mapItem.Name, "eng", Options{}), auth.NewTestingAuthentication(false, ""))
		m, err := handler.Handle()
		gomega.Expect(err).Should(gomega.BeNil())
		model := m.(View)

		gomega.Expect(model.ProjectID).Should(gomega.Equal(projectId))
		gomega.Expect(model.Behaviour).ShouldNot(gomega.BeEmpty())
		gomega.Expect(model.Name).ShouldNot(gomega.BeEmpty())
		gomega.Expect(model.Groups).ShouldNot(gomega.BeEmpty())
		gomega.Expect(model.StructureShortID).ShouldNot(gomega.BeEmpty())
		gomega.Expect(model.StructureID).ShouldNot(gomega.BeEmpty())
		gomega.Expect(model.StructureName).ShouldNot(gomega.BeEmpty())
		gomega.Expect(model.ShortID).ShouldNot(gomega.BeEmpty())

		gomega.Expect(model.ID).Should(gomega.Equal(mapItem.ID))
	})
})
