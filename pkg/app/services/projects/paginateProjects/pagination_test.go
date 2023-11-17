package paginateProjects

import (
	"creatif/pkg/app/auth"
	"creatif/pkg/lib/logger"
	"fmt"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Projects pagination tests", func() {
	ginkgo.It("should paginate through variables", func() {
		for i := 0; i < 10; i++ {
			testCreateProject(fmt.Sprintf("name-%d", i))
		}

		handler := New(NewModel("created_at", "", "desc", 2, 1), auth.NewTestingAuthentication(false), logger.NewLogBuilder())
		views, err := handler.Handle()
		testAssertErrNil(err)

		gomega.Expect(len(views.Data)).Should(gomega.Equal(2))
		gomega.Expect(views.Total).Should(gomega.Equal(int64(10)))
	})

	ginkgo.It("should get an empty result from the end of the variables listing", func() {
		for i := 0; i < 10; i++ {
			testCreateProject(fmt.Sprintf("name-%d", i))
		}

		handler := New(NewModel("created_at", "", "desc", 2, 50), auth.NewTestingAuthentication(false), logger.NewLogBuilder())

		views, err := handler.Handle()
		testAssertErrNil(err)

		gomega.Expect(len(views.Data)).Should(gomega.Equal(0))
		gomega.Expect(views.Total).Should(gomega.Equal(int64(10)))
	})

	ginkgo.It("should return variables search by name without groups", func() {
		ginkgo.Skip("")
		for i := 0; i < 100; i++ {
			testCreateProject(fmt.Sprintf("name-%d", i))
		}

		handler := New(NewModel("created_at", "", "desc", 10, 1), auth.NewTestingAuthentication(false), logger.NewLogBuilder())

		views, err := handler.Handle()
		testAssertErrNil(err)

		gomega.Expect(len(views.Data)).Should(gomega.Equal(10))
		gomega.Expect(views.Total).Should(gomega.Equal(int64(19)))
	})
})
