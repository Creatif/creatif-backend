package mapCreate

import (
	"creatif/pkg/app/auth"
	"encoding/json"
	"fmt"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Map variable tests", func() {
	ginkgo.It("should create multiple maps with different name with only variable entries", ginkgo.Label("map"), func() {
		projectId := testCreateProject("project")
		entries := make([]VariableModel, 0)

		m := map[string]interface{}{
			"one":   "one",
			"two":   []string{"one", "two", "three"},
			"three": []int{1, 2, 3},
			"four":  453,
		}

		b, err := json.Marshal(m)
		gomega.Expect(err).Should(gomega.BeNil())

		for i := 0; i < 10; i++ {
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
				Locale:    "eng",
				Behaviour: "modifiable",
				Value:     v,
			}

			entries = append(entries, variableModel)
		}

		handler := New(NewModel(projectId, "mapName", entries), auth.NewTestingAuthentication(false, ""))
		view, err := handler.Handle()
		testAssertErrNil(err)
		testAssertIDValid(view.ID)

		gomega.Expect(view.Name).Should(gomega.Equal("mapName"))
		gomega.Expect(view.Variables).Should(gomega.HaveLen(10))

		handler = New(NewModel(projectId, "otherMapName", entries), auth.NewTestingAuthentication(false, ""))
		view, err = handler.Handle()
		testAssertErrNil(err)
		testAssertIDValid(view.ID)

		gomega.Expect(view.Name).Should(gomega.Equal("otherMapName"))
		gomega.Expect(view.Variables).Should(gomega.HaveLen(10))
	})

	ginkgo.It("should fail on the database level when trying to create a map with name that already exists", func() {
		projectId := testCreateProject("project")
		entries := make([]VariableModel, 0)

		m := map[string]interface{}{
			"one":   "one",
			"two":   []string{"one", "two", "three"},
			"three": []int{1, 2, 3},
			"four":  453,
		}

		b, err := json.Marshal(m)
		gomega.Expect(err).Should(gomega.BeNil())

		for i := 0; i < 10; i++ {
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
				Locale:    "eng",
				Behaviour: "modifiable",
				Value:     v,
			}

			entries = append(entries, variableModel)
		}

		handler := New(NewModel(projectId, "mapName", entries), auth.NewTestingAuthentication(false, ""))
		view, err := handler.Handle()
		testAssertErrNil(err)
		testAssertIDValid(view.ID)

		gomega.Expect(view.Name).Should(gomega.Equal("mapName"))
		gomega.Expect(view.Variables).Should(gomega.HaveLen(10))

		handler = New(NewModel(projectId, "mapName", entries), auth.NewTestingAuthentication(false, ""))
		_, err = handler.Logic()
		gomega.Expect(err).ShouldNot(gomega.BeNil())
	})

	ginkgo.It("should fail on the application level when trying to create a map with name that already exists", func() {
		projectId := testCreateProject("project")
		entries := make([]VariableModel, 0)

		m := map[string]interface{}{
			"one":   "one",
			"two":   []string{"one", "two", "three"},
			"three": []int{1, 2, 3},
			"four":  453,
		}

		b, err := json.Marshal(m)
		gomega.Expect(err).Should(gomega.BeNil())

		for i := 0; i < 10; i++ {
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
				Locale:    "eng",
				Behaviour: "modifiable",
				Value:     v,
			}

			entries = append(entries, variableModel)
		}

		handler := New(NewModel(projectId, "mapName", entries), auth.NewTestingAuthentication(false, ""))
		view, err := handler.Handle()
		testAssertErrNil(err)
		testAssertIDValid(view.ID)

		gomega.Expect(view.Name).Should(gomega.Equal("mapName"))
		gomega.Expect(view.Variables).Should(gomega.HaveLen(10))

		handler = New(NewModel(projectId, "mapName", entries), auth.NewTestingAuthentication(false, ""))
		_, err = handler.Handle()
		gomega.Expect(err).ShouldNot(gomega.BeNil())
	})
})
