package create

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Declaration node tests", func() {
	ginkgo.It("should create a text declaration node", func() {
		name := uuid.NewString()
		b, _ := json.Marshal(map[string]interface{}{
			"one":  1,
			"two":  "three",
			"four": "six",
		})
		handler := New(NewCreateNodeModel(name, "modifiable", []string{"one", "two", "three"}, b, NodeValidation{}))

		view, err := handler.Handle()
		testAssertErrNil(err)
		testAssertIDValid(view.ID.String())

		gomega.Expect(view.Name).ShouldNot(gomega.BeEmpty())
		gomega.Expect(view.Behaviour).Should(gomega.Equal("modifiable"))
		gomega.Expect(view.Metadata).ShouldNot(gomega.BeEmpty())
		gomega.Expect(view.Groups).ShouldNot(gomega.BeEmpty())
		gomega.Expect(view.CreatedAt).ShouldNot(gomega.BeNil())
		gomega.Expect(view.UpdatedAt).ShouldNot(gomega.BeNil())
	})

	ginkgo.It("should create a boolean declaration node", func() {
		name := uuid.NewString()
		b, _ := json.Marshal(map[string]interface{}{
			"one":  1,
			"two":  "three",
			"four": "six",
		})
		handler := New(NewCreateNodeModel(name, "modifiable", []string{"one", "two", "three"}, b, NodeValidation{}))

		view, err := handler.Handle()
		testAssertErrNil(err)
		testAssertIDValid(view.ID.String())

		gomega.Expect(view.Name).ShouldNot(gomega.BeEmpty())
		gomega.Expect(view.Behaviour).Should(gomega.Equal("modifiable"))
		gomega.Expect(view.Metadata).ShouldNot(gomega.BeEmpty())
		gomega.Expect(view.Groups).ShouldNot(gomega.BeEmpty())
		gomega.Expect(view.CreatedAt).ShouldNot(gomega.BeNil())
		gomega.Expect(view.UpdatedAt).ShouldNot(gomega.BeNil())
	})
})
