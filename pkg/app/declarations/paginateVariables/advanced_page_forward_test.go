package paginateVariables

import (
	"creatif/pkg/app/declarations/paginateVariables/pagination"
	"fmt"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Advanced variable paginateVariables tests", func() {
	ginkgo.It("should advance cursor multiple times desc", func() {
		ginkgo.Skip("")
		projectId := testCreateProject("project1")
		limit := 10
		for i := 0; i < 40; i++ {
			testCreateBasicDeclarationTextVariable(projectId, fmt.Sprintf("name-%d", i), "modifiable")
		}
		// 39-30
		// 29-19
		// 18-8
		var paginationId string
		paginationId = testAdvanceCursor(paginationId, projectId, pagination.DIRECTION_FORWARD, pagination.DESC, 10)
		paginationId = testAdvanceCursor(paginationId, projectId, pagination.DIRECTION_FORWARD, pagination.DESC, 10)
		paginationId = testAdvanceCursor(paginationId, projectId, pagination.DIRECTION_FORWARD, pagination.DESC, 10)

		handler := New(NewModel(projectId, paginationId, "created_at", pagination.DESC, pagination.DIRECTION_FORWARD, limit, []string{}))
		views, err := handler.Handle()
		testAssertErrNil(err)

		gomega.Expect(len(views.Items)).Should(gomega.Equal(limit))
		gomega.Expect(views.Items[0].Name).Should(gomega.Equal("name-9"))
		gomega.Expect(views.Items[len(views.Items)-1].Name).Should(gomega.Equal("name-0"))

		gomega.Expect(views.PaginationInfo.Next).ShouldNot(gomega.BeEmpty())
		gomega.Expect(views.PaginationInfo.Prev).Should(gomega.BeEmpty())
	})

	ginkgo.It("should advance cursor until the end and return empty results with the same paginationId desc", func() {
		ginkgo.Skip("")
		projectId := testCreateProject("project2")
		limit := 10
		for i := 0; i < 40; i++ {
			testCreateBasicDeclarationTextVariable(projectId, fmt.Sprintf("name-%d", i), "modifiable")
		}

		var paginationId string
		paginationId = testAdvanceCursor(paginationId, projectId, pagination.DIRECTION_FORWARD, pagination.DESC, 10)
		paginationId = testAdvanceCursor(paginationId, projectId, pagination.DIRECTION_FORWARD, pagination.DESC, 10)
		paginationId = testAdvanceCursor(paginationId, projectId, pagination.DIRECTION_FORWARD, pagination.DESC, 10)
		paginationId = testAdvanceCursor(paginationId, projectId, pagination.DIRECTION_FORWARD, pagination.DESC, 10)

		handler := New(NewModel(projectId, paginationId, "created_at", pagination.DESC, pagination.DIRECTION_FORWARD, limit, []string{}))
		views, err := handler.Handle()
		testAssertErrNil(err)

		gomega.Expect(len(views.Items)).Should(gomega.Equal(0))

		gomega.Expect(views.PaginationInfo.Parameters.PaginationID).Should(gomega.Equal(paginationId))
		gomega.Expect(views.PaginationInfo.Next).ShouldNot(gomega.BeEmpty())
		gomega.Expect(views.PaginationInfo.Prev).Should(gomega.BeEmpty())
	})

	ginkgo.It("should advance cursor until the end multiple times and still return the same pagination id desc", func() {
		ginkgo.Skip("")
		projectId := testCreateProject("project3")
		limit := 10
		for i := 0; i < 40; i++ {
			testCreateBasicDeclarationTextVariable(projectId, fmt.Sprintf("name-%d", i), "modifiable")
		}

		var paginationId string
		paginationId = testAdvanceCursor(paginationId, projectId, pagination.DIRECTION_FORWARD, pagination.DESC, 10)
		paginationId = testAdvanceCursor(paginationId, projectId, pagination.DIRECTION_FORWARD, pagination.DESC, 10)
		paginationId = testAdvanceCursor(paginationId, projectId, pagination.DIRECTION_FORWARD, pagination.DESC, 10)
		paginationId = testAdvanceCursor(paginationId, projectId, pagination.DIRECTION_FORWARD, pagination.DESC, 10)
		paginationId = testAdvanceCursor(paginationId, projectId, pagination.DIRECTION_FORWARD, pagination.DESC, 10)

		handler := New(NewModel(projectId, paginationId, "created_at", pagination.DESC, pagination.DIRECTION_FORWARD, limit, []string{}))
		views, err := handler.Handle()
		testAssertErrNil(err)

		gomega.Expect(len(views.Items)).Should(gomega.Equal(0))

		gomega.Expect(views.PaginationInfo.Parameters.PaginationID).Should(gomega.Equal(paginationId))
		gomega.Expect(views.PaginationInfo.Next).ShouldNot(gomega.BeEmpty())
		gomega.Expect(views.PaginationInfo.Prev).Should(gomega.BeEmpty())
	})

	ginkgo.It("should advance cursor multiple times asc", func() {
		ginkgo.Skip("")

		projectId := testCreateProject("project4")
		limit := 10
		for i := 0; i < 40; i++ {
			testCreateBasicDeclarationTextVariable(projectId, fmt.Sprintf("name-%d", i), "modifiable")
		}
		// 39-30
		// 29-19
		// 18-8
		nextPaginationId := testAdvanceCursor("", projectId, pagination.DIRECTION_FORWARD, pagination.ASC, 10)

		handler := New(NewModel(projectId, nextPaginationId, "created_at", pagination.ASC, pagination.DIRECTION_FORWARD, limit, []string{}))
		views, err := handler.Handle()
		testAssertErrNil(err)

		gomega.Expect(len(views.Items)).Should(gomega.Equal(limit))
		gomega.Expect(views.Items[0].Name).Should(gomega.Equal("name-30"))
		gomega.Expect(views.Items[len(views.Items)-1].Name).Should(gomega.Equal("name-39"))

		gomega.Expect(views.PaginationInfo.Next).ShouldNot(gomega.BeEmpty())
		gomega.Expect(views.PaginationInfo.Prev).Should(gomega.BeEmpty())
	})

	ginkgo.It("should advance cursor until the end and return empty results with the same paginationId asc", func() {
		ginkgo.Skip("")

		projectId := testCreateProject("project5")
		limit := 10
		for i := 0; i < 40; i++ {
			testCreateBasicDeclarationTextVariable(projectId, fmt.Sprintf("name-%d", i), "modifiable")
		}

		nextPaginationId := testAdvanceCursor("", projectId, pagination.DIRECTION_FORWARD, pagination.ASC, 10)

		handler := New(NewModel(projectId, nextPaginationId, "created_at", pagination.ASC, pagination.DIRECTION_FORWARD, limit, []string{}))
		views, err := handler.Handle()
		testAssertErrNil(err)

		gomega.Expect(len(views.Items)).Should(gomega.Equal(0))

		gomega.Expect(views.PaginationInfo.Parameters.PaginationID).Should(gomega.Equal(nextPaginationId))
		gomega.Expect(views.PaginationInfo.Next).ShouldNot(gomega.BeEmpty())
		gomega.Expect(views.PaginationInfo.Prev).Should(gomega.BeEmpty())
	})

	ginkgo.It("should advance cursor until the end multiple times and still return the same pagination id asc", func() {
		ginkgo.Skip("")

		projectId := testCreateProject("project6")
		limit := 10
		for i := 0; i < 40; i++ {
			testCreateBasicDeclarationTextVariable(projectId, fmt.Sprintf("name-%d", i), "modifiable")
		}
		nextPaginationId := testAdvanceCursor("", projectId, pagination.DIRECTION_FORWARD, pagination.ASC, 10)

		for i := 0; i < 10; i++ {
			handler := New(NewModel(projectId, nextPaginationId, "created_at", pagination.ASC, pagination.DIRECTION_FORWARD, limit, []string{}))
			views, err := handler.Handle()
			testAssertErrNil(err)

			gomega.Expect(len(views.Items)).Should(gomega.Equal(0))

			gomega.Expect(views.PaginationInfo.Parameters.PaginationID).Should(gomega.Equal(nextPaginationId))
			gomega.Expect(views.PaginationInfo.Next).ShouldNot(gomega.BeEmpty())
			gomega.Expect(views.PaginationInfo.Prev).Should(gomega.BeEmpty())
		}
	})
})
