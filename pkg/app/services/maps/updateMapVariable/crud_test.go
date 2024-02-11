package updateMapVariable

import (
	"creatif/pkg/app/auth"
	declarations2 "creatif/pkg/app/domain/declarations"
	"creatif/pkg/app/services/maps/addToMap"
	"creatif/pkg/app/services/shared"
	"creatif/pkg/lib/appErrors"
	"creatif/pkg/lib/logger"
	"creatif/pkg/lib/storage"
	"encoding/json"
	"fmt"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"gorm.io/datatypes"
)

var _ = ginkgo.Describe("Declaration (UPDATE) map entry tests", func() {
	ginkgo.It("should update an entry in the map by replacing it completely", ginkgo.Label("map", "map_update"), func() {
		projectId := testCreateProject("project")
		groups := testCreateGroups(projectId, []string{"one", "two", "three", "four", "five"})
		m := testCreateMap(projectId, "map")
		referenceMap := testCreateMap(projectId, "referenceMap")

		referenceVar1 := testAddToMap(projectId, referenceMap.ID, "name-1", []shared.Reference{}, groups, "")
		referenceVar2 := testAddToMap(projectId, referenceMap.ID, "name-2", []shared.Reference{}, groups, "")
		referenceVar3 := testAddToMap(projectId, referenceMap.ID, "name-3", []shared.Reference{}, groups, "")

		referenceVar4 := testAddToMap(projectId, referenceMap.ID, "name-4", []shared.Reference{}, groups, "")
		referenceVar5 := testAddToMap(projectId, referenceMap.ID, "name-5", []shared.Reference{}, groups, "")

		addToMapView := testAddToMap(projectId, m.ID, "name-0", []shared.Reference{
			{
				Name:          "myName",
				StructureName: referenceMap.Name,
				StructureType: "map",
				VariableID:    referenceVar1.ID,
			},
			{
				Name:          "another",
				StructureName: referenceMap.Name,
				StructureType: "map",
				VariableID:    referenceVar2.ID,
			},
			{
				Name:          "three",
				StructureName: referenceMap.Name,
				StructureType: "map",
				VariableID:    referenceVar3.ID,
			},
		}, groups, "")

		b, err := json.Marshal("this is metadata")
		gomega.Expect(err).Should(gomega.BeNil())

		v, err := json.Marshal("this is value")
		gomega.Expect(err).Should(gomega.BeNil())

		handler := New(NewModel(projectId, m.ShortID, addToMapView.ID, []string{"metadata", "groups", "behaviour", "value", "name"}, VariableModel{
			Name:      "new name",
			Metadata:  b,
			Groups:    []string{groups[1], groups[2]},
			Behaviour: "readonly",
			Value:     v,
		}, []shared.UpdateReference{
			{
				Name:          addToMapView.References[0].Name,
				StructureName: "referenceMap",
				StructureType: "map",
				VariableID:    referenceVar4.ID,
			},
			{
				Name:          "new entry",
				StructureName: "referenceMap",
				StructureType: "map",
				VariableID:    referenceVar5.ID,
			},
		}), auth.NewTestingAuthentication(false, ""), logger.NewLogBuilder())

		view, err := handler.Handle()

		testAssertErrNil(err)
		testAssertIDValid(view.ID)

		var metadata string
		gomega.Expect(json.Unmarshal(view.Metadata.(datatypes.JSON), &metadata)).Should(gomega.BeNil())

		var value string
		gomega.Expect(json.Unmarshal(view.Value.(datatypes.JSON), &value)).Should(gomega.BeNil())

		gomega.Expect(view.Name).Should(gomega.Equal("new name"))
		gomega.Expect(metadata).Should(gomega.Equal("this is metadata"))
		gomega.Expect(value).Should(gomega.Equal("this is value"))
		gomega.Expect(view.Behaviour).Should(gomega.Equal("readonly"))

		var count int
		res := storage.Gorm().Raw("SELECT count(id) AS count FROM declarations.references").Scan(&count)
		testAssertErrNil(res.Error)
		gomega.Expect(count).Should(gomega.Equal(2))

		var groupCount int
		res = storage.Gorm().Raw(fmt.Sprintf("SELECT COUNT(variable_id) FROM %s WHERE variable_id = ? GROUP BY variable_id", (declarations2.VariableGroup{}).TableName()), addToMapView.ID).Scan(&groupCount)
		gomega.Expect(res.Error).Should(gomega.BeNil())

		gomega.Expect(count).Should(gomega.Equal(2))
	})

	ginkgo.It("should fail updating a map variable because of invalid number of groups", ginkgo.Label("map"), func() {
		projectId := testCreateProject("project")
		groups := testCreateGroups(projectId, []string{"one", "two", "three"})
		m := testCreateMap(projectId, "map")

		variables := make([]addToMap.View, 0)
		for i := 0; i < 10; i++ {
			variables = append(variables, testAddToMap(projectId, m.ID, fmt.Sprintf("name-%d", i), []shared.Reference{}, groups, ""))
		}

		b, err := json.Marshal("this is metadata")
		gomega.Expect(err).Should(gomega.BeNil())

		v, err := json.Marshal("this is value")
		gomega.Expect(err).Should(gomega.BeNil())

		handler := New(NewModel(projectId, m.ID, variables[5].ID, []string{"metadata", "groups", "behaviour", "value"}, VariableModel{
			Name:      "new name",
			Metadata:  b,
			Groups:    []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15", "16", "17", "18"},
			Behaviour: "readonly",
			Value:     v,
		}, nil), auth.NewTestingAuthentication(false, ""), logger.NewLogBuilder())

		_, err = handler.Handle()
		gomega.Expect(err).ShouldNot(gomega.BeNil())
		validationError, ok := err.(appErrors.AppError[map[string]string])
		gomega.Expect(ok).Should(gomega.Equal(true))

		errs := validationError.Data()
		gomega.Expect(errs["groupsExist"]).ShouldNot(gomega.BeEmpty())
	})

	ginkgo.It("should fail updating a readonly map variable", ginkgo.Label("map"), func() {
		projectId := testCreateProject("project")
		groups := testCreateGroups(projectId, []string{"one", "two", "three"})
		m := testCreateMap(projectId, "map")

		variables := make([]addToMap.View, 0)
		for i := 0; i < 10; i++ {
			variables = append(variables, testAddToMap(projectId, m.ID, fmt.Sprintf("name-%d", i), []shared.Reference{}, groups, "readonly"))
		}

		b, err := json.Marshal("this is metadata")
		gomega.Expect(err).Should(gomega.BeNil())

		v, err := json.Marshal("this is value")
		gomega.Expect(err).Should(gomega.BeNil())

		handler := New(NewModel(projectId, m.ShortID, variables[5].ID, []string{"metadata", "groups", "behaviour", "value"}, VariableModel{
			Name:      variables[6].ID,
			Metadata:  b,
			Groups:    []string{groups[0], groups[1]},
			Behaviour: "readonly",
			Value:     v,
		}, nil), auth.NewTestingAuthentication(false, ""), logger.NewLogBuilder())

		_, err = handler.Handle()
		gomega.Expect(err).ShouldNot(gomega.BeNil())
		validationError, ok := err.(appErrors.AppError[map[string]string])
		gomega.Expect(ok).Should(gomega.Equal(true))

		errs := validationError.Data()
		gomega.Expect(errs["behaviourReadonly"]).ShouldNot(gomega.BeEmpty())
	})

	ginkgo.It("should fail updating a name map variable if it exists", ginkgo.Label("map"), func() {
		projectId := testCreateProject("project")
		groups := testCreateGroups(projectId, []string{"one", "two", "three"})
		m := testCreateMap(projectId, "map")

		variables := make([]addToMap.View, 0)
		for i := 0; i < 10; i++ {
			variables = append(variables, testAddToMap(projectId, m.ID, fmt.Sprintf("name-%d", i), []shared.Reference{}, groups, ""))
		}

		b, err := json.Marshal("this is metadata")
		gomega.Expect(err).Should(gomega.BeNil())

		v, err := json.Marshal("this is value")
		gomega.Expect(err).Should(gomega.BeNil())

		variableId := variables[5].ID
		handler := New(NewModel(projectId, m.ID, variableId, []string{"metadata", "groups", "behaviour", "value", "name"}, VariableModel{
			Name:      "name-0",
			Metadata:  b,
			Groups:    []string{groups[0]},
			Behaviour: "modifiable",
			Value:     v,
		}, nil), auth.NewTestingAuthentication(false, ""), logger.NewLogBuilder())

		_, err = handler.Handle()
		gomega.Expect(err).ShouldNot(gomega.BeNil())
		validationError, ok := err.(appErrors.AppError[map[string]string])
		gomega.Expect(ok).Should(gomega.Equal(true))

		errs := validationError.Data()
		gomega.Expect(errs["exists"]).ShouldNot(gomega.BeEmpty())
	})
})
