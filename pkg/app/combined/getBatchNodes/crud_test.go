package getBatchNodes

import (
	create "creatif/pkg/app/declarations/createNode"
	"creatif/pkg/lib/sdk"
	"fmt"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Batch nodes tests", func() {
	ginkgo.It("should getNode a batch of declaration nodes with full data", func() {
		declarationsNodes := make([]create.View, 0)
		for i := 0; i < 30; i++ {
			declarationsNodes = append(declarationsNodes, testCreateDeclarationNode(fmt.Sprintf("name-%d", i), "modifiable"))
		}

		names := make([]string, 0)
		names = append(names, sdk.Map(declarationsNodes, func(idx int, value create.View) string {
			return value.Name
		})...)

		model := make(map[string]string)
		for _, name := range names {
			model[name] = "node"
		}

		handler := New(NewModel(model))
		views, err := handler.Handle()
		testAssertErrNil(err)

		nodes := views["nodes"]
		viewKeys := sdk.Keys(nodes.(map[string][]View))

		for _, viewName := range viewKeys {
			gomega.Expect(sdk.Includes(names, viewName)).Should(gomega.BeTrue())
		}
	})

	ginkgo.It("should get a batch of maps with full data", func() {
		maps := make([]string, 0)
		for i := 0; i < 100; i++ {
			view := testCreateMap(fmt.Sprintf("name-%d", i), 100)
			maps = append(maps, view.Name)
		}

		model := make(map[string]string)
		for _, name := range maps {
			model[name] = "map"
		}

		handler := New(NewModel(model))
		views, err := handler.Handle()
		testAssertErrNil(err)

		gomega.Expect(views).Should(gomega.HaveKey("nodes"))
		gomega.Expect(views).Should(gomega.HaveKey("maps"))
		gomega.Expect(views["maps"]).Should(gomega.HaveLen(100))
	})

	ginkgo.It("should get a batch of maps and variables with full data", func() {
		declarationsNodes := make([]create.View, 0)
		for i := 0; i < 30; i++ {
			declarationsNodes = append(declarationsNodes, testCreateDeclarationNode(fmt.Sprintf("name-%d", i), "modifiable"))
		}

		names := make([]string, 0)
		names = append(names, sdk.Map(declarationsNodes, func(idx int, value create.View) string {
			return value.Name
		})...)

		maps := make([]string, 0)
		for i := 0; i < 100; i++ {
			view := testCreateMap(fmt.Sprintf("name-%d", i), 100)
			maps = append(maps, view.Name)
		}

		model := make(map[string]string)
		for _, name := range maps {
			model[name] = "map"
		}

		for _, name := range names {
			model[name] = "node"
		}

		handler := New(NewModel(model))
		views, err := handler.Handle()
		testAssertErrNil(err)

		gomega.Expect(views).Should(gomega.HaveKey("nodes"))
		gomega.Expect(views).Should(gomega.HaveKey("maps"))
		gomega.Expect(views["maps"]).Should(gomega.HaveLen(100))
		gomega.Expect(views["nodes"]).Should(gomega.HaveLen(30))
	})
})
