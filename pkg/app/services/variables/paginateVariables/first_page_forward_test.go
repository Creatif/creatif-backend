package paginateVariables

import (
	"creatif/pkg/app/services/variables/paginateVariables/pagination"
	"fmt"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("First page variable paginateVariables tests", func() {
	ginkgo.It("should return the first row of results by created_at field desc going forward", func() {
		projectId := testCreateProject("project")
		limit := 10
		for i := 0; i < 40; i++ {
			testCreateBasicDeclarationTextVariable(projectId, fmt.Sprintf("name-%d", i), "modifiable")
		}

		handler := New(NewModel(projectId, "", "created_at", pagination.DESC, pagination.DIRECTION_FORWARD, limit, []string{}))
		views, err := handler.Handle()
		testAssertErrNil(err)

		gomega.Expect(len(views.Items)).Should(gomega.Equal(limit))
		gomega.Expect(views.Items[0].Name).Should(gomega.Equal("name-39"))
		gomega.Expect(views.Items[len(views.Items)-1].Name).Should(gomega.Equal("name-30"))

		gomega.Expect(views.PaginationInfo.Next).ShouldNot(gomega.BeEmpty())
		gomega.Expect(views.PaginationInfo.Prev).Should(gomega.BeEmpty())
	})

	ginkgo.It("should return the first row of results by created_at field desc going forward and with specifying all the groups and some that do not exist", func() {
		projectId := testCreateProject("project")
		limit := 10
		for i := 0; i < 20; i++ {
			testCreateBasicDeclarationTextVariable(projectId, fmt.Sprintf("name-%d", i), "modifiable")
		}

		handler := New(NewModel(projectId, "", "created_at", pagination.DESC, pagination.DIRECTION_FORWARD, limit, []string{}))
		views, err := handler.Handle()
		testAssertErrNil(err)

		gomega.Expect(len(views.Items)).Should(gomega.Equal(limit))
		gomega.Expect(views.Items[0].Name).Should(gomega.Equal("name-19"))
		gomega.Expect(views.Items[len(views.Items)-1].Name).Should(gomega.Equal("name-10"))

		gomega.Expect(views.PaginationInfo.Next).ShouldNot(gomega.BeEmpty())
		gomega.Expect(views.PaginationInfo.Prev).Should(gomega.BeEmpty())
	})

	ginkgo.It("should return the first row of results by created_at field asc going forward", func() {
		projectId := testCreateProject("project")
		limit := 10
		for i := 0; i < 20; i++ {
			testCreateBasicDeclarationTextVariable(projectId, fmt.Sprintf("name-%d", i), "modifiable")
		}

		handler := New(NewModel(projectId, "", "created_at", pagination.ASC, pagination.DIRECTION_FORWARD, limit, []string{"one"}))
		views, err := handler.Handle()
		testAssertErrNil(err)

		gomega.Expect(len(views.Items)).Should(gomega.Equal(limit))
		gomega.Expect(views.Items[0].Name).Should(gomega.Equal("name-0"))
		gomega.Expect(views.Items[len(views.Items)-1].Name).Should(gomega.Equal("name-9"))

		gomega.Expect(views.PaginationInfo.Next).ShouldNot(gomega.BeEmpty())
		gomega.Expect(views.PaginationInfo.Prev).Should(gomega.BeEmpty())
	})

	ginkgo.It("should give an empty result for a non existent group", func() {
		projectId := testCreateProject("project")
		limit := 10
		for i := 0; i < 5; i++ {
			testCreateBasicDeclarationTextVariable(projectId, fmt.Sprintf("name-%d", i), "modifiable")
		}

		handler := New(NewModel(projectId, "", "created_at", pagination.ASC, pagination.DIRECTION_FORWARD, limit, []string{"six"}))
		views, err := handler.Handle()
		testAssertErrNil(err)

		gomega.Expect(len(views.Items)).Should(gomega.Equal(0))

		gomega.Expect(views.PaginationInfo.Next).Should(gomega.BeEmpty())
		gomega.Expect(views.PaginationInfo.Prev).Should(gomega.BeEmpty())
	})
})
