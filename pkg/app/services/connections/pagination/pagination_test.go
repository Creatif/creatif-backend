package pagination

import (
	"creatif/pkg/app/auth"
	"creatif/pkg/app/services/groups/addGroups"
	"creatif/pkg/app/services/lists/addToList"
	"creatif/pkg/app/services/locales"
	"creatif/pkg/app/services/maps/addToMap"
	"creatif/pkg/app/services/shared/connections"
	"creatif/pkg/lib/sdk"
	"fmt"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Connections pagination tests", func() {
	ginkgo.It("should get list connections", func() {
		projectId := testCreateProject("project")
		groups := testCreateGroups(projectId, 0)
		list := testCreateList(projectId, "name")

		connectionList := testCreateList(projectId, "connection list")

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

		listConnections := createListConnections()

		singleListItem := testAddToList(projectId, list.ID, "single list item", listConnections, sdk.Map(groups, func(idx int, value addGroups.View) string {
			return value.ID
		}))

		localeId, _ := locales.GetIDWithAlpha("eng")
		handler := New(NewModel(projectId, "list", singleListItem.ID, []string{localeId}, connectionList.ID, "created_at", "", "desc", 10, 1, []string{}, nil, "", []string{}), auth.NewTestingAuthentication(false, ""))
		views, err := handler.Handle()
		testAssertErrNil(err)

		gomega.Expect(len(views)).Should(gomega.Equal(5))
	})

	ginkgo.It("should get map connections", func() {
		projectId := testCreateProject("project")
		groups := testCreateGroups(projectId, 0)
		m := testCreateMap(projectId, "name")

		connectionMap := testCreateMap(projectId, "connection map")

		createMapConnections := func() []connections.Connection {
			connsViews := make([]addToMap.View, 5)
			for i := 0; i < 5; i++ {
				connsViews[i] = testAddToMap(
					projectId,
					connectionMap.ID,
					fmt.Sprintf("variable-%d", i), []connections.Connection{},
					[]string{},
				)
			}

			conns := make([]connections.Connection, 5)
			for i, c := range connsViews {
				conns[i] = connections.Connection{
					Path:          fmt.Sprintf("conn-list-%d", i),
					StructureType: "map",
					VariableID:    c.ID,
				}
			}

			return conns
		}

		mapConnections := createMapConnections()

		singleListItem := testAddToMap(projectId, m.ID, "single list item", mapConnections, sdk.Map(groups, func(idx int, value addGroups.View) string {
			return value.ID
		}))

		localeId, _ := locales.GetIDWithAlpha("eng")
		handler := New(NewModel(projectId, "map", singleListItem.ID, []string{localeId}, connectionMap.ID, "created_at", "", "desc", 10, 1, []string{}, nil, "", []string{}), auth.NewTestingAuthentication(false, ""))
		views, err := handler.Handle()
		testAssertErrNil(err)

		gomega.Expect(len(views)).Should(gomega.Equal(5))
	})
})
