package paginateVariables

import (
	"fmt"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Variable pagination tests", func() {
	ginkgo.It("should paginate through variables", func() {
		projectId := testCreateProject("project")
		for i := 0; i < 100; i++ {
			testCreateBasicDeclarationTextVariable(projectId, fmt.Sprintf("one-%d", i), "modifiable")
		}

		handler := New(NewModel(projectId, "eng", "created_at", "desc", 10, 1, []string{"one"}, nil))
		views, err := handler.Handle()
		testAssertErrNil(err)

		gomega.Expect(len(views.Data)).Should(gomega.Equal(10))
		gomega.Expect(views.Total).Should(gomega.Equal(int64(100)))
	})

	ginkgo.It("should get an empty result from the end of the variables listing", func() {
		projectId := testCreateProject("project")
		for i := 0; i < 100; i++ {
			testCreateBasicDeclarationTextVariable(projectId, fmt.Sprintf("one-%d", i), "modifiable")
		}

		handler := New(NewModel(projectId, "eng", "created_at", "desc", 10, 50, []string{"one"}, nil))
		views, err := handler.Handle()
		testAssertErrNil(err)

		gomega.Expect(len(views.Data)).Should(gomega.Equal(0))
		gomega.Expect(views.Total).Should(gomega.Equal(int64(100)))
	})

	ginkgo.It("should return empty result for group that does not exist", func() {
		projectId := testCreateProject("project")
		for i := 0; i < 100; i++ {
			testCreateBasicDeclarationTextVariable(projectId, fmt.Sprintf("one-%d", i), "modifiable")
		}

		handler := New(NewModel(projectId, "eng", "created_at", "desc", 10, 1, []string{"not_exists"}, nil))
		views, err := handler.Handle()
		testAssertErrNil(err)

		gomega.Expect(len(views.Data)).Should(gomega.Equal(0))
		gomega.Expect(views.Total).Should(gomega.Equal(int64(100)))
	})
})
