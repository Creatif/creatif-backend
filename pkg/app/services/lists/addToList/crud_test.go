package addToList

import (
	"creatif/pkg/app/auth"
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/app/services/shared"
	"creatif/pkg/lib/logger"
	"creatif/pkg/lib/storage"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Declaration (ADD) list entry tests", func() {
	ginkgo.It("should add an entry to a list by name", func() {
		projectId := testCreateProject("project")
		m := testCreateList(projectId, "listName", 10)
		reference := testCreateList(projectId, "referenceMap", 10)

		var listVariables []declarations.ListVariable
		res := storage.Gorm().Where("list_id = ?", reference.ID).Find(&listVariables)
		testAssertErrNil(res.Error)

		handler := New(NewModel(projectId, m.ID, VariableModel{
			Name:      "newEntry",
			Metadata:  nil,
			Groups:    nil,
			Locale:    "eng",
			Behaviour: "readonly",
			Value:     nil,
		}, []shared.Reference{
			{
				Name:          "first",
				StructureName: reference.Name,
				StructureType: "map",
				VariableID:    listVariables[0].ID,
			},
			{
				Name:          "second",
				StructureName: reference.Name,
				StructureType: "map",
				VariableID:    listVariables[1].ID,
			},
			{
				Name:          "third",
				StructureName: reference.Name,
				StructureType: "map",
				VariableID:    listVariables[2].ID,
			},
		}), auth.NewTestingAuthentication(false), logger.NewLogBuilder())

		_, err := handler.Handle()
		testAssertErrNil(err)
	})

	ginkgo.It("should fail to add an entry because of a duplicate", func() {
		projectId := testCreateProject("project")
		m := testCreateList(projectId, "mapName", 10)
		reference := testCreateList(projectId, "referenceMap", 10)

		var listVariables []declarations.ListVariable
		res := storage.Gorm().Where("list_id = ?", reference.ID).Find(&listVariables)
		testAssertErrNil(res.Error)

		handler := New(NewModel(projectId, m.ShortID, VariableModel{
			Name:      "newEntry",
			Metadata:  nil,
			Groups:    nil,
			Locale:    "eng",
			Behaviour: "readonly",
			Value:     nil,
		}, []shared.Reference{
			{
				Name:          "first",
				StructureName: reference.Name,
				StructureType: "map",
				VariableID:    listVariables[0].ID,
			},
			{
				Name:          "second",
				StructureName: reference.Name,
				StructureType: "map",
				VariableID:    listVariables[0].ID,
			},
			{
				Name:          "first",
				StructureName: reference.Name,
				StructureType: "map",
				VariableID:    listVariables[2].ID,
			},
		}), auth.NewTestingAuthentication(false), logger.NewLogBuilder())

		_, err := handler.Handle()
		gomega.Expect(err).ShouldNot(gomega.BeNil())
	})

	ginkgo.It("should add an entry to the map by id", func() {
		projectId := testCreateProject("project")
		m := testCreateList(projectId, "mapName", 10)

		handler := New(NewModel(projectId, m.ShortID, VariableModel{
			Name:      "newEntry",
			Metadata:  nil,
			Groups:    nil,
			Locale:    "eng",
			Behaviour: "readonly",
			Value:     nil,
		}, []shared.Reference{}), auth.NewTestingAuthentication(false), logger.NewLogBuilder())

		_, err := handler.Handle()
		testAssertErrNil(err)
	})

	ginkgo.It("should add an entry to the map by shortID", func() {
		projectId := testCreateProject("project")
		m := testCreateList(projectId, "mapName", 10)

		handler := New(NewModel(projectId, m.ID, VariableModel{
			Name:      "newEntry",
			Metadata:  nil,
			Groups:    nil,
			Locale:    "eng",
			Behaviour: "readonly",
			Value:     nil,
		}, []shared.Reference{}), auth.NewTestingAuthentication(false), logger.NewLogBuilder())

		_, err := handler.Handle()
		testAssertErrNil(err)
	})
})
