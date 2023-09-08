package pagination

import (
	"fmt"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Declaration node pagination tests", func() {
	ginkgo.It("should return the first row of results by created_at field desc", func() {
		limit := 10
		for i := 0; i < 20; i++ {
			testCreateBasicAssignmentTextNode(fmt.Sprintf("name-%d", i))
		}

		handler := New(NewModel(false, "created_at", "desc", limit))
		views, err := handler.Handle()
		testAssertErrNil(err)

		gomega.Expect(len(views.Items)).Should(gomega.Equal(limit))
		gomega.Expect(views.Items[0].Name).Should(gomega.Equal("name-19"))
		gomega.Expect(views.Items[len(views.Items)-1].Name).Should(gomega.Equal("name-10"))

		gomega.Expect(views.PaginationInfo.Next).ShouldNot(gomega.BeEmpty())
		gomega.Expect(views.PaginationInfo.Prev).ShouldNot(gomega.BeEmpty())
	})

	ginkgo.It("should return the first row of results by created_at field asc", func() {
		limit := 10
		for i := 0; i < 20; i++ {
			testCreateBasicAssignmentTextNode(fmt.Sprintf("name-%d", i))
		}

		handler := New(NewModel(false, "created_at", "asc", limit))
		views, err := handler.Handle()
		testAssertErrNil(err)

		gomega.Expect(len(views.Items)).Should(gomega.Equal(limit))
		gomega.Expect(views.Items[0].Name).Should(gomega.Equal("name-0"))
		gomega.Expect(views.Items[len(views.Items)-1].Name).Should(gomega.Equal("name-9"))

		gomega.Expect(views.PaginationInfo.Next).ShouldNot(gomega.BeEmpty())
		gomega.Expect(views.PaginationInfo.Prev).ShouldNot(gomega.BeEmpty())
	})
})
