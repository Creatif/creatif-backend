package getMap

import (
	"creatif/pkg/app/auth"
	"creatif/pkg/lib/logger"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("GET map tests", func() {
	ginkgo.It("should getVariable names only (default) representation of map of variables", func() {
		projectId := testCreateProject("project")
		view := testCreateMap(projectId, "mapName", 10)

		handler := New(NewModel(projectId, "eng", view.Name, []string{}, []string{}), auth.NewNoopAuthentication(), logger.NewLogBuilder())

		mapVariablesView, err := handler.Handle()
		testAssertErrNil(err)
		testAssertIDValid(mapVariablesView.ID)
		gomega.Expect(mapVariablesView.Variables).Should(gomega.HaveLen(10))
		gomega.Expect(mapVariablesView.Locale).Should(gomega.Equal("eng"))
	})

	ginkgo.It("should get specific fields from a map variable", func() {
		projectId := testCreateProject("project")
		view := testCreateMap(projectId, "mapName", 10)

		handler := New(NewModel(projectId, "eng", view.Name, []string{"groups", "value"}, []string{}), auth.NewNoopAuthentication(), logger.NewLogBuilder())

		mapVariablesView, err := handler.Handle()
		testAssertErrNil(err)
		testAssertIDValid(mapVariablesView.ID)
		gomega.Expect(mapVariablesView.Variables).Should(gomega.HaveLen(10))
		gomega.Expect(mapVariablesView.Locale).Should(gomega.Equal("eng"))

		for _, n := range mapVariablesView.Variables {
			gomega.Expect(n["id"]).ShouldNot(gomega.BeEmpty())
			gomega.Expect(n["name"]).ShouldNot(gomega.BeEmpty())
			gomega.Expect(n["value"]).ShouldNot(gomega.BeEmpty())
			gomega.Expect(n["groups"]).ShouldNot(gomega.BeEmpty())
		}
	})

	ginkgo.It("should get map variables of a specific group", func() {
		projectId := testCreateProject("project")
		view := testCreateMap(projectId, "mapName", 100)

		handler := New(NewModel(projectId, "eng", view.Name, []string{"groups", "value"}, []string{"one"}), auth.NewNoopAuthentication(), logger.NewLogBuilder())

		mapVariablesView, err := handler.Handle()
		testAssertErrNil(err)
		testAssertIDValid(mapVariablesView.ID)
		gomega.Expect(len(mapVariablesView.Variables)).Should(gomega.Equal(50))
		gomega.Expect(mapVariablesView.Locale).Should(gomega.Equal("eng"))

		for _, n := range mapVariablesView.Variables {
			gomega.Expect(n["id"]).ShouldNot(gomega.BeEmpty())
			gomega.Expect(n["name"]).ShouldNot(gomega.BeEmpty())
			gomega.Expect(n["value"]).ShouldNot(gomega.BeEmpty())
			gomega.Expect(n["groups"]).ShouldNot(gomega.BeEmpty())
		}
	})
})
