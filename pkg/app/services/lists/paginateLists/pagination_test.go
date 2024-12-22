package paginateLists

import (
	"creatif/pkg/app/auth"
	"fmt"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Lists pagination tests", func() {
	ginkgo.It("should paginate through lists", func() {
		projectId := testCreateProject("project")
		for i := 0; i < 100; i++ {
			testCreateListAndReturnNameAndID(projectId, fmt.Sprintf("name-%d", i), 10)
		}

		handler := New(NewModel(projectId, "created_at", "", "desc", 10, 1), auth.NewTestingAuthentication(false, ""))
		views, err := handler.Handle()
		testAssertErrNil(err)

		gomega.Expect(len(views.Data)).Should(gomega.Equal(100))
		gomega.Expect(views.Total).Should(gomega.Equal(int64(100)))
	})

	ginkgo.It("should return lists search by name with regex", func() {
		projectId := testCreateProject("project")
		for i := 0; i < 100; i++ {
			testCreateListAndReturnNameAndID(projectId, fmt.Sprintf("name-%d", i), 10)
		}

		handler := New(NewModel(projectId, "created_at", "1", "desc", 10, 1), auth.NewTestingAuthentication(false, ""))
		views, err := handler.Handle()
		testAssertErrNil(err)

		gomega.Expect(len(views.Data)).Should(gomega.Equal(19))
		gomega.Expect(views.Total).Should(gomega.Equal(int64(19)))
	})
})
