package addToMap

import (
	"creatif/pkg/app/auth"
	"creatif/pkg/app/services/groups/addGroups"
	getMap2 "creatif/pkg/app/services/maps/getMap"
	"creatif/pkg/app/services/shared"
	"creatif/pkg/lib/sdk"
	"creatif/pkg/lib/storage"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Declaration (UPDATE) map entry tests", func() {
	ginkgo.It("should add an entry to the map by name with references", ginkgo.Label("map"), func() {
		projectId := testCreateProject("project")
		groups := testCreateGroups(projectId, 5)
		m := testCreateMap(projectId, "mapName", 10, sdk.Map(groups, func(idx int, value addGroups.View) string {
			return value.ID
		}))
		reference := testCreateMap(projectId, "referenceMap", 10, sdk.Map(groups, func(idx int, value addGroups.View) string {
			return value.ID
		}))

		handler := New(NewModel(projectId, m.ID, VariableModel{
			Name:     "newEntry",
			Metadata: nil,
			Groups: sdk.Map(groups, func(idx int, value addGroups.View) string {
				return value.ID
			}),
			Locale:    "eng",
			Behaviour: "readonly",
			Value:     nil,
		}, []shared.Reference{
			{
				Name:          "first",
				StructureName: reference.Name,
				StructureType: "map",
				VariableID:    reference.Variables[0].ID,
			},
			{
				Name:          "second",
				StructureName: reference.Name,
				StructureType: "map",
				VariableID:    reference.Variables[1].ID,
			},
			{
				Name:          "third",
				StructureName: reference.Name,
				StructureType: "map",
				VariableID:    reference.Variables[2].ID,
			},
		}, []string{}), auth.NewTestingAuthentication(false, ""))

		_, err := handler.Handle()
		gomega.Expect(err).Should(gomega.BeNil())

		getMapHandler := getMap2.New(getMap2.NewModel(projectId, m.Name), auth.NewTestingAuthentication(false, ""))
		maps, err := getMapHandler.Handle()
		gomega.Expect(err).Should(gomega.BeNil())
		testAssertIDValid(maps.ID)
		gomega.Expect(maps.ProjectID).Should(gomega.Equal(projectId))

		var count int
		res := storage.Gorm().Raw("SELECT count(id) AS count FROM declarations.references").Scan(&count)
		gomega.Expect(res.Error).Should(gomega.BeNil())
		gomega.Expect(count).Should(gomega.Equal(3))
	})

	ginkgo.It("should fail to add an entry because of a duplicate", ginkgo.Label("map"), func() {
		projectId := testCreateProject("project")
		m := testCreateMap(projectId, "mapName", 10, nil)
		reference := testCreateMap(projectId, "referenceMap", 10, nil)

		handler := New(NewModel(projectId, m.Name, VariableModel{
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
				VariableID:    reference.Variables[0].ID,
			},
			{
				Name:          "second",
				StructureName: reference.Name,
				StructureType: "map",
				VariableID:    reference.Variables[0].ID,
			},
			{
				Name:          "first",
				StructureName: reference.Name,
				StructureType: "map",
				VariableID:    reference.Variables[2].ID,
			},
		}, []string{}), auth.NewTestingAuthentication(false, ""))

		_, err := handler.Handle()
		gomega.Expect(err).ShouldNot(gomega.BeNil())
	})

	ginkgo.It("should add an entry to the map by id", ginkgo.Label("map"), func() {
		projectId := testCreateProject("project")
		m := testCreateMap(projectId, "mapName", 10, nil)

		handler := New(NewModel(projectId, m.ID, VariableModel{
			Name:      "newEntry",
			Metadata:  nil,
			Groups:    nil,
			Locale:    "eng",
			Behaviour: "readonly",
			Value:     nil,
		}, []shared.Reference{}, []string{}), auth.NewTestingAuthentication(false, ""))

		_, err := handler.Handle()
		gomega.Expect(err).Should(gomega.BeNil())

		getMapHandler := getMap2.New(getMap2.NewModel(projectId, m.Name), auth.NewTestingAuthentication(false, ""))
		maps, err := getMapHandler.Handle()
		gomega.Expect(err).Should(gomega.BeNil())
		testAssertIDValid(maps.ID)
		gomega.Expect(maps.ProjectID).Should(gomega.Equal(projectId))
	})

	ginkgo.It("should add an entry to the map by shortID", ginkgo.Label("map"), func() {
		projectId := testCreateProject("project")
		m := testCreateMap(projectId, "mapName", 10, nil)

		handler := New(NewModel(projectId, m.ID, VariableModel{
			Name:      "newEntry",
			Metadata:  nil,
			Groups:    nil,
			Locale:    "eng",
			Behaviour: "readonly",
			Value:     nil,
		}, []shared.Reference{}, []string{}), auth.NewTestingAuthentication(false, ""))

		_, err := handler.Handle()
		gomega.Expect(err).Should(gomega.BeNil())

		getMapHandler := getMap2.New(getMap2.NewModel(projectId, m.ShortID), auth.NewTestingAuthentication(false, ""))
		maps, err := getMapHandler.Handle()
		gomega.Expect(err).Should(gomega.BeNil())
		testAssertIDValid(maps.ID)
		gomega.Expect(maps.ProjectID).Should(gomega.Equal(projectId))
	})
})
