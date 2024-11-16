package getMapGroups

import (
	"creatif/pkg/app/auth"
	"creatif/pkg/app/services/groups/addGroups"
	"creatif/pkg/app/services/shared/connections"
	"creatif/pkg/lib/sdk"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Declaration maps variable tests", func() {
	ginkgo.It("should get all distinct groups from a map", ginkgo.Label("map"), func() {
		projectId := testCreateProject("project")
		groups := testCreateGroups(projectId, 5)
		view := testCreateMap(projectId, "map")
		variable := testAddToMap(projectId, view.ID, []connections.Connection{}, sdk.Map(groups, func(idx int, value addGroups.View) string {
			return value.ID
		}))

		handler := New(NewModel(view.Name, variable.Variable.ID, projectId), auth.NewTestingAuthentication(true, ""))
		fetchedGroups, err := handler.Handle()
		testAssertErrNil(err)

		gomega.Expect(len(fetchedGroups)).To(gomega.Equal(5))
	})
})
