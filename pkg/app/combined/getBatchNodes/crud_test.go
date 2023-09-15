package getBatchNodes

import (
	"creatif/pkg/app/assignments/create"
	"creatif/pkg/lib/sdk"
	"fmt"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Batch nodes tests", func() {
	ginkgo.It("should getNode a batch of declaration nodes with full data", func() {
		textNodes := make([]create.View, 0)
		booleanNodes := make([]create.View, 0)
		jsonNodes := make([]create.View, 0)

		for i := 0; i < 10; i++ {
			textNodes = append(textNodes, testCreateBasicAssignmentTextNode(fmt.Sprintf("name-%d", i), "this is a text node"))
		}

		for i := 10; i < 20; i++ {
			booleanNodes = append(booleanNodes, testCreateBasicAssignmentBooleanNode(fmt.Sprintf("name-%d", i), false))
		}

		for i := 20; i < 30; i++ {
			jsonNodes = append(jsonNodes, testCreateBasicAssignmentTextNode(fmt.Sprintf("name-%d", i), map[string]interface{}{
				"one":   "one",
				"two":   []string{"one", "two"},
				"three": []int{1, 2, 3, 4},
				"four":  583,
			}))
		}

		names := make([]string, 0)
		names = append(names, sdk.Map(textNodes, func(idx int, value create.View) string {
			return value.Name
		})...)

		names = append(names, sdk.Map(booleanNodes, func(idx int, value create.View) string {
			return value.Name
		})...)

		names = append(names, sdk.Map(jsonNodes, func(idx int, value create.View) string {
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
		textNodes := make([]create.View, 0)
		booleanNodes := make([]create.View, 0)
		jsonNodes := make([]create.View, 0)

		for i := 0; i < 10; i++ {
			textNodes = append(textNodes, testCreateBasicAssignmentTextNode(fmt.Sprintf("name-%d", i), "this is a text node"))
		}

		for i := 10; i < 20; i++ {
			booleanNodes = append(booleanNodes, testCreateBasicAssignmentBooleanNode(fmt.Sprintf("name-%d", i), false))
		}

		for i := 20; i < 30; i++ {
			jsonNodes = append(jsonNodes, testCreateBasicAssignmentTextNode(fmt.Sprintf("name-%d", i), map[string]interface{}{
				"one":   "one",
				"two":   []string{"one", "two"},
				"three": []int{1, 2, 3, 4},
				"four":  583,
			}))
		}

		names := make([]string, 0)
		names = append(names, sdk.Map(textNodes, func(idx int, value create.View) string {
			return value.Name
		})...)

		names = append(names, sdk.Map(booleanNodes, func(idx int, value create.View) string {
			return value.Name
		})...)

		names = append(names, sdk.Map(jsonNodes, func(idx int, value create.View) string {
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
