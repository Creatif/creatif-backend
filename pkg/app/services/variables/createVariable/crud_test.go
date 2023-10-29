package createVariable

import (
	"creatif/pkg/app/auth"
	"creatif/pkg/lib/logger"
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
		handler := New(NewModel(projectId, "eng", name, "modifiable", []string{"one", "two", "three"}, b, b), auth.NewNoopAuthentication(), logger.NewLogBuilder())

		view, err := handler.Handle()
		testAssertErrNil(err)
		testAssertIDValid(view.ID)

		gomega.Expect(view.Name).Should(gomega.Equal(name))
		gomega.Expect(view.Behaviour).Should(gomega.Equal("modifiable"))
		gomega.Expect(view.Metadata).ShouldNot(gomega.BeEmpty())
		gomega.Expect(view.Groups).ShouldNot(gomega.BeEmpty())
		gomega.Expect(view.Locale).ShouldNot(gomega.BeEmpty())
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
		handler := New(NewModel(projectId, "eng", name, "modifiable", []string{"one", "two", "three"}, b, b), auth.NewNoopAuthentication(), logger.NewLogBuilder())

		view, err := handler.Handle()
		testAssertErrNil(err)
		testAssertIDValid(view.ID)

		gomega.Expect(view.Name).Should(gomega.Equal(name))
		gomega.Expect(view.Behaviour).Should(gomega.Equal("modifiable"))
		gomega.Expect(view.Metadata).ShouldNot(gomega.BeEmpty())
		gomega.Expect(view.Groups).ShouldNot(gomega.BeEmpty())
		gomega.Expect(view.Locale).ShouldNot(gomega.BeEmpty())
		gomega.Expect(view.CreatedAt).ShouldNot(gomega.BeNil())
		gomega.Expect(view.UpdatedAt).ShouldNot(gomega.BeNil())
		gomega.Expect(view.ProjectID).ShouldNot(gomega.BeEmpty())
	})

	ginkgo.It("should fail on database to create a variable with the same name on a same project", func() {
		projectId := testCreateProject("project")
		name, _ := sdk.NewULID()
		b, _ := json.Marshal(map[string]interface{}{
			"one":  1,
			"two":  "three",
			"four": "six",
		})
		handler := New(NewModel(projectId, "eng", name, "modifiable", []string{"one", "two", "three"}, b, b), auth.NewNoopAuthentication(), logger.NewLogBuilder())

		view, err := handler.Handle()
		testAssertErrNil(err)
		testAssertIDValid(view.ID)

		gomega.Expect(view.Name).Should(gomega.Equal(name))
		gomega.Expect(view.Behaviour).Should(gomega.Equal("modifiable"))
		gomega.Expect(view.Metadata).ShouldNot(gomega.BeEmpty())
		gomega.Expect(view.Groups).ShouldNot(gomega.BeEmpty())
		gomega.Expect(view.CreatedAt).ShouldNot(gomega.BeNil())
		gomega.Expect(view.Locale).ShouldNot(gomega.BeEmpty())
		gomega.Expect(view.UpdatedAt).ShouldNot(gomega.BeNil())
		gomega.Expect(view.ProjectID).ShouldNot(gomega.BeEmpty())

		handler = New(NewModel(projectId, "eng", name, "modifiable", []string{"one", "two", "three"}, b, b), auth.NewNoopAuthentication(), logger.NewLogBuilder())

		// skipping validation
		_, err = handler.Logic()
		gomega.Expect(err).ShouldNot(gomega.BeNil())
	})

	ginkgo.It("should create variables with equal name on different projects when skipping validation", func() {
		projectId := testCreateProject("project")
		name, _ := sdk.NewULID()
		b, _ := json.Marshal(map[string]interface{}{
			"one":  1,
			"two":  "three",
			"four": "six",
		})
		handler := New(NewModel(projectId, "eng", name, "modifiable", []string{"one", "two", "three"}, b, b), auth.NewNoopAuthentication(), logger.NewLogBuilder())

		view, err := handler.Logic()
		testAssertErrNil(err)
		testAssertIDValid(view.ID)

		gomega.Expect(view.Name).Should(gomega.Equal(name))
		gomega.Expect(view.Behaviour).Should(gomega.Equal("modifiable"))
		gomega.Expect(view.Metadata).ShouldNot(gomega.BeEmpty())
		gomega.Expect(view.Groups).ShouldNot(gomega.BeEmpty())
		gomega.Expect(view.CreatedAt).ShouldNot(gomega.BeNil())
		gomega.Expect(view.UpdatedAt).ShouldNot(gomega.BeNil())
		gomega.Expect(view.ProjectID).ShouldNot(gomega.BeEmpty())

		projectId = testCreateProject("different project")
		handler = New(NewModel(projectId, "eng", name, "modifiable", []string{"one", "two", "three"}, b, b), auth.NewNoopAuthentication(), logger.NewLogBuilder())

		logicView, err := handler.Logic()
		testAssertErrNil(err)
		testAssertIDValid(logicView.ID)

		gomega.Expect(view.Name).Should(gomega.Equal(name))
		gomega.Expect(view.Behaviour).Should(gomega.Equal("modifiable"))
		gomega.Expect(view.Metadata).ShouldNot(gomega.BeEmpty())
		gomega.Expect(view.Groups).ShouldNot(gomega.BeEmpty())
		gomega.Expect(view.CreatedAt).ShouldNot(gomega.BeNil())
		gomega.Expect(view.UpdatedAt).ShouldNot(gomega.BeNil())
		gomega.Expect(view.ProjectID).ShouldNot(gomega.BeEmpty())
	})
})
