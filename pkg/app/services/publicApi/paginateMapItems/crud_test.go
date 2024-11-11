package paginateMapItems

import (
	"creatif/pkg/app/auth"
	"creatif/pkg/app/services/maps/addToMap"
	"creatif/pkg/app/services/shared/queryProcessor"
	"creatif/pkg/lib/sdk"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Public API", func() {
	ginkgo.It("should get paginated list of map items", ginkgo.Label("public_api"), func() {
		projectId := testCreateProject("project")
		items, publishView := publishFullProject(projectId)

		handler := New(NewModel(publishView.Name, projectId, "paginationMap", 1, 100, "desc", "index", "", []string{}, []string{}, Options{}, []queryProcessor.Query{}), auth.NewTestingAuthentication(false, ""))
		m, err := handler.Handle()
		models := m.([]View)
		gomega.Expect(err).Should(gomega.BeNil())

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

			sdk.IncludesFn(items, func(item addToMap.View) bool {
				return item.ID == model.ID
			})
		}
	})

	ginkgo.It("should get paginated list of map items with a custom limit and fetch a different set of items", ginkgo.Label("public_api"), func() {
		projectId := testCreateProject("project")
		items, publishView := publishFullProject(projectId)

		handler := New(NewModel(publishView.Name, projectId, "paginationMap", 1, 10, "desc", "index", "", []string{}, []string{}, Options{}, []queryProcessor.Query{}), auth.NewTestingAuthentication(false, ""))
		m, err := handler.Handle()
		models := m.([]View)
		gomega.Expect(err).Should(gomega.BeNil())

		gomega.Expect(len(models)).Should(gomega.Equal(10))

		modelIds := make([]string, len(models))
		for i, model := range models {
			modelIds[i] = model.ID
			gomega.Expect(model.ProjectID).Should(gomega.Equal(projectId))
			gomega.Expect(model.Behaviour).ShouldNot(gomega.BeEmpty())
			gomega.Expect(model.Name).ShouldNot(gomega.BeEmpty())
			gomega.Expect(model.Groups).ShouldNot(gomega.BeEmpty())
			gomega.Expect(model.StructureShortID).ShouldNot(gomega.BeEmpty())
			gomega.Expect(model.StructureID).ShouldNot(gomega.BeEmpty())
			gomega.Expect(model.StructureName).ShouldNot(gomega.BeEmpty())
			gomega.Expect(model.ShortID).ShouldNot(gomega.BeEmpty())

			sdk.IncludesFn(items, func(item addToMap.View) bool {
				return item.ID == model.ID
			})
		}

		handler = New(NewModel(publishView.Name, projectId, "paginationMap", 2, 10, "desc", "index", "", []string{}, []string{}, Options{}, []queryProcessor.Query{}), auth.NewTestingAuthentication(false, ""))
		m, err = handler.Handle()
		models = m.([]View)
		gomega.Expect(err).Should(gomega.BeNil())

		for _, model := range models {
			for _, id := range modelIds {
				gomega.Expect(model.ID).ShouldNot(gomega.Equal(id))
			}
		}
	})

	ginkgo.It("should return empty result when there aren't enough items in page", ginkgo.Label("public_api"), func() {
		projectId := testCreateProject("project")
		_, publishView := publishFullProject(projectId)

		handler := New(NewModel(publishView.Name, projectId, "paginationMap", 3, 100, "desc", "index", "", []string{}, []string{}, Options{}, []queryProcessor.Query{}), auth.NewTestingAuthentication(false, ""))
		m, err := handler.Handle()
		models := m.([]View)
		gomega.Expect(err).Should(gomega.BeNil())

		gomega.Expect(len(models)).Should(gomega.Equal(0))
	})

	ginkgo.It("should get paginated list of map items based on group", ginkgo.Label("public_api"), func() {
		ginkgo.Skip("")

		projectId := testCreateProject("project")
		items, publishView := publishFullProject(projectId)

		handler := New(NewModel(publishView.Name, projectId, "paginationMap", 1, 100, "desc", "index", "", []string{}, []string{"group-0"}, Options{}, []queryProcessor.Query{}), auth.NewTestingAuthentication(false, ""))
		m, err := handler.Handle()
		models := m.([]View)
		gomega.Expect(err).Should(gomega.BeNil())

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

			sdk.IncludesFn(items, func(item addToMap.View) bool {
				return item.ID == model.ID
			})
		}
	})

	ginkgo.It("should get paginated list of map items based on group and locale", ginkgo.Label("public_api"), func() {
		ginkgo.Skip("")

		projectId := testCreateProject("project")
		items, publishView := publishFullProject(projectId)

		handler := New(NewModel(publishView.Name, projectId, "paginationMap", 1, 100, "desc", "index", "", []string{"eng"}, []string{"group-0"}, Options{}, []queryProcessor.Query{}), auth.NewTestingAuthentication(false, ""))
		m, err := handler.Handle()
		models := m.([]View)
		gomega.Expect(err).Should(gomega.BeNil())

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

			sdk.IncludesFn(items, func(item addToMap.View) bool {
				return item.ID == model.ID
			})
		}
	})

	ginkgo.It("should get paginated list of map items based on group, locale and search", ginkgo.Label("public_api"), func() {
		ginkgo.Skip("")

		projectId := testCreateProject("project")
		items, publishView := publishFullProject(projectId)

		handler := New(NewModel(publishView.Name, projectId, "paginationMap", 1, 100, "desc", "index", "0", []string{"eng"}, []string{"group-0"}, Options{}, []queryProcessor.Query{}), auth.NewTestingAuthentication(false, ""))
		m, err := handler.Handle()
		models := m.([]View)
		gomega.Expect(err).Should(gomega.BeNil())
		gomega.Expect(len(models)).Should(gomega.Equal(29))

		for _, model := range models {
			gomega.Expect(model.ProjectID).Should(gomega.Equal(projectId))
			gomega.Expect(model.Behaviour).ShouldNot(gomega.BeEmpty())
			gomega.Expect(model.Name).ShouldNot(gomega.BeEmpty())
			gomega.Expect(model.Groups).ShouldNot(gomega.BeEmpty())
			gomega.Expect(model.StructureShortID).ShouldNot(gomega.BeEmpty())
			gomega.Expect(model.StructureID).ShouldNot(gomega.BeEmpty())
			gomega.Expect(model.StructureName).ShouldNot(gomega.BeEmpty())
			gomega.Expect(model.ShortID).ShouldNot(gomega.BeEmpty())

			sdk.IncludesFn(items, func(item addToMap.View) bool {
				return item.ID == model.ID
			})
		}
	})
})
