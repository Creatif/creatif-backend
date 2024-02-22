package paginateListItems

import (
	"creatif/pkg/app/auth"
	"creatif/pkg/app/services/groups/addGroups"
	"creatif/pkg/app/services/locales"
	"creatif/pkg/app/services/shared"
	"creatif/pkg/lib/logger"
	"creatif/pkg/lib/sdk"
	"fmt"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("List items pagination tests", func() {
	ginkgo.It("should paginate through list variables", func() {
		projectId := testCreateProject("project")
		groups := testCreateGroups(projectId, 0)
		list := testCreateList(projectId, "name")

		for i := 0; i < 100; i++ {
			testAddToList(projectId, list.ID, fmt.Sprintf("name-%d", i), []shared.Reference{}, sdk.Map(groups, func(idx int, value addGroups.View) string {
				return value.ID
			}))
		}

		localeId, _ := locales.GetIDWithAlpha("eng")
		handler := New(NewModel(projectId, []string{localeId}, list.ID, "created_at", "", "desc", 10, 1, []string{}, nil, "", []string{}), auth.NewTestingAuthentication(false, ""), logger.NewLogBuilder())
		views, err := handler.Handle()
		testAssertErrNil(err)

		gomega.Expect(len(views.Data)).Should(gomega.Equal(10))
		gomega.Expect(views.Total).Should(gomega.Equal(int64(100)))
	})

	ginkgo.It("should get an empty result from the end of the list variables listing", func() {
		projectId := testCreateProject("project")
		groups := testCreateGroups(projectId, 5)
		list := testCreateList(projectId, "name")

		for i := 0; i < 50; i++ {
			testAddToList(projectId, list.ID, fmt.Sprintf("name-%d", i), []shared.Reference{}, sdk.Map(groups, func(idx int, value addGroups.View) string {
				return value.ID
			}))
		}

		handler := New(NewModel(projectId, []string{}, list.ID, "created_at", "", "desc", 10, 50, []string{}, nil, "", []string{}), auth.NewTestingAuthentication(false, ""), logger.NewLogBuilder())
		views, err := handler.Handle()
		testAssertErrNil(err)

		gomega.Expect(len(views.Data)).Should(gomega.Equal(0))
		gomega.Expect(views.Total).Should(gomega.Equal(int64(50)))
	})

	ginkgo.It("should return empty result for group that does not exist", func() {
		projectId := testCreateProject("project")
		groups := testCreateGroups(projectId, 5)
		list := testCreateList(projectId, "name")

		for i := 0; i < 100; i++ {
			testAddToList(projectId, list.ID, fmt.Sprintf("name-%d", i), []shared.Reference{}, sdk.Map(groups, func(idx int, value addGroups.View) string {
				return value.ID
			}))
		}

		localeId, _ := locales.GetIDWithAlpha("eng")
		handler := New(NewModel(projectId, []string{localeId}, list.ID, "created_at", "", "desc", 10, 1, []string{"not_exists"}, nil, "", []string{}), auth.NewTestingAuthentication(false, ""), logger.NewLogBuilder())
		views, err := handler.Handle()
		testAssertErrNil(err)

		gomega.Expect(len(views.Data)).Should(gomega.Equal(0))
		gomega.Expect(views.Total).Should(gomega.Equal(int64(0)))
	})

	ginkgo.It("should return the exact number of items by group", func() {
		projectId := testCreateProject("project")
		groups := testCreateGroups(projectId, 5)
		list := testCreateList(projectId, "name")

		for i := 0; i < 100; i++ {
			testAddToList(projectId, list.ID, fmt.Sprintf("name-%d", i), []shared.Reference{}, sdk.Map(groups, func(idx int, value addGroups.View) string {
				return value.ID
			}))
		}

		handler := New(NewModel(projectId, []string{}, list.ID, "created_at", "", "desc", 50, 1, []string{groups[0].ID}, nil, "", []string{}), auth.NewTestingAuthentication(false, ""), logger.NewLogBuilder())
		views, err := handler.Handle()
		testAssertErrNil(err)

		gomega.Expect(len(views.Data)).Should(gomega.Equal(50))
		gomega.Expect(views.Total).Should(gomega.Equal(int64(100)))
	})

	ginkgo.It("should return items search by name with regex", func() {
		projectId := testCreateProject("project")
		groups := testCreateGroups(projectId, 5)
		list := testCreateList(projectId, "name")

		for i := 0; i < 100; i++ {
			testAddToList(projectId, list.ID, fmt.Sprintf("name-%d", i), []shared.Reference{}, sdk.Map(groups, func(idx int, value addGroups.View) string {
				return value.ID
			}))
		}

		handler := New(NewModel(projectId, []string{}, list.ID, "created_at", "1", "desc", 10, 1, []string{}, nil, "", []string{}), auth.NewTestingAuthentication(false, ""), logger.NewLogBuilder())
		views, err := handler.Handle()
		testAssertErrNil(err)

		gomega.Expect(len(views.Data)).Should(gomega.Equal(10))
		gomega.Expect(views.Total).Should(gomega.Equal(int64(19)))
	})

	ginkgo.It("should return items search by name with regex with groups", func() {
		projectId := testCreateProject("project")
		groups := testCreateGroups(projectId, 5)
		list := testCreateList(projectId, "name")

		for i := 0; i < 100; i++ {
			testAddToList(projectId, list.ID, fmt.Sprintf("name-%d", i), []shared.Reference{}, sdk.Map(groups, func(idx int, value addGroups.View) string {
				return value.ID
			}))
		}

		localeId, _ := locales.GetIDWithAlpha("eng")
		handler := New(NewModel(projectId, []string{localeId}, list.ID, "created_at", "1", "desc", 10, 1, []string{groups[0].ID}, nil, "", []string{}), auth.NewTestingAuthentication(false, ""), logger.NewLogBuilder())
		views, err := handler.Handle()
		testAssertErrNil(err)

		gomega.Expect(len(views.Data)).Should(gomega.Equal(10))
		gomega.Expect(views.Total).Should(gomega.Equal(int64(19)))
	})
})
