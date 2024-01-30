package deleteRangeByID

import (
	"creatif/pkg/app/auth"
	declarations2 "creatif/pkg/app/domain/declarations"
	"creatif/pkg/lib/logger"
	"creatif/pkg/lib/sdk"
	"creatif/pkg/lib/storage"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Declaration map item delete tests", func() {
	ginkgo.It("should delete a range of map items by name", func() {
		projectId := testCreateProject("project")
		mapView := testCreateMap(projectId, "name", 15)

		var listItems []declarations2.MapVariable
		res := storage.Gorm().Where("map_id = ?", mapView.ID).Select("ID").Limit(10).Find(&listItems)
		gomega.Expect(res.Error).Should(gomega.BeNil())

		ids := sdk.Map(listItems, func(idx int, value declarations2.MapVariable) string {
			return value.ID
		})

		handler := New(NewModel(projectId, mapView.Name, ids), auth.NewTestingAuthentication(false, ""), logger.NewLogBuilder())
		model, err := handler.Handle()
		testAssertErrNil(err)
		gomega.Expect(model).Should(gomega.BeNil())

		var remainingItems []declarations2.MapVariable
		res = storage.Gorm().Where("map_id = ?", mapView.ID).Select("ID").Find(&remainingItems)
		gomega.Expect(res.Error).Should(gomega.BeNil())
		gomega.Expect(len(remainingItems)).Should(gomega.Equal(5))
	})

	ginkgo.It("should delete a range of map items by ID", func() {
		projectId := testCreateProject("project")
		mapView := testCreateMap(projectId, "name", 15)

		var listItems []declarations2.MapVariable
		res := storage.Gorm().Where("map_id = ?", mapView.ID).Select("ID").Limit(10).Find(&listItems)
		gomega.Expect(res.Error).Should(gomega.BeNil())

		ids := sdk.Map(listItems, func(idx int, value declarations2.MapVariable) string {
			return value.ID
		})

		handler := New(NewModel(projectId, mapView.ID, ids), auth.NewTestingAuthentication(false, ""), logger.NewLogBuilder())
		model, err := handler.Handle()
		testAssertErrNil(err)
		gomega.Expect(model).Should(gomega.BeNil())

		var remainingItems []declarations2.MapVariable
		res = storage.Gorm().Where("map_id = ?", mapView.ID).Select("ID").Find(&remainingItems)
		gomega.Expect(res.Error).Should(gomega.BeNil())
		gomega.Expect(len(remainingItems)).Should(gomega.Equal(5))
	})

	ginkgo.It("should delete a range of map items by shortID", func() {
		projectId := testCreateProject("project")
		mapView := testCreateMap(projectId, "name", 15)

		var listItems []declarations2.MapVariable
		res := storage.Gorm().Where("map_id = ?", mapView.ID).Select("ID").Limit(10).Find(&listItems)
		gomega.Expect(res.Error).Should(gomega.BeNil())

		ids := sdk.Map(listItems, func(idx int, value declarations2.MapVariable) string {
			return value.ID
		})

		handler := New(NewModel(projectId, mapView.ShortID, ids), auth.NewTestingAuthentication(false, ""), logger.NewLogBuilder())
		model, err := handler.Handle()
		testAssertErrNil(err)
		gomega.Expect(model).Should(gomega.BeNil())

		var remainingItems []declarations2.MapVariable
		res = storage.Gorm().Where("map_id = ?", mapView.ID).Select("ID").Find(&remainingItems)
		gomega.Expect(res.Error).Should(gomega.BeNil())
		gomega.Expect(len(remainingItems)).Should(gomega.Equal(5))
	})
})
