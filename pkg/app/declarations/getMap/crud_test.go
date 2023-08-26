package getMap

import (
	"creatif/pkg/app/assignments/create"
	"creatif/pkg/lib/sdk"
	"fmt"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("GET map tests", func() {
	ginkgo.It("should get names only (default) representation of map of nodes", func() {
		nodes := make([]create.View, 0)
		for i := 0; i < 10; i++ {
			nodes = append(nodes, testCreateBasicAssignmentTextNode(fmt.Sprintf("name-%d", i)))
		}

		view := testCreateMap("mapName", sdk.Map(nodes, func(idx int, value create.View) string {
			return value.ID
		}))

		handler := New(NewGetMapModel(view.ID, "", []string{}))

		mapNodesView, err := handler.Handle()
		testAssertErrNil(err)
		testAssertIDValid(mapNodesView.ID)
		gomega.Expect(mapNodesView.Nodes).Should(gomega.HaveLen(len(nodes)))

		viewNodes := mapNodesView.Nodes.([]FullNode)
		gomega.Expect(viewNodes).Should(gomega.HaveLen(10))
	})

	ginkgo.It("should get full representation of map of nodes", func() {
		nodes := make([]create.View, 0)
		for i := 0; i < 10; i++ {
			nodes = append(nodes, testCreateBasicAssignmentTextNode(fmt.Sprintf("name-%d", i)))
		}

		view := testCreateMap("mapName", sdk.Map(nodes, func(idx int, value create.View) string {
			return value.ID
		}))

		handler := New(NewGetMapModel(view.ID, "full", []string{}))

		mapNodesView, err := handler.Handle()
		testAssertErrNil(err)
		testAssertIDValid(mapNodesView.ID)
		gomega.Expect(mapNodesView.Nodes).Should(gomega.HaveLen(len(nodes)))

		viewNodes := mapNodesView.Nodes.([]FullNode)
		gomega.Expect(viewNodes).Should(gomega.HaveLen(10))
	})

	ginkgo.It("should get representation of map of nodes by custom fields", func() {
		nodes := make([]create.View, 0)
		for i := 0; i < 10; i++ {
			nodes = append(nodes, testCreateBasicAssignmentTextNode(fmt.Sprintf("name-%d", i)))
		}

		view := testCreateMap("mapName", sdk.Map(nodes, func(idx int, value create.View) string {
			return value.ID
		}))

		handler := New(NewGetMapModel(view.ID, "", []string{"groups", "value"}))

		mapNodesView, err := handler.Handle()
		testAssertErrNil(err)
		testAssertIDValid(mapNodesView.ID)
		gomega.Expect(mapNodesView.Nodes).Should(gomega.HaveLen(len(nodes)))

		viewNodes := mapNodesView.Nodes.([]CustomNode)
		gomega.Expect(viewNodes).Should(gomega.HaveLen(10))
	})
})
