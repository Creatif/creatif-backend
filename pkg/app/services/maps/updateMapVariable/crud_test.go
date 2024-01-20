package updateMapVariable

import (
	"creatif/pkg/app/auth"
	"creatif/pkg/app/services/shared"
	"creatif/pkg/lib/appErrors"
	"creatif/pkg/lib/logger"
	"creatif/pkg/lib/sdk"
	"creatif/pkg/lib/storage"
	"encoding/json"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"gorm.io/datatypes"
)

var _ = ginkgo.Describe("Declaration (UPDATE) map entry tests", func() {
	ginkgo.It("should update an entry in the map by replacing it completely", func() {
		projectId := testCreateProject("project")
		m := testCreateMap(projectId, "map", 10, "modifiable")
		referenceMap := testCreateMap(projectId, "referenceMap", 10, "modifiable")
		addToMapView := testAddToMap(projectId, "map", []shared.Reference{
			{
				Name:          "myName",
				StructureName: referenceMap.Name,
				StructureType: "map",
				VariableID:    referenceMap.Variables[0].ID,
			},
			{
				Name:          "another",
				StructureName: referenceMap.Name,
				StructureType: "map",
				VariableID:    referenceMap.Variables[4].ID,
			},
			{
				Name:          "three",
				StructureName: referenceMap.Name,
				StructureType: "map",
				VariableID:    referenceMap.Variables[5].ID,
			},
		})

		b, err := json.Marshal("this is metadata")
		gomega.Expect(err).Should(gomega.BeNil())

		v, err := json.Marshal("this is value")
		gomega.Expect(err).Should(gomega.BeNil())

		handler := New(NewModel(projectId, m.Name, addToMapView.Variable.ID, []string{"metadata", "groups", "behaviour", "value", "name"}, VariableModel{
			Name:      "new name",
			Metadata:  b,
			Groups:    []string{"updated1", "updated2", "updated3"},
			Behaviour: "readonly",
			Value:     v,
		}, []shared.UpdateReference{
			{
				Name:          addToMapView.References[0].Name,
				StructureName: "referenceMap",
				StructureType: "map",
				VariableID:    referenceMap.Variables[0].ID,
			},
			{
				Name:          "new entry",
				StructureName: "referenceMap",
				StructureType: "map",
				VariableID:    referenceMap.Variables[3].ID,
			},
		}), auth.NewTestingAuthentication(false), logger.NewLogBuilder())

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
		gomega.Expect(sdk.Includes(view.Groups, "updated1")).Should(gomega.Equal(true))
		gomega.Expect(sdk.Includes(view.Groups, "updated2")).Should(gomega.Equal(true))
		gomega.Expect(sdk.Includes(view.Groups, "updated3")).Should(gomega.Equal(true))

		var count int
		res := storage.Gorm().Raw("SELECT count(id) AS count FROM declarations.references").Scan(&count)
		testAssertErrNil(res.Error)
		gomega.Expect(count).Should(gomega.Equal(2))
	})

	ginkgo.It("should fail updating a map variable because of invalid number of groups", func() {
		projectId := testCreateProject("project")
		m := testCreateMap(projectId, "map", 10, "modifiable")

		b, err := json.Marshal("this is metadata")
		gomega.Expect(err).Should(gomega.BeNil())

		v, err := json.Marshal("this is value")
		gomega.Expect(err).Should(gomega.BeNil())

		handler := New(NewModel(projectId, m.Name, m.Variables[5].ShortID, []string{"metadata", "groups", "behaviour", "value"}, VariableModel{
			Name:      "new name",
			Metadata:  b,
			Groups:    []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15", "16", "17", "18"},
			Behaviour: "readonly",
			Value:     v,
		}, nil), auth.NewTestingAuthentication(false), logger.NewLogBuilder())

		_, err = handler.Handle()
		gomega.Expect(err).ShouldNot(gomega.BeNil())
		validationError, ok := err.(appErrors.AppError[map[string]string])
		gomega.Expect(ok).Should(gomega.Equal(true))

		errs := validationError.Data()
		gomega.Expect(errs["groups"]).ShouldNot(gomega.BeEmpty())
	})

	ginkgo.It("should fail updating a readonly map variable", func() {
		projectId := testCreateProject("project")
		m := testCreateMap(projectId, "map", 10, "readonly")

		b, err := json.Marshal("this is metadata")
		gomega.Expect(err).Should(gomega.BeNil())

		v, err := json.Marshal("this is value")
		gomega.Expect(err).Should(gomega.BeNil())

		handler := New(NewModel(projectId, m.Name, m.Variables[5].ID, []string{"metadata", "groups", "behaviour", "value"}, VariableModel{
			Name:      m.Variables[6].ID,
			Metadata:  b,
			Groups:    []string{"1", "2", "3", "4", "5"},
			Behaviour: "readonly",
			Value:     v,
		}, nil), auth.NewTestingAuthentication(false), logger.NewLogBuilder())

		_, err = handler.Handle()
		gomega.Expect(err).ShouldNot(gomega.BeNil())
		validationError, ok := err.(appErrors.AppError[map[string]string])
		gomega.Expect(ok).Should(gomega.Equal(true))

		errs := validationError.Data()
		gomega.Expect(errs["behaviourReadonly"]).ShouldNot(gomega.BeEmpty())
	})

	ginkgo.It("should fail updating a name map variable if it exists", func() {
		projectId := testCreateProject("project")
		m := testCreateMap(projectId, "map", 10, "modifiable")

		b, err := json.Marshal("this is metadata")
		gomega.Expect(err).Should(gomega.BeNil())

		v, err := json.Marshal("this is value")
		gomega.Expect(err).Should(gomega.BeNil())

		handler := New(NewModel(projectId, m.Name, m.Variables[5].ID, []string{"metadata", "groups", "behaviour", "value", "name"}, VariableModel{
			Name:      "name-0",
			Metadata:  b,
			Groups:    []string{"1", "2", "3", "4", "5"},
			Behaviour: "modifiable",
			Value:     v,
		}, nil), auth.NewTestingAuthentication(false), logger.NewLogBuilder())

		_, err = handler.Handle()
		gomega.Expect(err).ShouldNot(gomega.BeNil())
		validationError, ok := err.(appErrors.AppError[map[string]string])
		gomega.Expect(ok).Should(gomega.Equal(true))

		errs := validationError.Data()
		gomega.Expect(errs["exists"]).ShouldNot(gomega.BeEmpty())
	})
})
