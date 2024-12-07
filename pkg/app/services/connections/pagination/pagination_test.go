package pagination

import (
	"creatif/pkg/app/auth"
	"creatif/pkg/app/services/groups/addGroups"
	"creatif/pkg/app/services/lists/addToList"
	"creatif/pkg/app/services/locales"
	"creatif/pkg/app/services/shared/connections"
	"creatif/pkg/lib/sdk"
	"fmt"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Connections pagination tests", func() {
	ginkgo.It("should paginate through list connections", func() {
		projectId := testCreateProject("project")
		groups := testCreateGroups(projectId, 0)
		list := testCreateList(projectId, "name")

		connectionList := testCreateList(projectId, "connection list")

		createListConnections := func() []connections.Connection {
			connsViews := make([]addToList.View, 100)
			for i := 0; i < 100; i++ {
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

		listConnections := createListConnections()

		singleListItem := testAddToList(projectId, list.ID, "single list item", listConnections, sdk.Map(groups, func(idx int, value addGroups.View) string {
			return value.ID
		}))

		localeId, _ := locales.GetIDWithAlpha("eng")
		handler := New(NewModel(projectId, "list", singleListItem.ID, []string{localeId}, list.ID, "created_at", "", "desc", 10, 1, []string{}, nil, "", []string{}), auth.NewTestingAuthentication(false, ""))
		views, err := handler.Handle()
		testAssertErrNil(err)

		gomega.Expect(len(views.Data)).Should(gomega.Equal(10))
		gomega.Expect(views.Total).Should(gomega.Equal(int64(100)))
	})
})
