package queryMapVariable

import (
	"creatif/pkg/app/auth"
	"creatif/pkg/app/services/groups/addGroups"
	"creatif/pkg/app/services/shared/connections"
	"creatif/pkg/lib/sdk"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Declaration maps variable tests", func() {
	ginkgo.It("should query a map variable by ID", ginkgo.Label("map", "query_single_variable"), func() {
		projectId := testCreateProject("project")
		groups := testCreateGroups(projectId, 5)
		mapView := testCreateMap(projectId, "name")
		variable := testAddToMap(projectId, mapView.ID, []connections.Connection{}, sdk.Map(groups, func(idx int, value addGroups.View) string {
			return value.ID
		}))

		handler := New(NewModel(projectId, mapView.ID, variable.Variable.ID), auth.NewTestingAuthentication(false, ""))
		view, err := handler.Handle()
		testAssertErrNil(err)
		testAssertIDValid(view.ID)

		gomega.Expect(view.ID).Should(gomega.Equal(variable.Variable.ID))
		gomega.Expect(view.Name).Should(gomega.Equal(variable.Variable.Name))
		gomega.Expect(len(view.Groups)).Should(gomega.Equal(5))
	})
})
