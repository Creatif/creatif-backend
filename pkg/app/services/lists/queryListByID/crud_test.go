package queryListByID

import (
	"creatif/pkg/app/auth"
	"creatif/pkg/app/services/lists/addToList"
	"creatif/pkg/app/services/maps/addToMap"
	"creatif/pkg/app/services/shared/connections"
	"creatif/pkg/lib/sdk"
	"fmt"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Declaration list variable tests", func() {
	ginkgo.It("should query a list variable by ID", func() {
		projectId := testCreateProject("project")
		list := testCreateList(projectId, "name")

		variables := make([]addToList.View, 5)
		for i := 0; i < 5; i++ {
			variables[i] = testAddToList(projectId, fmt.Sprintf("variable-%d", i), list.ID, []connections.Connection{}, []string{})
		}

		selectedVariableId := variables[4].ID

		handler := New(NewModel(projectId, list.ID, selectedVariableId, "connection"), auth.NewTestingAuthentication(false, ""))
		view, err := handler.Handle()
		gomega.Expect(err).Should(gomega.BeNil())
		testAssertIDValid(view.ID)

		gomega.Expect(view.ID).Should(gomega.Equal(selectedVariableId))
		gomega.Expect(view.Name).Should(gomega.Equal(variables[4].Name))
		gomega.Expect(view.Locale).Should(gomega.Equal("eng"))
		gomega.Expect(len(view.Connections)).Should(gomega.Equal(0))
	})

	ginkgo.It("should query a list variable by ID with list connections", func() {
		projectId := testCreateProject("project")
		list := testCreateList(projectId, "name")
		connectionList := testCreateList(projectId, "connection list")

		connsViews := make([]addToList.View, 5)
		for i := 0; i < 5; i++ {
			connsViews[i] = testAddToList(projectId, fmt.Sprintf("variable-%d", i), connectionList.ID, []connections.Connection{}, []string{})
		}

		conns := make([]connections.Connection, 5)
		for i, c := range connsViews {
			conns[i] = connections.Connection{
				Path:          fmt.Sprintf("conn-%d", i),
				StructureType: "list",
				VariableID:    c.ID,
			}
		}

		variables := make([]addToList.View, 5)
		for i := 0; i < 5; i++ {
			variables[i] = testAddToList(projectId, fmt.Sprintf("variable-%d", i), list.ID, conns, []string{})
		}

		selectedVariableId := variables[4].ID

		handler := New(NewModel(projectId, list.ID, selectedVariableId, "connection"), auth.NewTestingAuthentication(false, ""))
		view, err := handler.Handle()
		gomega.Expect(err).Should(gomega.BeNil())
		testAssertIDValid(view.ID)

		gomega.Expect(view.ID).Should(gomega.Equal(selectedVariableId))
		gomega.Expect(view.Name).Should(gomega.Equal(variables[4].Name))
		gomega.Expect(view.Locale).Should(gomega.Equal("eng"))
		gomega.Expect(len(view.Connections)).Should(gomega.Equal(5))
		gomega.Expect(len(view.ChildStructures)).Should(gomega.Equal(1))

		gomega.Expect(view.ChildStructures[0].StructureID).ShouldNot(gomega.BeEmpty())
		gomega.Expect(view.ChildStructures[0].StructureType).Should(gomega.Equal("list"))
		gomega.Expect(view.ChildStructures[0].StructureName).Should(gomega.Equal("connection list"))

		for _, c := range view.Connections {
			gomega.Expect(c.ParentVariableID).Should(gomega.Equal(selectedVariableId))
			gomega.Expect(c.ParentStructureType).Should(gomega.Equal("list"))
			gomega.Expect(c.ChildStructureType).Should(gomega.Equal("list"))

			found := false
			for _, v := range connsViews {
				if v.ID == c.ChildVariableID {
					found = true
					break
				}
			}

			gomega.Expect(found).Should(gomega.BeTrue())
		}
	})

	ginkgo.It("should query a list variable by ID with map connections", func() {
		projectId := testCreateProject("project")
		list := testCreateList(projectId, "name")
		connectionList := testCreateMap(projectId, "connection map")

		connsViews := make([]addToMap.View, 5)
		for i := 0; i < 5; i++ {
			connsViews[i] = testAddToMap(projectId, fmt.Sprintf("variable-%d", i), connectionList.ID, []connections.Connection{}, []string{})
		}

		conns := make([]connections.Connection, 5)
		for i, c := range connsViews {
			conns[i] = connections.Connection{
				Path:          fmt.Sprintf("conn-%d", i),
				StructureType: "map",
				VariableID:    c.ID,
			}
		}

		variables := make([]addToList.View, 5)
		for i := 0; i < 5; i++ {
			variables[i] = testAddToList(projectId, fmt.Sprintf("variable-%d", i), list.ID, conns, []string{})
		}

		selectedVariableId := variables[4].ID

		handler := New(NewModel(projectId, list.ID, selectedVariableId, "connection"), auth.NewTestingAuthentication(false, ""))
		view, err := handler.Handle()
		gomega.Expect(err).Should(gomega.BeNil())
		testAssertIDValid(view.ID)

		gomega.Expect(view.ID).Should(gomega.Equal(selectedVariableId))
		gomega.Expect(view.Name).Should(gomega.Equal(variables[4].Name))
		gomega.Expect(view.Locale).Should(gomega.Equal("eng"))
		gomega.Expect(len(view.Connections)).Should(gomega.Equal(5))
		gomega.Expect(len(view.ChildStructures)).Should(gomega.Equal(1))

		gomega.Expect(view.ChildStructures[0].StructureID).ShouldNot(gomega.BeEmpty())
		gomega.Expect(view.ChildStructures[0].StructureType).Should(gomega.Equal("map"))
		gomega.Expect(view.ChildStructures[0].StructureName).Should(gomega.Equal("connection map"))

		for _, c := range view.Connections {
			gomega.Expect(c.ParentVariableID).Should(gomega.Equal(selectedVariableId))
			gomega.Expect(c.ParentStructureType).Should(gomega.Equal("list"))
			gomega.Expect(c.ChildStructureType).Should(gomega.Equal("map"))

			found := false
			for _, v := range connsViews {
				if v.ID == c.ChildVariableID {
					found = true
					break
				}
			}

			gomega.Expect(found).Should(gomega.BeTrue())
		}
	})

	ginkgo.It("should query a list variable by ID with mixed connections", func() {
		projectId := testCreateProject("project")
		list := testCreateList(projectId, "name")
		connectionList := testCreateList(projectId, "connection list")
		connectionMap := testCreateMap(projectId, "connection map")

		createMapConnections := func() []connections.Connection {
			connsViews := make([]addToMap.View, 5)
			for i := 0; i < 5; i++ {
				connsViews[i] = testAddToMap(projectId, fmt.Sprintf("variable-%d", i), connectionMap.ID, []connections.Connection{}, []string{})
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
				connsViews[i] = testAddToList(projectId, fmt.Sprintf("variable-%d", i), connectionList.ID, []connections.Connection{}, []string{})
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

		variables := make([]addToList.View, 5)
		for i := 0; i < 5; i++ {
			variables[i] = testAddToList(projectId, fmt.Sprintf("variable-%d", i), list.ID, conns, []string{})
		}

		selectedVariableId := variables[4].ID

		handler := New(NewModel(projectId, list.ID, selectedVariableId, "connection"), auth.NewTestingAuthentication(false, ""))
		view, err := handler.Handle()
		gomega.Expect(err).Should(gomega.BeNil())
		testAssertIDValid(view.ID)

		gomega.Expect(view.ID).Should(gomega.Equal(selectedVariableId))
		gomega.Expect(view.Name).Should(gomega.Equal(variables[4].Name))
		gomega.Expect(view.Locale).Should(gomega.Equal("eng"))
		gomega.Expect(len(view.Connections)).Should(gomega.Equal(10))
		gomega.Expect(len(view.ChildStructures)).Should(gomega.Equal(2))

		listConnections := sdk.Filter(view.Connections, func(idx int, value ConnectionView) bool {
			return value.ChildStructureType == "list"
		})
		gomega.Expect(len(listConnections)).Should(gomega.Equal(5))

		mapConnections := sdk.Filter(view.Connections, func(idx int, value ConnectionView) bool {
			return value.ChildStructureType == "map"
		})
		gomega.Expect(len(mapConnections)).Should(gomega.Equal(5))
	})

	ginkgo.It("should query a list variable by ID with mixed connections and with each connection variable view", func() {
		projectId := testCreateProject("project")
		list := testCreateList(projectId, "name")
		connectionList := testCreateList(projectId, "connection list")
		connectionMap := testCreateMap(projectId, "connection map")

		createMapConnections := func() []connections.Connection {
			connsViews := make([]addToMap.View, 5)
			for i := 0; i < 5; i++ {
				connsViews[i] = testAddToMap(projectId, fmt.Sprintf("variable-%d", i), connectionMap.ID, []connections.Connection{}, []string{})
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
				connsViews[i] = testAddToList(projectId, fmt.Sprintf("variable-%d", i), connectionList.ID, []connections.Connection{}, []string{})
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

		variables := make([]addToList.View, 5)
		for i := 0; i < 5; i++ {
			variables[i] = testAddToList(projectId, fmt.Sprintf("variable-%d", i), list.ID, conns, []string{})
		}

		selectedVariableId := variables[4].ID

		handler := New(NewModel(projectId, list.ID, selectedVariableId, "connection"), auth.NewTestingAuthentication(false, ""))
		view, err := handler.Handle()
		gomega.Expect(err).Should(gomega.BeNil())
		testAssertIDValid(view.ID)

		handler = New(NewModel(projectId, list.ID, selectedVariableId, "value"), auth.NewTestingAuthentication(false, ""))
		view, err = handler.Handle()
		gomega.Expect(err).Should(gomega.BeNil())
		testAssertIDValid(view.ID)
	})
})
