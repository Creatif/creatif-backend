package addToList

import (
	"creatif/pkg/app/auth"
	"creatif/pkg/app/services/shared"
	"creatif/pkg/lib/logger"
	"creatif/pkg/lib/storage"
	"fmt"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Declaration (ADD) list entry tests", func() {
	ginkgo.It("should add an entry to a list by name", func() {
		projectId := testCreateProject("project")
		groups := testCreateGroups(projectId, 5)
		m := testCreateList(projectId, "listName")
		reference := testCreateList(projectId, "referenceMap")

		listVariables := make([]View, 0)
		for i := 0; i < 10; i++ {
			listVariables = append(listVariables, testAddToList(projectId, m.ID, fmt.Sprintf("name-%d", i), []shared.Reference{}, groups))
		}

		handler := New(NewModel(projectId, reference.ID, VariableModel{
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
				StructureType: "list",
				VariableID:    listVariables[0].ID,
			},
			{
				Name:          "second",
				StructureName: reference.Name,
				StructureType: "list",
				VariableID:    listVariables[1].ID,
			},
			{
				Name:          "third",
				StructureName: reference.Name,
				StructureType: "list",
				VariableID:    listVariables[2].ID,
			},
		}), auth.NewTestingAuthentication(false, ""), logger.NewLogBuilder())

		_, err := handler.Handle()
		gomega.Expect(err).Should(gomega.BeNil())

		var count int
		res := storage.Gorm().Raw("SELECT count(id) AS count FROM declarations.references").Scan(&count)
		gomega.Expect(res.Error).Should(gomega.BeNil())
		gomega.Expect(count).Should(gomega.Equal(3))
	})

	ginkgo.It("should fail to add an entry because of a duplicate reference", func() {
		projectId := testCreateProject("project")
		m := testCreateList(projectId, "mapName")
		groups := testCreateGroups(projectId, 5)
		reference := testCreateList(projectId, "referenceMap")

		listVariables := make([]View, 0)
		for i := 0; i < 10; i++ {
			listVariables = append(listVariables, testAddToList(projectId, m.ID, fmt.Sprintf("name-%d", i), []shared.Reference{}, groups))
		}

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
				StructureType: "list",
				VariableID:    listVariables[0].ID,
			},
			{
				Name:          "second",
				StructureName: reference.Name,
				StructureType: "list",
				VariableID:    listVariables[0].ID,
			},
			{
				Name:          "first",
				StructureName: reference.Name,
				StructureType: "list",
				VariableID:    listVariables[2].ID,
			},
		}), auth.NewTestingAuthentication(false, ""), logger.NewLogBuilder())

		_, err := handler.Handle()
		gomega.Expect(err).ShouldNot(gomega.BeNil())
	})

	ginkgo.It("should add an entry to a list by id", func() {
		projectId := testCreateProject("project")
		groups := testCreateGroups(projectId, 5)
		m := testCreateList(projectId, "mapName")

		listVariables := make([]View, 0)
		for i := 0; i < 10; i++ {
			listVariables = append(listVariables, testAddToList(projectId, m.ID, fmt.Sprintf("name-%d", i), []shared.Reference{}, groups))
		}

		handler := New(NewModel(projectId, m.ShortID, VariableModel{
			Name:      "newEntry",
			Metadata:  nil,
			Groups:    nil,
			Locale:    "eng",
			Behaviour: "readonly",
			Value:     nil,
		}, []shared.Reference{}), auth.NewTestingAuthentication(false, ""), logger.NewLogBuilder())

		_, err := handler.Handle()
		testAssertErrNil(err)
	})

	ginkgo.It("should add an entry to a list by shortID", func() {
		projectId := testCreateProject("project")
		m := testCreateList(projectId, "mapName")
		groups := testCreateGroups(projectId, 5)

		listVariables := make([]View, 0)
		for i := 0; i < 10; i++ {
			listVariables = append(listVariables, testAddToList(projectId, m.ID, fmt.Sprintf("name-%d", i), []shared.Reference{}, groups))
		}

		handler := New(NewModel(projectId, m.ID, VariableModel{
			Name:      "newEntry",
			Metadata:  nil,
			Groups:    nil,
			Locale:    "eng",
			Behaviour: "readonly",
			Value:     nil,
		}, []shared.Reference{}), auth.NewTestingAuthentication(false, ""), logger.NewLogBuilder())

		_, err := handler.Handle()
		testAssertErrNil(err)
	})
})
