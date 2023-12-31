package paginateMapVariables

import (
	"creatif/pkg/app/auth"
	"creatif/pkg/app/services/locales"
	"creatif/pkg/lib/logger"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Map variables pagination tests", func() {
	ginkgo.It("should paginate through map variables", func() {
		projectId := testCreateProject("project")
		mapView, _ := testCreateMap(projectId, "name", 100)

		localeId, _ := locales.GetIDWithAlpha("eng")
		handler := New(NewModel(projectId, []string{localeId}, mapView.Name, "created_at", "", "desc", 10, 1, []string{"one"}, nil, "", []string{}), auth.NewTestingAuthentication(false), logger.NewLogBuilder())
		views, err := handler.Handle()
		testAssertErrNil(err)

		gomega.Expect(len(views.Data)).Should(gomega.Equal(10))
		gomega.Expect(views.Total).Should(gomega.Equal(int64(50)))
	})

	ginkgo.It("should get an empty result from the end of the map variables listing", func() {
		projectId := testCreateProject("project")
		mapView, _ := testCreateMap(projectId, "name", 100)

		handler := New(NewModel(projectId, []string{}, mapView.ID, "created_at", "", "desc", 10, 50, []string{}, nil, "", []string{}), auth.NewTestingAuthentication(false), logger.NewLogBuilder())
		views, err := handler.Handle()
		testAssertErrNil(err)

		gomega.Expect(len(views.Data)).Should(gomega.Equal(0))
		gomega.Expect(views.Total).Should(gomega.Equal(int64(100)))
	})

	ginkgo.It("should return empty result for group that does not exist", func() {
		projectId := testCreateProject("project")
		mapView, _ := testCreateMap(projectId, "name", 100)

		localeId, _ := locales.GetIDWithAlpha("eng")
		handler := New(NewModel(projectId, []string{localeId}, mapView.ShortID, "created_at", "", "desc", 10, 1, []string{"not_exists"}, nil, "", []string{}), auth.NewTestingAuthentication(false), logger.NewLogBuilder())
		views, err := handler.Handle()
		testAssertErrNil(err)

		gomega.Expect(len(views.Data)).Should(gomega.Equal(0))
		gomega.Expect(views.Total).Should(gomega.Equal(int64(0)))
	})

	ginkgo.It("should return the exact number of items by group", func() {
		projectId := testCreateProject("project")
		mapView, groups := testCreateMap(projectId, "name", 100)

		handler := New(NewModel(projectId, []string{}, mapView.Name, "created_at", "", "desc", 50, 1, []string{"one"}, nil, "", []string{}), auth.NewTestingAuthentication(false), logger.NewLogBuilder())
		views, err := handler.Handle()
		testAssertErrNil(err)

		gomega.Expect(len(views.Data)).Should(gomega.Equal(50))
		gomega.Expect(views.Total).Should(gomega.Equal(int64(groups["one"])))
	})

	ginkgo.It("should return items search by name with regex", func() {
		projectId := testCreateProject("project")
		mapView, _ := testCreateMap(projectId, "name", 100)

		handler := New(NewModel(projectId, []string{}, mapView.Name, "created_at", "1", "desc", 10, 1, []string{}, nil, "", []string{}), auth.NewTestingAuthentication(false), logger.NewLogBuilder())
		views, err := handler.Handle()
		testAssertErrNil(err)

		gomega.Expect(len(views.Data)).Should(gomega.Equal(10))
		gomega.Expect(views.Total).Should(gomega.Equal(int64(19)))
	})

	ginkgo.It("should return items search by name with regex with groups", func() {
		projectId := testCreateProject("project")
		mapView, _ := testCreateMap(projectId, "name", 100)

		localeId, _ := locales.GetIDWithAlpha("eng")
		handler := New(NewModel(projectId, []string{localeId}, mapView.ID, "created_at", "1", "desc", 10, 1, []string{"one"}, nil, "", []string{}), auth.NewTestingAuthentication(false), logger.NewLogBuilder())
		views, err := handler.Handle()
		testAssertErrNil(err)

		gomega.Expect(len(views.Data)).Should(gomega.Equal(5))
		gomega.Expect(views.Total).Should(gomega.Equal(int64(5)))
	})
})
