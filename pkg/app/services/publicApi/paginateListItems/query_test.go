package paginateListItems

import (
	"creatif/pkg/app/auth"
	"creatif/pkg/app/services/lists/addToList"
	"creatif/pkg/app/services/shared/queryProcessor"
	"creatif/pkg/lib/sdk"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Public API list query tests", func() {
	ginkgo.It("should return a paginated list of map items based on EQUAL query, the result should not be empty", ginkgo.Label("public_api", "list_query"), func() {
		projectId := testCreateProject("project")
		items, publishView := publishFullProject(projectId)

		handler := New(NewModel(publishView.Name, projectId, "paginationList", 1, 100, "desc", "index", "", []string{}, []string{}, Options{}, []queryProcessor.Query{
			{
				Column:   "one",
				Value:    "one",
				Operator: "equal",
				Type:     "string",
			},
			{
				Column:   "four",
				Value:    "453",
				Operator: "equal",
				Type:     "int",
			},
		}), auth.NewTestingAuthentication(false, ""))
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

			sdk.IncludesFn(items, func(item addToList.View) bool {
				return item.ID == model.ID
			})
		}
	})

	ginkgo.It("should return a paginated list of map items based on EQUAL query, the result should not be empty and decimal should be used", ginkgo.Label("public_api", "list_query"), func() {
		projectId := testCreateProject("project")
		items, publishView := publishFullProject(projectId)

		handler := New(NewModel(publishView.Name, projectId, "paginationList", 1, 100, "desc", "index", "", []string{}, []string{}, Options{}, []queryProcessor.Query{
			{
				Column:   "one",
				Value:    "one",
				Operator: "equal",
				Type:     "string",
			},
			{
				Column:   "five",
				Value:    "456.43",
				Operator: "equal",
				Type:     "float",
			},
		}), auth.NewTestingAuthentication(false, ""))
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

			sdk.IncludesFn(items, func(item addToList.View) bool {
				return item.ID == model.ID
			})
		}
	})

	ginkgo.It("should return a paginated list of list items based on UNEQUAL query, the result should be empty", ginkgo.Label("public_api", "list_query"), func() {
		projectId := testCreateProject("project")
		_, publishView := publishFullProject(projectId)

		handler := New(NewModel(publishView.Name, projectId, "paginationList", 1, 100, "desc", "index", "", []string{}, []string{}, Options{}, []queryProcessor.Query{
			{
				Column:   "one",
				Value:    "one",
				Operator: "unequal",
				Type:     "string",
			},
			{
				Column:   "four",
				Value:    "453",
				Operator: "equal",
				Type:     "int",
			},
		}), auth.NewTestingAuthentication(false, ""))
		m, err := handler.Handle()
		gomega.Expect(m).ShouldNot(gomega.BeNil())
		gomega.Expect(err).Should(gomega.BeNil())

		models := m.([]View)
		gomega.Expect(err).Should(gomega.BeNil())

		gomega.Expect(len(models)).Should(gomega.Equal(0))
	})

	ginkgo.It("should return a paginated list of map items based on EQUAL query, the result should be empty", ginkgo.Label("public_api", "list_query"), func() {
		projectId := testCreateProject("project")
		_, publishView := publishFullProject(projectId)

		handler := New(NewModel(publishView.Name, projectId, "paginationList", 1, 100, "desc", "index", "", []string{}, []string{}, Options{}, []queryProcessor.Query{
			{
				Column:   "one",
				Value:    "one",
				Operator: "equal",
				Type:     "string",
			},
			{
				Column:   "four",
				Value:    "34534345",
				Operator: "equal",
				Type:     "int",
			},
		}), auth.NewTestingAuthentication(false, ""))
		m, err := handler.Handle()
		gomega.Expect(m).ShouldNot(gomega.BeNil())
		gomega.Expect(err).Should(gomega.BeNil())

		models := m.([]View)
		gomega.Expect(err).Should(gomega.BeNil())

		gomega.Expect(len(models)).Should(gomega.Equal(0))
	})

	ginkgo.It("should return a paginated list of map items based on UNEQUAL query, the result should be empty", ginkgo.Label("public_api", "list_query"), func() {
		projectId := testCreateProject("project")
		_, publishView := publishFullProject(projectId)

		handler := New(NewModel(publishView.Name, projectId, "paginationList", 1, 100, "desc", "index", "", []string{}, []string{}, Options{}, []queryProcessor.Query{
			{
				Column:   "one",
				Value:    "one",
				Operator: "unequal",
				Type:     "string",
			},
			{
				Column:   "four",
				Value:    "453",
				Operator: "equal",
				Type:     "int",
			},
		}), auth.NewTestingAuthentication(false, ""))
		m, err := handler.Handle()
		gomega.Expect(m).ShouldNot(gomega.BeNil())
		gomega.Expect(err).Should(gomega.BeNil())

		models := m.([]View)
		gomega.Expect(err).Should(gomega.BeNil())

		gomega.Expect(len(models)).Should(gomega.Equal(0))
	})

	ginkgo.It("should return a paginated list of map items based on UNEQUAL query, the result should be empty when using decimal numbers", ginkgo.Label("public_api", "list_query"), func() {
		projectId := testCreateProject("project")
		_, publishView := publishFullProject(projectId)

		handler := New(NewModel(publishView.Name, projectId, "paginationList", 1, 100, "desc", "index", "", []string{}, []string{}, Options{}, []queryProcessor.Query{
			{
				Column:   "five",
				Value:    "456.4345",
				Operator: "equal",
				Type:     "float",
			},
		}), auth.NewTestingAuthentication(false, ""))
		m, err := handler.Handle()
		gomega.Expect(m).ShouldNot(gomega.BeNil())
		gomega.Expect(err).Should(gomega.BeNil())

		models := m.([]View)
		gomega.Expect(err).Should(gomega.BeNil())

		gomega.Expect(len(models)).Should(gomega.Equal(0))
	})

	ginkgo.It("should return a paginated list of map items based on GREATER THAN query, the result should not be empty,", ginkgo.Label("public_api", "list_query"), func() {
		projectId := testCreateProject("project")
		items, publishView := publishFullProject(projectId)

		handler := New(NewModel(publishView.Name, projectId, "paginationList", 1, 100, "desc", "index", "", []string{}, []string{}, Options{}, []queryProcessor.Query{
			{
				Column:   "one",
				Value:    "one",
				Operator: "equal",
				Type:     "string",
			},
			{
				Column:   "four",
				Value:    "400",
				Operator: "greaterThan",
				Type:     "int",
			},
		}), auth.NewTestingAuthentication(false, ""))
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

			sdk.IncludesFn(items, func(item addToList.View) bool {
				return item.ID == model.ID
			})
		}
	})

	ginkgo.It("should return a paginated list of map items based on GREATER THAN query with double precision value, the result should not be empty,", ginkgo.Label("public_api", "list_query"), func() {
		projectId := testCreateProject("project")
		items, publishView := publishFullProject(projectId)

		handler := New(NewModel(publishView.Name, projectId, "paginationList", 1, 100, "desc", "index", "", []string{}, []string{}, Options{}, []queryProcessor.Query{
			{
				Column:   "one",
				Value:    "one",
				Operator: "equal",
				Type:     "string",
			},
			{
				Column:   "five",
				Value:    "43.56",
				Operator: "greaterThan",
				Type:     "float",
			},
		}), auth.NewTestingAuthentication(false, ""))
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

			sdk.IncludesFn(items, func(item addToList.View) bool {
				return item.ID == model.ID
			})
		}
	})

	ginkgo.It("should return a paginated list of map items based on GREATER THAN OR EQUAL with both numeric values, the result should not be empty,", ginkgo.Label("public_api", "list_query"), func() {
		projectId := testCreateProject("project")
		items, publishView := publishFullProject(projectId)

		handler := New(NewModel(publishView.Name, projectId, "paginationList", 1, 100, "desc", "index", "", []string{}, []string{}, Options{}, []queryProcessor.Query{
			{
				Column:   "four",
				Value:    "453",
				Operator: "greaterThanOrEqual",
				Type:     "int",
			},
			{
				Column:   "five",
				Value:    "456.43",
				Operator: "greaterThanOrEqual",
				Type:     "float",
			},
		}), auth.NewTestingAuthentication(false, ""))
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

			sdk.IncludesFn(items, func(item addToList.View) bool {
				return item.ID == model.ID
			})
		}
	})

	ginkgo.It("should return a paginated list of map items based on LESS THAN query with double precision value, the result should not be empty,", ginkgo.Label("public_api", "list_query"), func() {
		projectId := testCreateProject("project")
		items, publishView := publishFullProject(projectId)

		handler := New(NewModel(publishView.Name, projectId, "paginationList", 1, 100, "desc", "index", "", []string{}, []string{}, Options{}, []queryProcessor.Query{
			{
				Column:   "one",
				Value:    "one",
				Operator: "equal",
				Type:     "string",
			},
			{
				Column:   "five",
				Value:    "567.498",
				Operator: "lessThan",
				Type:     "float",
			},
		}), auth.NewTestingAuthentication(false, ""))
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

			sdk.IncludesFn(items, func(item addToList.View) bool {
				return item.ID == model.ID
			})
		}
	})

	ginkgo.It("should return a paginated list of map items based on LESS THAN OR EQUAL with both numeric values, the result should not be empty,", ginkgo.Label("public_api", "list_query"), func() {
		projectId := testCreateProject("project")
		items, publishView := publishFullProject(projectId)

		handler := New(NewModel(publishView.Name, projectId, "paginationList", 1, 100, "desc", "index", "", []string{}, []string{}, Options{}, []queryProcessor.Query{
			{
				Column:   "four",
				Value:    "600",
				Operator: "lessThanOrEqual",
				Type:     "int",
			},
			{
				Column:   "five",
				Value:    "456.43",
				Operator: "lessThanOrEqual",
				Type:     "float",
			},
		}), auth.NewTestingAuthentication(false, ""))
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

			sdk.IncludesFn(items, func(item addToList.View) bool {
				return item.ID == model.ID
			})
		}
	})
})
