package removeMapVariable

import (
	"creatif/pkg/app/auth"
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/app/services/shared"
	"creatif/pkg/lib/logger"
	"creatif/pkg/lib/storage"
	"errors"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"gorm.io/gorm"
)

var _ = ginkgo.Describe("Declaration (DELETE) map entry tests", func() {
	ginkgo.It("should delete a map entry by name", func() {
		projectId := testCreateProject("project")
		view := testCreateMap(projectId, "mapName", 10)
		referenceView := testCreateMap(projectId, "referenceMap", 10)
		addToMapVariable := testAddToMap(projectId, view.Name, []shared.Reference{
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

		handler := New(NewModel(projectId, "mapName", addToMapVariable.Variable.ShortID), auth.NewTestingAuthentication(false), logger.NewLogBuilder())

		_, err := handler.Handle()
		testAssertErrNil(err)

		res := storage.Gorm().Where("map_id = ? AND short_id = ?", view.ID, addToMapVariable.Variable.ShortID).First(&declarations.MapVariable{})
		gomega.Expect(errors.Is(res.Error, gorm.ErrRecordNotFound)).Should(gomega.BeTrue())

		var count int
		res = storage.Gorm().Raw("SELECT count(id) AS count FROM declarations.references").Scan(&count)
		testAssertErrNil(res.Error)
		gomega.Expect(count).Should(gomega.Equal(0))
	})

	ginkgo.It("should delete a map entry multiple fields", func() {
		projectId := testCreateProject("project")
		view := testCreateMap(projectId, "mapName", 10)

		variables := view.Variables
		variableName := variables[0].ID
		handler := New(NewModel(projectId, view.ShortID, variableName), auth.NewTestingAuthentication(false), logger.NewLogBuilder())

		_, err := handler.Handle()
		testAssertErrNil(err)

		res := storage.Gorm().Where("map_id = ? AND id = ?", view.ID, variableName).First(&declarations.MapVariable{})
		gomega.Expect(errors.Is(res.Error, gorm.ErrRecordNotFound)).Should(gomega.BeTrue())
	})

	ginkgo.It("should delete a map entry multiple fields (3)", func() {
		projectId := testCreateProject("project")
		view := testCreateMap(projectId, "mapName", 10)

		variables := view.Variables
		shortId := variables[0].ShortID
		handler := New(NewModel(projectId, view.Name, shortId), auth.NewTestingAuthentication(false), logger.NewLogBuilder())

		_, err := handler.Handle()
		testAssertErrNil(err)

		res := storage.Gorm().Where("map_id = ? AND short_id = ?", view.ID, shortId).First(&declarations.MapVariable{})
		gomega.Expect(errors.Is(res.Error, gorm.ErrRecordNotFound)).Should(gomega.BeTrue())
	})
})
