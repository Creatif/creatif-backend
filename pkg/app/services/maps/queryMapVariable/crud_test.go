package queryMapVariable

import (
	"creatif/pkg/app/auth"
	"creatif/pkg/app/services/groups/addGroups"
	"creatif/pkg/app/services/lists/addToList"
	"creatif/pkg/app/services/maps/addToMap"
	"creatif/pkg/app/services/shared/connections"
	"creatif/pkg/lib/sdk"
	"fmt"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Declaration maps variable tests", func() {
	ginkgo.It("should query a map variable by ID", ginkgo.Label("map", "query_single_variable"), func() {
		projectId := testCreateProject("project")
		groups := testCreateGroups(projectId, 5)
		mapView := testCreateMap(projectId, "name")
		variable := testAddToMap(projectId, mapView.ID, "my variable", []connections.Connection{}, sdk.Map(groups, func(idx int, value addGroups.View) string {
			return value.ID
		}))

		handler := New(NewModel(projectId, mapView.ID, variable.ID), auth.NewTestingAuthentication(false, ""))
		view, err := handler.Handle()
		testAssertErrNil(err)
		testAssertIDValid(view.ID)

		gomega.Expect(view.ID).Should(gomega.Equal(variable.ID))
		gomega.Expect(view.Name).Should(gomega.Equal(variable.Name))
		gomega.Expect(len(view.Groups)).Should(gomega.Equal(5))
		gomega.Expect(len(view.Connections)).Should(gomega.Equal(0))
	})

	ginkgo.It("should query a map variable by ID with map connections", ginkgo.Label("map", "query_single_variable"), func() {
		projectId := testCreateProject("project")
		groups := testCreateGroups(projectId, 5)
		mapView := testCreateMap(projectId, "name")

		connectionMap := testCreateMap(projectId, "connection map")
		createMapConnections := func() ([]connections.Connection, []addToMap.View) {
			conns := make([]connections.Connection, 5)
			variables := make([]addToMap.View, 5)
			for i := 0; i < 5; i++ {
				m := testAddToMap(projectId, connectionMap.ID, fmt.Sprintf("map-%d", i), []connections.Connection{}, sdk.Map(groups, func(idx int, value addGroups.View) string {
					return value.ID
				}))

				variables[i] = m
				conns[i] = connections.Connection{
					Path:          fmt.Sprintf("conn-map-%d", i),
					StructureType: "map",
					VariableID:    m.ID,
				}
			}

			return conns, variables
		}

		conns, connectionVariables := createMapConnections()
		variable := testAddToMap(projectId, mapView.ID, "my variable", conns, sdk.Map(groups, func(idx int, value addGroups.View) string {
			return value.ID
		}))

		handler := New(NewModel(projectId, mapView.ID, variable.ID), auth.NewTestingAuthentication(false, ""))
		view, err := handler.Handle()
		testAssertErrNil(err)
		testAssertIDValid(view.ID)

		gomega.Expect(view.ID).Should(gomega.Equal(variable.ID))
		gomega.Expect(view.Name).Should(gomega.Equal(variable.Name))
		gomega.Expect(len(view.Groups)).Should(gomega.Equal(5))
		gomega.Expect(len(view.Connections)).Should(gomega.Equal(5))

		for _, c := range view.Connections {
			gomega.Expect(c.ParentVariableID).Should(gomega.Equal(view.ID))
			gomega.Expect(c.ParentStructureType).Should(gomega.Equal("map"))
			gomega.Expect(c.ChildStructureType).Should(gomega.Equal("map"))

			found := false
			for _, v := range connectionVariables {
				if v.ID == c.ChildVariableID {
					found = true
					break
				}
			}

			gomega.Expect(found).Should(gomega.BeTrue())
		}
	})

	ginkgo.It("should query a map variable by ID with list connections", ginkgo.Label("map", "query_single_variable"), func() {
		projectId := testCreateProject("project")
		groups := testCreateGroups(projectId, 5)
		mapView := testCreateMap(projectId, "name")

		connectionList := testCreateList(projectId, "connection list")
		createListConnections := func() ([]connections.Connection, []addToList.View) {
			conns := make([]connections.Connection, 5)
			variables := make([]addToList.View, 5)
			for i := 0; i < 5; i++ {
				m := testAddToList(projectId, connectionList.ID, fmt.Sprintf("list-%d", i), []connections.Connection{}, sdk.Map(groups, func(idx int, value addGroups.View) string {
					return value.ID
				}))

				variables[i] = m
				conns[i] = connections.Connection{
					Path:          fmt.Sprintf("conn-list-%d", i),
					StructureType: "list",
					VariableID:    m.ID,
				}
			}

			return conns, variables
		}

		conns, connectionVariables := createListConnections()
		variable := testAddToMap(projectId, mapView.ID, "my variable", conns, sdk.Map(groups, func(idx int, value addGroups.View) string {
			return value.ID
		}))

		handler := New(NewModel(projectId, mapView.ID, variable.ID), auth.NewTestingAuthentication(false, ""))
		view, err := handler.Handle()
		testAssertErrNil(err)
		testAssertIDValid(view.ID)

		gomega.Expect(view.ID).Should(gomega.Equal(variable.ID))
		gomega.Expect(view.Name).Should(gomega.Equal(variable.Name))
		gomega.Expect(len(view.Groups)).Should(gomega.Equal(5))
		gomega.Expect(len(view.Connections)).Should(gomega.Equal(5))

		for _, c := range view.Connections {
			gomega.Expect(c.ParentVariableID).Should(gomega.Equal(view.ID))
			gomega.Expect(c.ParentStructureType).Should(gomega.Equal("map"))
			gomega.Expect(c.ChildStructureType).Should(gomega.Equal("list"))

			found := false
			for _, v := range connectionVariables {
				if v.ID == c.ChildVariableID {
					found = true
					break
				}
			}

			gomega.Expect(found).Should(gomega.BeTrue())
		}
	})

	ginkgo.It("should query a map variable by ID with mixed connections", ginkgo.Label("map", "query_single_variable"), func() {
		projectId := testCreateProject("project")
		groups := testCreateGroups(projectId, 5)
		mapView := testCreateMap(projectId, "name")

		connectionMap := testCreateMap(projectId, "connection map")
		createMapConnections := func() ([]connections.Connection, []addToMap.View) {
			conns := make([]connections.Connection, 5)
			variables := make([]addToMap.View, 5)
			for i := 0; i < 5; i++ {
				m := testAddToMap(projectId, connectionMap.ID, fmt.Sprintf("map-%d", i), []connections.Connection{}, sdk.Map(groups, func(idx int, value addGroups.View) string {
					return value.ID
				}))

				variables[i] = m
				conns[i] = connections.Connection{
					Path:          fmt.Sprintf("conn-map-%d", i),
					StructureType: "map",
					VariableID:    m.ID,
				}
			}

			return conns, variables
		}

		mapConns, _ := createMapConnections()
		connectionList := testCreateList(projectId, "connection list")
		createListConnections := func() ([]connections.Connection, []addToList.View) {
			conns := make([]connections.Connection, 5)
			variables := make([]addToList.View, 5)
			for i := 0; i < 5; i++ {
				m := testAddToList(projectId, connectionList.ID, fmt.Sprintf("list-%d", i), []connections.Connection{}, sdk.Map(groups, func(idx int, value addGroups.View) string {
					return value.ID
				}))

				variables[i] = m
				conns[i] = connections.Connection{
					Path:          fmt.Sprintf("conn-list-%d", i),
					StructureType: "list",
					VariableID:    m.ID,
				}
			}

			return conns, variables
		}

		listConns, _ := createListConnections()
		conns := make([]connections.Connection, 0)
		conns = append(conns, mapConns...)
		conns = append(conns, listConns...)

		variable := testAddToMap(projectId, mapView.ID, "my variable", conns, sdk.Map(groups, func(idx int, value addGroups.View) string {
			return value.ID
		}))

		handler := New(NewModel(projectId, mapView.ID, variable.ID), auth.NewTestingAuthentication(false, ""))
		view, err := handler.Handle()
		testAssertErrNil(err)
		testAssertIDValid(view.ID)

		gomega.Expect(view.ID).Should(gomega.Equal(variable.ID))
		gomega.Expect(view.Name).Should(gomega.Equal(variable.Name))
		gomega.Expect(len(view.Groups)).Should(gomega.Equal(5))
		gomega.Expect(len(view.Connections)).Should(gomega.Equal(10))

		listConnections := sdk.Filter(view.Connections, func(idx int, value ConnectionView) bool {
			return value.ChildStructureType == "list"
		})
		gomega.Expect(len(listConnections)).Should(gomega.Equal(5))

		mapConnections := sdk.Filter(view.Connections, func(idx int, value ConnectionView) bool {
			return value.ChildStructureType == "map"
		})
		gomega.Expect(len(mapConnections)).Should(gomega.Equal(5))
	})
})
