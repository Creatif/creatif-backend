package getMap

import (
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("GET map tests", func() {
	ginkgo.It("should getVariable names only (default) representation of map of variables", func() {
		projectId := testCreateProject("project")
		view := testCreateMap(projectId, "mapName", 10)

		handler := New(NewModel(projectId, "eng", view.Name, []string{}))

		mapVariablesView, err := handler.Handle()
		testAssertErrNil(err)
		testAssertIDValid(mapVariablesView.ID)
		gomega.Expect(mapVariablesView.Variables).Should(gomega.HaveLen(10))
		gomega.Expect(mapVariablesView.Locale).Should(gomega.Equal("eng"))
	})

	ginkgo.It("should get specific fields from a map variable", func() {
		projectId := testCreateProject("project")
		view := testCreateMap(projectId, "mapName", 10)

		handler := New(NewModel(projectId, "eng", view.Name, []string{"groups", "value"}))

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
})
