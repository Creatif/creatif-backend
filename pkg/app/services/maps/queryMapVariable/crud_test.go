package queryMapVariable

import (
	"creatif/pkg/app/auth"
	"creatif/pkg/lib/logger"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Declaration maps variable tests", func() {
	ginkgo.It("should query a map variable by ID", func() {
		projectId := testCreateProject("project")
		mapView := testCreateMap(projectId, "name", 6)

		selectedVariable := mapView.Variables[3]

		handler := New(NewModel(projectId, mapView.ID, selectedVariable.ID), auth.NewTestingAuthentication(false, ""), logger.NewLogBuilder())
		view, err := handler.Handle()
		testAssertErrNil(err)
		testAssertIDValid(view.ID)

		gomega.Expect(view.ID).Should(gomega.Equal(selectedVariable.ID))
		gomega.Expect(view.Name).Should(gomega.Equal(selectedVariable.Name))
	})
})
