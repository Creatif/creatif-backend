package paginateLists

import (
	"creatif/pkg/app/auth"
	"creatif/pkg/lib/logger"
	"fmt"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("List items pagination tests", func() {
	ginkgo.It("should paginate through list variables", func() {
		projectId := testCreateProject("project")
		for i := 0; i < 100; i++ {
			testCreateListAndReturnNameAndID(projectId, fmt.Sprintf("name-%d", i), 10)
		}

		handler := New(NewModel(projectId, "created_at", "", "desc", 10, 1), auth.NewTestingAuthentication(false), logger.NewLogBuilder())
		views, err := handler.Handle()
		testAssertErrNil(err)

		gomega.Expect(len(views.Data)).Should(gomega.Equal(10))
		gomega.Expect(views.Total).Should(gomega.Equal(int64(100)))
	})

	ginkgo.It("should get an empty result from the end of the list variables listing", func() {
		projectId := testCreateProject("project")
		for i := 0; i < 100; i++ {
			testCreateListAndReturnNameAndID(projectId, fmt.Sprintf("name-%d", i), 10)
		}

		handler := New(NewModel(projectId, "created_at", "", "desc", 10, 50), auth.NewTestingAuthentication(false), logger.NewLogBuilder())
		views, err := handler.Handle()
		testAssertErrNil(err)

		gomega.Expect(len(views.Data)).Should(gomega.Equal(0))
		gomega.Expect(views.Total).Should(gomega.Equal(int64(100)))
	})

	ginkgo.It("should return items search by name with regex", func() {
		projectId := testCreateProject("project")
		for i := 0; i < 100; i++ {
			testCreateListAndReturnNameAndID(projectId, fmt.Sprintf("name-%d", i), 10)
		}

		handler := New(NewModel(projectId, "created_at", "1", "desc", 10, 1), auth.NewTestingAuthentication(false), logger.NewLogBuilder())
		views, err := handler.Handle()
		testAssertErrNil(err)

		gomega.Expect(len(views.Data)).Should(gomega.Equal(10))
		gomega.Expect(views.Total).Should(gomega.Equal(int64(19)))
	})
})
