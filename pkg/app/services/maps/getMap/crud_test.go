package getMap

import (
	"creatif/pkg/app/auth"
	"github.com/onsi/ginkgo/v2"
)

var _ = ginkgo.Describe("GET map tests", func() {
	ginkgo.It("should getVariable names only (default) representation of map of variables by name", ginkgo.Label("map"), func() {
		projectId := testCreateProject("project")
		view := testCreateMap(projectId, "mapName", 10)

		handler := New(NewModel(projectId, view.ID), auth.NewTestingAuthentication(false, ""))

		mapVariablesView, err := handler.Handle()
		testAssertErrNil(err)
		testAssertIDValid(mapVariablesView.ID)
	})
})
