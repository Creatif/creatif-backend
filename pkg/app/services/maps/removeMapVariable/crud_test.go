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
	ginkgo.It("should delete a map entry by name", ginkgo.Label("map"), func() {
		projectId := testCreateProject("project")
		view := testCreateMap(projectId, "mapName")
		referenceView := testCreateMap(projectId, "referenceMap")

		referenceVar1 := testAddToMap(projectId, referenceView.ID, []shared.Reference{})
		referenceVar2 := testAddToMap(projectId, referenceView.ID, []shared.Reference{})

		addToMapVariable := testAddToMap(projectId, view.ID, []shared.Reference{
			{
				Name:          "one",
				StructureName: referenceView.Name,
				StructureType: "map",
				VariableID:    referenceVar1.Variable.ID,
			},
			{
				Name:          "two",
				StructureName: referenceView.Name,
				StructureType: "map",
				VariableID:    referenceVar2.Variable.ID,
			},
		})

		handler := New(NewModel(projectId, view.ID, addToMapVariable.Variable.ID), auth.NewTestingAuthentication(false, ""), logger.NewLogBuilder())

		_, err := handler.Handle()
		testAssertErrNil(err)

		res := storage.Gorm().Where("map_id = ? AND id = ?", view.ID, addToMapVariable.Variable.ID).First(&declarations.MapVariable{})
		gomega.Expect(errors.Is(res.Error, gorm.ErrRecordNotFound)).Should(gomega.BeTrue())

		var count int
		res = storage.Gorm().Raw("SELECT count(id) AS count FROM declarations.references").Scan(&count)
		testAssertErrNil(res.Error)
		gomega.Expect(count).Should(gomega.Equal(0))
	})
})
