package updateMapVariable

import (
	"creatif/pkg/app/auth"
	declarations2 "creatif/pkg/app/domain/declarations"
	"creatif/pkg/app/services/lists/addToList"
	"creatif/pkg/app/services/maps/addToMap"
	"creatif/pkg/app/services/shared/connections"
	"creatif/pkg/lib/appErrors"
	"creatif/pkg/lib/storage"
	"encoding/json"
	"fmt"
	"github.com/lib/pq"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"gorm.io/datatypes"
)

var _ = ginkgo.Describe("Declaration (UPDATE) map entry tests", func() {
	ginkgo.It("should update an entry in the map by replacing it completely", ginkgo.Label("map", "map_update"), func() {
		projectId := testCreateProject("project")
		groups := testCreateGroups(projectId, []string{"one", "two", "three", "four", "five"})
		m := testCreateMap(projectId, "map")
		addToMapView := testAddToMap(projectId, m.ID, "name-0", []connections.Connection{}, groups, "")

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
		}, []connections.Connection{}, []string{}), auth.NewTestingAuthentication(false, ""))

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

		type SelectedGroups struct {
			Groups pq.StringArray `gorm:"column:groups;type:text[]"`
		}

		var selGroups SelectedGroups
		res := storage.Gorm().Raw(fmt.Sprintf("SELECT groups::text[] FROM %s WHERE variable_id = ?", (declarations2.VariableGroup{}).TableName()), addToMapView.ID).Scan(&selGroups)
		gomega.Expect(res.Error).Should(gomega.BeNil())

		gomega.Expect(len(selGroups.Groups)).Should(gomega.Equal(2))
	})

	ginkgo.It("should fail updating a map variable because of invalid number of groups", ginkgo.Label("map"), func() {
		projectId := testCreateProject("project")
		groups := testCreateGroups(projectId, []string{"one", "two", "three"})
		m := testCreateMap(projectId, "map")

		variables := make([]addToMap.View, 0)
		for i := 0; i < 10; i++ {
			variables = append(variables, testAddToMap(projectId, m.ID, fmt.Sprintf("name-%d", i), []connections.Connection{}, groups, ""))
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
		}, nil, []string{}), auth.NewTestingAuthentication(false, ""))

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
			variables = append(variables, testAddToMap(projectId, m.ID, fmt.Sprintf("name-%d", i), []connections.Connection{}, groups, "readonly"))
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
		}, nil, []string{}), auth.NewTestingAuthentication(false, ""))

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
			variables = append(variables, testAddToMap(projectId, m.ID, fmt.Sprintf("name-%d", i), []connections.Connection{}, groups, ""))
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
		}, nil, []string{}), auth.NewTestingAuthentication(false, ""))

		_, err = handler.Handle()
		gomega.Expect(err).ShouldNot(gomega.BeNil())
		validationError, ok := err.(appErrors.AppError[map[string]string])
		gomega.Expect(ok).Should(gomega.Equal(true))

		errs := validationError.Data()
		gomega.Expect(errs["exists"]).ShouldNot(gomega.BeEmpty())
	})

	ginkgo.It("should update an entry in the map with mixed connections", ginkgo.Label("map", "map_update"), func() {
		projectId := testCreateProject("project")
		groups := testCreateGroups(projectId, []string{"one", "two", "three", "four", "five"})
		m := testCreateMap(projectId, "map")
		addToMapView := testAddToMap(projectId, m.ID, "name-0", []connections.Connection{}, groups, "")

		connectionList := testCreateList(projectId, "connection list")
		connectionMap := testCreateMap(projectId, "connection map")

		createMapConnections := func() []connections.Connection {
			connsViews := make([]addToMap.View, 5)
			for i := 0; i < 5; i++ {
				connsViews[i] = testAddToMap(projectId, connectionMap.ID, fmt.Sprintf("variable-%d", i), []connections.Connection{}, []string{}, "")
			}

			conns := make([]connections.Connection, 5)
			for i, c := range connsViews {
				conns[i] = connections.Connection{
					Path:          fmt.Sprintf("conn-map-%d", i),
					StructureType: "map",
					VariableID:    c.ID,
				}
			}

			return conns
		}

		createListConnections := func() []connections.Connection {
			connsViews := make([]addToList.View, 5)
			for i := 0; i < 5; i++ {
				connsViews[i] = testAddToList(projectId, connectionList.ID, fmt.Sprintf("variable-%d", i), []connections.Connection{}, []string{})
			}

			conns := make([]connections.Connection, 5)
			for i, c := range connsViews {
				conns[i] = connections.Connection{
					Path:          fmt.Sprintf("conn-list-%d", i),
					StructureType: "list",
					VariableID:    c.ID,
				}
			}

			return conns
		}

		var conns []connections.Connection
		conns = append(conns, createListConnections()...)
		conns = append(conns, createMapConnections()...)

		b, err := json.Marshal("this is metadata")
		gomega.Expect(err).Should(gomega.BeNil())

		v, err := json.Marshal("this is value")
		gomega.Expect(err).Should(gomega.BeNil())

		handler := New(NewModel(projectId, m.ShortID, addToMapView.ID, []string{"connections"}, VariableModel{
			Name:      "new name",
			Metadata:  b,
			Groups:    []string{groups[1], groups[2]},
			Behaviour: "modifiable",
			Value:     v,
		}, conns, []string{}), auth.NewTestingAuthentication(false, ""))

		view, err := handler.Handle()

		testAssertErrNil(err)
		testAssertIDValid(view.ID)

		var connectionsCount int
		res := storage.Gorm().Raw(
			fmt.Sprintf("SELECT COUNT(*) as count FROM %s WHERE parent_variable_id = ?", (declarations2.Connection{}).TableName()),
			view.ID,
		).Scan(&connectionsCount)
		gomega.Expect(res.Error).Should(gomega.BeNil())
		gomega.Expect(connectionsCount).Should(gomega.Equal(10))
	})

	ginkgo.It("should update an entry in the map with mixed connections and remove them if connections is empty", ginkgo.Label("map", "map_update"), func() {
		projectId := testCreateProject("project")
		groups := testCreateGroups(projectId, []string{"one", "two", "three", "four", "five"})
		m := testCreateMap(projectId, "map")
		addToMapView := testAddToMap(projectId, m.ID, "name-0", []connections.Connection{}, groups, "")

		connectionList := testCreateList(projectId, "connection list")
		connectionMap := testCreateMap(projectId, "connection map")

		createMapConnections := func() []connections.Connection {
			connsViews := make([]addToMap.View, 5)
			for i := 0; i < 5; i++ {
				connsViews[i] = testAddToMap(projectId, connectionMap.ID, fmt.Sprintf("variable-%d", i), []connections.Connection{}, []string{}, "")
			}

			conns := make([]connections.Connection, 5)
			for i, c := range connsViews {
				conns[i] = connections.Connection{
					Path:          fmt.Sprintf("conn-map-%d", i),
					StructureType: "map",
					VariableID:    c.ID,
				}
			}

			return conns
		}

		createListConnections := func() []connections.Connection {
			connsViews := make([]addToList.View, 5)
			for i := 0; i < 5; i++ {
				connsViews[i] = testAddToList(projectId, connectionList.ID, fmt.Sprintf("variable-%d", i), []connections.Connection{}, []string{})
			}

			conns := make([]connections.Connection, 5)
			for i, c := range connsViews {
				conns[i] = connections.Connection{
					Path:          fmt.Sprintf("conn-list-%d", i),
					StructureType: "list",
					VariableID:    c.ID,
				}
			}

			return conns
		}

		var conns []connections.Connection
		conns = append(conns, createListConnections()...)
		conns = append(conns, createMapConnections()...)

		b, err := json.Marshal("this is metadata")
		gomega.Expect(err).Should(gomega.BeNil())

		v, err := json.Marshal("this is value")
		gomega.Expect(err).Should(gomega.BeNil())

		handler := New(NewModel(projectId, m.ShortID, addToMapView.ID, []string{"connections"}, VariableModel{
			Name:      "new name",
			Metadata:  b,
			Groups:    []string{groups[1], groups[2]},
			Behaviour: "modifiable",
			Value:     v,
		}, conns, []string{}), auth.NewTestingAuthentication(false, ""))

		view, err := handler.Handle()

		testAssertErrNil(err)
		testAssertIDValid(view.ID)

		var connectionsCount int
		res := storage.Gorm().Raw(
			fmt.Sprintf("SELECT COUNT(*) as count FROM %s WHERE parent_variable_id = ?", (declarations2.Connection{}).TableName()),
			view.ID,
		).Scan(&connectionsCount)
		gomega.Expect(res.Error).Should(gomega.BeNil())
		gomega.Expect(connectionsCount).Should(gomega.Equal(10))

		handler = New(NewModel(projectId, m.ShortID, addToMapView.ID, []string{"connections"}, VariableModel{
			Name:      "new name",
			Metadata:  b,
			Groups:    []string{groups[1], groups[2]},
			Behaviour: "modifiable",
			Value:     v,
		}, []connections.Connection{}, []string{}), auth.NewTestingAuthentication(false, ""))

		view, err = handler.Handle()

		testAssertErrNil(err)
		testAssertIDValid(view.ID)

		res = storage.Gorm().Raw(
			fmt.Sprintf("SELECT COUNT(*) as count FROM %s WHERE parent_variable_id = ?", (declarations2.Connection{}).TableName()),
			view.ID,
		).Scan(&connectionsCount)
		gomega.Expect(res.Error).Should(gomega.BeNil())
		gomega.Expect(connectionsCount).Should(gomega.Equal(0))
	})
})
