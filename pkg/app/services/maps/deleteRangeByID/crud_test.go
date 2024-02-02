package deleteRangeByID

import (
	"creatif/pkg/app/auth"
	declarations2 "creatif/pkg/app/domain/declarations"
	"creatif/pkg/app/services/maps/addToMap"
	"creatif/pkg/app/services/shared"
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
		referenceView := testCreateMap(projectId, "referenceView", 15)

		addedMapsWithReferences := make([]addToMap.LogicModel, 0)
		for i := 0; i < 10; i++ {
			addToMapVariable := testAddToMap(projectId, mapView.ID, []shared.Reference{
				{
					StructureName: referenceView.Name,
					StructureType: "map",
					VariableID:    referenceView.Variables[0].ID,
				},
				{
					StructureName: referenceView.Name,
					StructureType: "map",
					VariableID:    referenceView.Variables[1].ID,
				},
			})

			addedMapsWithReferences = append(addedMapsWithReferences, addToMapVariable)
		}

		var listItems []declarations2.MapVariable
		res := storage.Gorm().Where("map_id = ?", mapView.ID).Select("ID").Find(&listItems)
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
		gomega.Expect(len(remainingItems)).Should(gomega.Equal(0))

		var count int
		res = storage.Gorm().Raw("SELECT count(id) AS count FROM declarations.references").Scan(&count)
		testAssertErrNil(res.Error)
		gomega.Expect(count).Should(gomega.Equal(0))
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
