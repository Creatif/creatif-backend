package paginateVariables

import (
	"creatif/pkg/lib/sdk/pagination"
	"fmt"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Declaration variable paginateVariables tests", func() {
	ginkgo.It("should return the first row of results by created_at field desc going forward", func() {
		limit := 10
		for i := 0; i < 20; i++ {
			testCreateBasicDeclarationTextVariable(fmt.Sprintf("name-%d", i), "modifiable")
		}

		handler := New(NewModel("", "", "created_at", pagination.DESC, pagination.DIRECTION_FORWARD, limit, []string{}))
		views, err := handler.Handle()
		testAssertErrNil(err)

		gomega.Expect(len(views.Items)).Should(gomega.Equal(limit))
		gomega.Expect(views.Items[0].Name).Should(gomega.Equal("name-19"))
		gomega.Expect(views.Items[len(views.Items)-1].Name).Should(gomega.Equal("name-10"))

		gomega.Expect(views.PaginationInfo.Next).ShouldNot(gomega.BeEmpty())
		gomega.Expect(views.PaginationInfo.Prev).ShouldNot(gomega.BeEmpty())
	})

	ginkgo.It("should return the first row of results by created_at field desc going forward and with specifying all the groups and some that do not exist", func() {
		limit := 10
		for i := 0; i < 20; i++ {
			testCreateBasicDeclarationTextVariable(fmt.Sprintf("name-%d", i), "modifiable")
		}

		handler := New(NewModel("", "", "created_at", pagination.DESC, pagination.DIRECTION_FORWARD, limit, []string{"one", "two", "three", "six"}))
		views, err := handler.Handle()
		testAssertErrNil(err)

		gomega.Expect(len(views.Items)).Should(gomega.Equal(limit))
		gomega.Expect(views.Items[0].Name).Should(gomega.Equal("name-19"))
		gomega.Expect(views.Items[len(views.Items)-1].Name).Should(gomega.Equal("name-10"))

		gomega.Expect(views.PaginationInfo.Next).ShouldNot(gomega.BeEmpty())
		gomega.Expect(views.PaginationInfo.Prev).ShouldNot(gomega.BeEmpty())
	})

	ginkgo.It("should return the first row of results by created_at field asc going forward", func() {
		limit := 10
		for i := 0; i < 20; i++ {
			testCreateBasicDeclarationTextVariable(fmt.Sprintf("name-%d", i), "modifiable")
		}

		handler := New(NewModel("", "", "created_at", pagination.ASC, pagination.DIRECTION_FORWARD, limit, []string{"one"}))
		views, err := handler.Handle()
		testAssertErrNil(err)

		gomega.Expect(len(views.Items)).Should(gomega.Equal(limit))
		gomega.Expect(views.Items[0].Name).Should(gomega.Equal("name-0"))
		gomega.Expect(views.Items[len(views.Items)-1].Name).Should(gomega.Equal("name-9"))

		gomega.Expect(views.PaginationInfo.Next).ShouldNot(gomega.BeEmpty())
		gomega.Expect(views.PaginationInfo.Prev).ShouldNot(gomega.BeEmpty())
	})

	ginkgo.It("nextUrl in paginateVariables info should be an empty string if number of items is less than limit", func() {
		limit := 10
		for i := 0; i < 5; i++ {
			testCreateBasicDeclarationTextVariable(fmt.Sprintf("name-%d", i), "modifiable")
		}

		handler := New(NewModel("", "", "created_at", pagination.ASC, pagination.DIRECTION_FORWARD, limit, []string{"one"}))
		views, err := handler.Handle()
		testAssertErrNil(err)

		gomega.Expect(len(views.Items)).Should(gomega.Equal(5))
		gomega.Expect(views.Items[0].Name).Should(gomega.Equal("name-0"))
		gomega.Expect(views.Items[len(views.Items)-1].Name).Should(gomega.Equal("name-4"))

		gomega.Expect(views.PaginationInfo.Next).Should(gomega.BeEmpty())
		gomega.Expect(views.PaginationInfo.Prev).ShouldNot(gomega.BeEmpty())
	})

	ginkgo.It("should give an empty result for a non existent group", func() {
		limit := 10
		for i := 0; i < 5; i++ {
			testCreateBasicDeclarationTextVariable(fmt.Sprintf("name-%d", i), "modifiable")
		}

		handler := New(NewModel("", "", "created_at", pagination.ASC, pagination.DIRECTION_FORWARD, limit, []string{"six"}))
		views, err := handler.Handle()
		testAssertErrNil(err)

		gomega.Expect(len(views.Items)).Should(gomega.Equal(0))

		gomega.Expect(views.PaginationInfo.Next).Should(gomega.BeEmpty())
		gomega.Expect(views.PaginationInfo.Prev).Should(gomega.BeEmpty())
	})
})
