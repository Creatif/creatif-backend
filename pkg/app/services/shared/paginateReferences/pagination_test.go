package paginateReferences

import (
	"creatif/pkg/app/auth"
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/app/services/shared"
	"creatif/pkg/lib/logger"
	"creatif/pkg/lib/storage"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Map variables pagination tests", func() {
	ginkgo.It("should paginate through map variables", func() {
		projectId := testCreateProject("project")
		mapView, _ := testCreateMap(projectId, "name", 100)
		referenceView, _ := testCreateMap(projectId, "referenceMap", 100)
		_ = testAddToMap(projectId, mapView.Name, []shared.Reference{
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

		var references []declarations.Reference
		res := storage.Gorm().Raw("SELECT parent_id, child_id FROM declarations.references").Scan(&references)
		testAssertErrNil(res.Error)
		gomega.Expect(len(references)).Should(gomega.Equal(2))

		handler := New(NewModel(
			projectId,
			references[0].ParentID,
			references[1].ChildID,
			"map",
			[]string{},
			"created_at",
			"",
			"desc",
			10,
			1,
			[]string{"one"},
			nil,
			"",
			[]string{},
		), auth.NewTestingAuthentication(false), logger.NewLogBuilder())
		views, err := handler.Handle()
		testAssertErrNil(err)

		gomega.Expect(len(views.Data)).Should(gomega.Equal(2))
		gomega.Expect(views.Total).Should(gomega.Equal(int64(50)))
	})
})
