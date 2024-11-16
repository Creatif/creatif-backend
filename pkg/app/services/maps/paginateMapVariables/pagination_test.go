package paginateMapVariables

import (
	"creatif/pkg/app/auth"
	"creatif/pkg/app/services/groups/addGroups"
	"creatif/pkg/app/services/locales"
	"creatif/pkg/app/services/shared/connections"
	"creatif/pkg/lib/sdk"
	"fmt"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Map variables pagination tests", func() {
	ginkgo.It("should paginate through map variables", ginkgo.Label("map", "maps_pagination"), func() {
		projectId := testCreateProject("project")
		groups := testCreateGroups(projectId, 5)
		mapView := testCreateMap(projectId, "name")

		for i := 0; i < 100; i++ {
			testAddToMap(projectId, mapView.ID, fmt.Sprintf("name-%d", i), []connections.Connection{}, sdk.Map(groups, func(idx int, value addGroups.View) string {
				return value.ID
			}))
		}

		localeId, _ := locales.GetIDWithAlpha("eng")
		handler := New(NewModel(projectId, []string{localeId}, mapView.ID, "created_at", "", "desc", 10, 1, []string{groups[0].ID}, nil, "", []string{}), auth.NewTestingAuthentication(false, ""))
		views, err := handler.Handle()
		testAssertErrNil(err)

		gomega.Expect(len(views.Data)).Should(gomega.Equal(10))
	})

	ginkgo.It("should get an empty result from the end of the map variables listing", ginkgo.Label("map", "maps_pagination"), func() {
		projectId := testCreateProject("project")
		groups := testCreateGroups(projectId, 5)
		mapView := testCreateMap(projectId, "name")

		for i := 0; i < 50; i++ {
			testAddToMap(projectId, mapView.ID, fmt.Sprintf("name-%d", i), []connections.Connection{}, sdk.Map(groups, func(idx int, value addGroups.View) string {
				return value.ID
			}))
		}

		handler := New(NewModel(projectId, []string{}, mapView.ID, "created_at", "", "desc", 10, 50, []string{}, nil, "", []string{}), auth.NewTestingAuthentication(false, ""))
		views, err := handler.Handle()
		testAssertErrNil(err)

		gomega.Expect(len(views.Data)).Should(gomega.Equal(0))
	})

	ginkgo.It("should return empty result for group that does not exist", ginkgo.Label("map", "maps_pagination"), func() {
		projectId := testCreateProject("project")
		mapView := testCreateMap(projectId, "name")

		for i := 0; i < 100; i++ {
			testAddToMap(projectId, mapView.ID, fmt.Sprintf("name-%d", i), []connections.Connection{}, []string{})
		}

		localeId, _ := locales.GetIDWithAlpha("eng")
		handler := New(NewModel(projectId, []string{localeId}, mapView.ShortID, "created_at", "", "desc", 10, 1, []string{"not_exists"}, nil, "", []string{}), auth.NewTestingAuthentication(false, ""))
		views, err := handler.Handle()
		testAssertErrNil(err)

		gomega.Expect(len(views.Data)).Should(gomega.Equal(0))
	})

	ginkgo.It("should return the exact number of items by group", ginkgo.Label("map", "maps_pagination"), func() {
		projectId := testCreateProject("project")
		groups := testCreateGroups(projectId, 10)
		mapView := testCreateMap(projectId, "name")

		for i := 0; i < 100; i++ {
			testAddToMap(projectId, mapView.ID, fmt.Sprintf("name-%d", i), []connections.Connection{}, sdk.Map(groups, func(idx int, value addGroups.View) string {
				return value.ID
			}))
		}

		handler := New(NewModel(projectId, []string{}, mapView.ID, "created_at", "", "desc", 50, 1, []string{groups[0].ID, groups[1].ID}, nil, "", []string{}), auth.NewTestingAuthentication(false, ""))
		views, err := handler.Handle()
		testAssertErrNil(err)

		gomega.Expect(len(views.Data)).Should(gomega.Equal(50))
	})

	ginkgo.It("should return items search by name with regex", ginkgo.Label("map", "maps_pagination"), func() {
		projectId := testCreateProject("project")
		mapView := testCreateMap(projectId, "name")

		for i := 0; i < 100; i++ {
			testAddToMap(projectId, mapView.ID, fmt.Sprintf("name-%d", i), []connections.Connection{}, []string{})
		}

		handler := New(NewModel(projectId, []string{}, mapView.ID, "created_at", "1", "desc", 10, 1, []string{}, nil, "", []string{}), auth.NewTestingAuthentication(false, ""))
		views, err := handler.Handle()
		testAssertErrNil(err)

		gomega.Expect(len(views.Data)).Should(gomega.Equal(10))
	})

	ginkgo.It("should return items search by name with regex with groups", ginkgo.Label("map", "maps_pagination"), func() {
		projectId := testCreateProject("project")
		groups := testCreateGroups(projectId, 10)
		mapView := testCreateMap(projectId, "name")

		for i := 0; i < 100; i++ {
			testAddToMap(projectId, mapView.ID, fmt.Sprintf("name-%d", i), []connections.Connection{}, sdk.Map(groups, func(idx int, value addGroups.View) string {
				return value.ID
			}))
		}

		localeId, _ := locales.GetIDWithAlpha("eng")
		handler := New(NewModel(projectId, []string{localeId}, mapView.ID, "created_at", "1", "desc", 10, 1, []string{groups[0].ID, groups[1].ID}, nil, "", []string{}), auth.NewTestingAuthentication(false, ""))
		views, err := handler.Handle()
		testAssertErrNil(err)

		gomega.Expect(len(views.Data)).Should(gomega.Equal(10))
	})
})
