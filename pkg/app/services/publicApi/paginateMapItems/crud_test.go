package paginateMapItems

import (
	"creatif/pkg/app/auth"
	"creatif/pkg/app/services/maps/addToMap"
	"creatif/pkg/lib/sdk"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Public API", func() {
	ginkgo.It("should get paginated list of map items", ginkgo.Label("public_api"), func() {
		projectId := testCreateProject("project")
		items, _ := publishFullProject(projectId)

		handler := New(NewModel("", projectId, "paginationMap", 1, "desc", "index", "", []string{}, []string{}, Options{}), auth.NewTestingAuthentication(false, ""))
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

	ginkgo.It("should return empty result when there aren't enough items in page", ginkgo.Label("public_api"), func() {
		projectId := testCreateProject("project")
		publishFullProject(projectId)

		handler := New(NewModel("", projectId, "paginationMap", 3, "desc", "index", "", []string{}, []string{}, Options{}), auth.NewTestingAuthentication(false, ""))
		m, err := handler.Handle()
		models := m.([]View)
		gomega.Expect(err).Should(gomega.BeNil())

		gomega.Expect(len(models)).Should(gomega.Equal(0))
	})

	ginkgo.It("should get paginated list of map items based on group", ginkgo.Label("public_api"), func() {
		projectId := testCreateProject("project")
		items, _ := publishFullProject(projectId)

		handler := New(NewModel("", projectId, "paginationMap", 1, "desc", "index", "", []string{}, []string{"group-0"}, Options{}), auth.NewTestingAuthentication(false, ""))
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
		projectId := testCreateProject("project")
		items, _ := publishFullProject(projectId)

		handler := New(NewModel("", projectId, "paginationMap", 1, "desc", "index", "", []string{"eng"}, []string{"group-0"}, Options{}), auth.NewTestingAuthentication(false, ""))
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
		projectId := testCreateProject("project")
		items, _ := publishFullProject(projectId)

		handler := New(NewModel("", projectId, "paginationMap", 1, "desc", "index", "0", []string{"eng"}, []string{"group-0"}, Options{}), auth.NewTestingAuthentication(false, ""))
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
