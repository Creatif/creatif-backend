package paginateMaps

import (
	"creatif/pkg/app/auth"
	"fmt"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Maps pagination tests", func() {
	ginkgo.It("should paginate through maps", ginkgo.Label("map"), func() {
		projectId := testCreateProject("project")
		for i := 0; i < 100; i++ {
			testCreateMap(projectId, fmt.Sprintf("name-%d", i), 10)
		}

		handler := New(NewModel(projectId, "created_at", "", "desc", 10, 1), auth.NewTestingAuthentication(false, ""))
		views, err := handler.Handle()
		testAssertErrNil(err)

		gomega.Expect(len(views.Data)).Should(gomega.Equal(10))
		gomega.Expect(views.Total).Should(gomega.Equal(int64(100)))
	})

	ginkgo.It("should get an empty result from the end of the maps listing", ginkgo.Label("map"), func() {
		projectId := testCreateProject("project")
		for i := 0; i < 100; i++ {
			testCreateMap(projectId, fmt.Sprintf("name-%d", i), 10)
		}

		handler := New(NewModel(projectId, "created_at", "", "desc", 10, 50), auth.NewTestingAuthentication(false, ""))
		views, err := handler.Handle()
		testAssertErrNil(err)

		gomega.Expect(len(views.Data)).Should(gomega.Equal(0))
		gomega.Expect(views.Total).Should(gomega.Equal(int64(100)))
	})

	ginkgo.It("should return maps search by name with regex", ginkgo.Label("map"), func() {
		projectId := testCreateProject("project")
		for i := 0; i < 100; i++ {
			testCreateMap(projectId, fmt.Sprintf("name-%d", i), 10)
		}

		handler := New(NewModel(projectId, "created_at", "1", "desc", 10, 1), auth.NewTestingAuthentication(false, ""))
		views, err := handler.Handle()
		testAssertErrNil(err)

		gomega.Expect(len(views.Data)).Should(gomega.Equal(10))
		gomega.Expect(views.Total).Should(gomega.Equal(int64(19)))
	})
})
