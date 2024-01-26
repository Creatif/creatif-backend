package paginateReferences

import (
	"creatif/pkg/app/auth"
	"creatif/pkg/app/services/locales"
	"creatif/pkg/app/services/shared"
	"creatif/pkg/lib/logger"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Map variables pagination tests", func() {
	ginkgo.It("should paginate through map variables", func() {
		projectId := testCreateProject("project")
		mapView, _ := testCreateMap(projectId, "name", 100)
		referenceView, _ := testCreateMap(projectId, "referenceView", 100)
		addToMapVariable := testAddToMap(projectId, mapView.Name, []shared.Reference{
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

		localeId, _ := locales.GetIDWithAlpha("eng")
		handler := New(NewModel(projectId, referenceView.Variables[0].ID, addToMapVariable.Variable.ID, "parent", "map", []string{localeId}, "created_at", "", "desc", 10, 1, []string{"one"}, nil, "", []string{}), auth.NewTestingAuthentication(false), logger.NewLogBuilder())
		views, err := handler.Handle()
		testAssertErrNil(err)

		gomega.Expect(len(views.Data)).Should(gomega.Equal(2))
	})
})
