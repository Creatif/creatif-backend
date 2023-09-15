package getMap

import (
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("GET map tests", func() {
	ginkgo.It("should getNode names only (default) representation of map of nodes", func() {
		view := testCreateMap("mapName", 100)

		handler := New(NewModel(view.Name, []string{}))

		mapNodesView, err := handler.Handle()
		testAssertErrNil(err)
		testAssertIDValid(mapNodesView.ID)
		gomega.Expect(mapNodesView.Nodes).Should(gomega.HaveLen(100))
	})

	ginkgo.It("should get specific fields from a map node", func() {
		view := testCreateMap("mapName", 100)

		handler := New(NewModel(view.Name, []string{"groups", "value"}))

		mapNodesView, err := handler.Handle()
		testAssertErrNil(err)
		testAssertIDValid(mapNodesView.ID)
		gomega.Expect(mapNodesView.Nodes).Should(gomega.HaveLen(100))

		for _, n := range mapNodesView.Nodes {
			gomega.Expect(n["id"]).ShouldNot(gomega.BeEmpty())
			gomega.Expect(n["name"]).ShouldNot(gomega.BeEmpty())
			gomega.Expect(n["value"]).ShouldNot(gomega.BeEmpty())
			gomega.Expect(n["groups"]).ShouldNot(gomega.BeEmpty())
		}
	})
})
