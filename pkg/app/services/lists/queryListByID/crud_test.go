package queryListByID

import (
	"creatif/pkg/app/auth"
	"creatif/pkg/lib/logger"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Declaration list variable tests", func() {
	ginkgo.It("should query a list variable by ID", func() {
		projectId := testCreateProject("project")
		listId, variableIds := testCreateListAndReturnIds(projectId, "name", 6)

		selectedVariable := variableIds[3]

		handler := New(NewModel(projectId, listId, selectedVariable["id"]), auth.NewTestingAuthentication(false, ""), logger.NewLogBuilder())
		view, err := handler.Handle()
		gomega.Expect(err).Should(gomega.BeNil())
		testAssertIDValid(view.ID)

		gomega.Expect(view.ID).Should(gomega.Equal(variableIds[3]["id"]))
		gomega.Expect(view.Name).Should(gomega.Equal(selectedVariable["name"]))
		gomega.Expect(view.Locale).Should(gomega.Equal("eng"))
	})
})
