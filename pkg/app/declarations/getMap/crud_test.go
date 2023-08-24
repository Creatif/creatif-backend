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

		viewNodes := mapNodesView.Nodes.([]map[string]interface{})
		for _, value := range viewNodes {
			gomega.Expect(value).Should(gomega.HaveLen(2))
			gomega.Expect(value).Should(gomega.HaveKey("name"))
			gomega.Expect(value).Should(gomega.HaveKey("id"))

			gomega.Expect(value["id"]).ShouldNot(gomega.BeEmpty())
			gomega.Expect(value["name"]).ShouldNot(gomega.BeEmpty())
		}
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

		viewNodes := mapNodesView.Nodes.([]map[string]interface{})
		for _, value := range viewNodes {
			gomega.Expect(value).Should(gomega.HaveLen(9))
			gomega.Expect(value).Should(gomega.HaveKey("name"))
			gomega.Expect(value).Should(gomega.HaveKey("id"))
			gomega.Expect(value).Should(gomega.HaveKey("behaviour"))
			gomega.Expect(value).Should(gomega.HaveKey("metadata"))
			gomega.Expect(value).Should(gomega.HaveKey("type"))
			gomega.Expect(value).Should(gomega.HaveKey("groups"))
			gomega.Expect(value).Should(gomega.HaveKey("value"))
			gomega.Expect(value).Should(gomega.HaveKey("created_at"))
			gomega.Expect(value).Should(gomega.HaveKey("updated_at"))

			gomega.Expect(value["id"]).ShouldNot(gomega.BeEmpty())
			gomega.Expect(value["name"]).ShouldNot(gomega.BeEmpty())
			gomega.Expect(value["behaviour"]).Should(gomega.Equal("modifiable"))
			gomega.Expect(value["metadata"]).Should(gomega.BeNil())
			gomega.Expect(value["type"]).Should(gomega.Equal("text"))
			gomega.Expect(value["groups"]).Should(gomega.HaveLen(0))
			gomega.Expect(value["value"]).Should(gomega.Equal("this is a text node"))
			gomega.Expect(value["updated_at"]).ShouldNot(gomega.BeNil())
			gomega.Expect(value["created_at"]).ShouldNot(gomega.BeNil())
		}
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

		viewNodes := mapNodesView.Nodes.([]map[string]interface{})
		for _, value := range viewNodes {
			gomega.Expect(value).Should(gomega.HaveLen(4))
			gomega.Expect(value).Should(gomega.HaveKey("name"))
			gomega.Expect(value).Should(gomega.HaveKey("id"))
			gomega.Expect(value).Should(gomega.HaveKey("groups"))
			gomega.Expect(value).Should(gomega.HaveKey("value"))

			gomega.Expect(value["id"]).ShouldNot(gomega.BeEmpty())
			gomega.Expect(value["name"]).ShouldNot(gomega.BeEmpty())
			gomega.Expect(value["groups"]).Should(gomega.HaveLen(0))
			gomega.Expect(value["value"]).Should(gomega.Equal("this is a text node"))
		}
	})
})
