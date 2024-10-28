package paginateListItems

import (
	"creatif/pkg/app/auth"
	"creatif/pkg/app/services/shared/queryProcessor"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Public API query validation tests", func() {
	ginkgo.It("should fail validation if Column is empty", ginkgo.Label("public_api", "list_query"), func() {
		projectId := testCreateProject("project")
		publishFullProject(projectId)

		handler := New(NewModel("", projectId, "paginationMap", 1, 100, "desc", "index", "", []string{}, []string{}, Options{}, []queryProcessor.Query{
			{
				Column:   "",
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
		gomega.Expect(m).Should(gomega.BeNil())
		gomega.Expect(err).ShouldNot(gomega.BeNil())
	})

	ginkgo.It("should fail validation if Operator is empty", ginkgo.Label("public_api", "list_query"), func() {
		projectId := testCreateProject("project")
		publishFullProject(projectId)

		handler := New(NewModel("", projectId, "paginationMap", 1, 100, "desc", "index", "", []string{}, []string{}, Options{}, []queryProcessor.Query{
			{
				Column:   "one",
				Value:    "one",
				Operator: "equal",
				Type:     "string",
			},
			{
				Column:   "four",
				Value:    "453",
				Operator: "",
				Type:     "int",
			},
		}), auth.NewTestingAuthentication(false, ""))
		m, err := handler.Handle()
		gomega.Expect(m).Should(gomega.BeNil())
		gomega.Expect(err).ShouldNot(gomega.BeNil())
	})

	ginkgo.It("should fail validation if Operator is invalid", ginkgo.Label("public_api", "list_query"), func() {
		projectId := testCreateProject("project")
		publishFullProject(projectId)

		handler := New(NewModel("", projectId, "paginationMap", 1, 100, "desc", "index", "", []string{}, []string{}, Options{}, []queryProcessor.Query{
			{
				Column:   "one",
				Value:    "one",
				Operator: "equal",
				Type:     "string",
			},
			{
				Column:   "four",
				Value:    "453",
				Operator: "invalid",
				Type:     "int",
			},
		}), auth.NewTestingAuthentication(false, ""))
		m, err := handler.Handle()
		gomega.Expect(m).Should(gomega.BeNil())
		gomega.Expect(err).ShouldNot(gomega.BeNil())
	})

	ginkgo.It("should fail validation if Type is invalid", ginkgo.Label("public_api", "list_query"), func() {
		projectId := testCreateProject("project")
		publishFullProject(projectId)

		handler := New(NewModel("", projectId, "paginationMap", 1, 100, "desc", "index", "", []string{}, []string{}, Options{}, []queryProcessor.Query{
			{
				Column:   "one",
				Value:    "one",
				Operator: "equal",
				Type:     "string",
			},
			{
				Column:   "four",
				Value:    "453",
				Operator: "unequal",
				Type:     "invalid",
			},
		}), auth.NewTestingAuthentication(false, ""))
		m, err := handler.Handle()
		gomega.Expect(m).Should(gomega.BeNil())
		gomega.Expect(err).ShouldNot(gomega.BeNil())
	})
})
