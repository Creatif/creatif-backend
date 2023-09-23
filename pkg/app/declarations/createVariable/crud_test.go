package createVariable

import (
	"creatif/pkg/lib/sdk"
	"encoding/json"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Declaration variable tests", func() {
	ginkgo.It("should create a text declaration variable", func() {
		projectId := testCreateProject("project")
		name, _ := sdk.NewULID()
		b, _ := json.Marshal(map[string]interface{}{
			"one":  1,
			"two":  "three",
			"four": "six",
		})
		handler := New(NewModel(projectId, name, "modifiable", []string{"one", "two", "three"}, b, b))

		view, err := handler.Handle()
		testAssertErrNil(err)
		testAssertIDValid(view.ID)

		gomega.Expect(view.Name).ShouldNot(gomega.BeEmpty())
		gomega.Expect(view.Behaviour).Should(gomega.Equal("modifiable"))
		gomega.Expect(view.Metadata).ShouldNot(gomega.BeEmpty())
		gomega.Expect(view.Groups).ShouldNot(gomega.BeEmpty())
		gomega.Expect(view.CreatedAt).ShouldNot(gomega.BeNil())
		gomega.Expect(view.UpdatedAt).ShouldNot(gomega.BeNil())
		gomega.Expect(view.ProjectID).ShouldNot(gomega.BeEmpty())
	})

	ginkgo.It("should create a boolean declaration variable", func() {
		projectId := testCreateProject("project")
		name, _ := sdk.NewULID()
		b, _ := json.Marshal(map[string]interface{}{
			"one":  1,
			"two":  "three",
			"four": "six",
		})
		handler := New(NewModel(projectId, name, "modifiable", []string{"one", "two", "three"}, b, b))

		view, err := handler.Handle()
		testAssertErrNil(err)
		testAssertIDValid(view.ID)

		gomega.Expect(view.Name).ShouldNot(gomega.BeEmpty())
		gomega.Expect(view.Behaviour).Should(gomega.Equal("modifiable"))
		gomega.Expect(view.Metadata).ShouldNot(gomega.BeEmpty())
		gomega.Expect(view.Groups).ShouldNot(gomega.BeEmpty())
		gomega.Expect(view.CreatedAt).ShouldNot(gomega.BeNil())
		gomega.Expect(view.UpdatedAt).ShouldNot(gomega.BeNil())
		gomega.Expect(view.ProjectID).ShouldNot(gomega.BeEmpty())
	})
})
