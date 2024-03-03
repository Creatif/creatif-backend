package paginateListItems

import (
	"creatif/pkg/app/auth"
	"creatif/pkg/app/services/lists/addToList"
	"creatif/pkg/lib/logger"
	"creatif/pkg/lib/sdk"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Public API", func() {
	ginkgo.It("should get paginated list of list items", func() {
		projectId := testCreateProject("project")
		items, version := publishFullProject(projectId)

		handler := New(NewModel(projectId, version.Name, "paginationList", 1, "desc", "index", "", []string{}, []string{}), auth.NewTestingAuthentication(false, ""), logger.NewLogBuilder())
		models, err := handler.Handle()
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

			sdk.IncludesFn(items, func(item addToList.View) bool {
				return item.ID == model.ID
			})
			gomega.Expect(len(model.Connections)).Should(gomega.Equal(4))
		}
	})

	ginkgo.It("should return empty result when there aren't enough items in page", func() {
		projectId := testCreateProject("project")
		_, version := publishFullProject(projectId)

		handler := New(NewModel(projectId, version.Name, "paginationList", 3, "desc", "index", "", []string{}, []string{}), auth.NewTestingAuthentication(false, ""), logger.NewLogBuilder())
		models, err := handler.Handle()
		gomega.Expect(err).Should(gomega.BeNil())

		gomega.Expect(len(models)).Should(gomega.Equal(0))
	})

	ginkgo.It("should get paginated list of list items based on group", func() {
		projectId := testCreateProject("project")
		items, version := publishFullProject(projectId)

		handler := New(NewModel(projectId, version.Name, "paginationList", 1, "desc", "index", "", []string{}, []string{"group-0"}), auth.NewTestingAuthentication(false, ""), logger.NewLogBuilder())
		models, err := handler.Handle()
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

			sdk.IncludesFn(items, func(item addToList.View) bool {
				return item.ID == model.ID
			})
			gomega.Expect(len(model.Connections)).Should(gomega.Equal(4))
		}
	})

	ginkgo.It("should get paginated list of list items based on group and locale", func() {
		projectId := testCreateProject("project")
		items, version := publishFullProject(projectId)

		handler := New(NewModel(projectId, version.Name, "paginationList", 1, "desc", "index", "", []string{"eng"}, []string{"group-0"}), auth.NewTestingAuthentication(false, ""), logger.NewLogBuilder())
		models, err := handler.Handle()
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

			sdk.IncludesFn(items, func(item addToList.View) bool {
				return item.ID == model.ID
			})
			gomega.Expect(len(model.Connections)).Should(gomega.Equal(4))
		}
	})

	ginkgo.It("should get paginated list of list items based on group, locale and search", func() {
		projectId := testCreateProject("project")
		items, version := publishFullProject(projectId)

		handler := New(NewModel(projectId, version.Name, "paginationList", 1, "desc", "index", "0", []string{"eng"}, []string{"group-0"}), auth.NewTestingAuthentication(false, ""), logger.NewLogBuilder())
		models, err := handler.Handle()
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

			sdk.IncludesFn(items, func(item addToList.View) bool {
				return item.ID == model.ID
			})
			gomega.Expect(len(model.Connections)).Should(gomega.Equal(4))
		}
	})
})