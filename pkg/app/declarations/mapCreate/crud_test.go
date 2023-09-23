package mapCreate

import (
	"encoding/json"
	"fmt"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Map variable tests", func() {
	ginkgo.It("should create a map with only variable entries", func() {
		projectId := testCreateProject("project")
		entries := make([]Entry, 0)

		m := map[string]interface{}{
			"one":   "one",
			"two":   []string{"one", "two", "three"},
			"three": []int{1, 2, 3},
			"four":  453,
		}

		b, err := json.Marshal(m)
		gomega.Expect(err).Should(gomega.BeNil())

		for i := 0; i < 100; i++ {
			var value interface{}
			value = "my value"
			if i%2 == 0 {
				value = true
			}

			if i%3 == 0 {
				value = map[string]interface{}{
					"one":   "one",
					"two":   []string{"one", "two", "three"},
					"three": []int{1, 2, 3},
					"four":  453,
				}
			}

			v, err := json.Marshal(value)
			gomega.Expect(err).Should(gomega.BeNil())

			variableModel := VariableModel{
				Name:     fmt.Sprintf("name-%d", i),
				Metadata: b,
				Groups: []string{
					"one",
					"two",
					"three",
				},
				Behaviour: "modifiable",
				Value:     v,
			}

			entries = append(entries, Entry{
				Type:  "variable",
				Model: variableModel,
			})
		}

		handler := New(NewModel(projectId, "mapName", entries))
		view, err := handler.Handle()
		testAssertErrNil(err)
		testAssertIDValid(view.ID)

		gomega.Expect(view.Name).Should(gomega.Equal("mapName"))
		gomega.Expect(view.Variables).Should(gomega.HaveLen(100))
	})
})
