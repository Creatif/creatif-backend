package queryListByID

import (
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Declaration list variable tests", func() {
	ginkgo.It("should query a list variable by ID", func() {
		projectId := testCreateProject("project")
		variableIds := testCreateListAndReturnIds(projectId, "name", 6)

		handler := New(NewModel(projectId, "name", variableIds[3]))
		view, err := handler.Handle()
		testAssertErrNil(err)
		testAssertIDValid(view.ID)

		gomega.Expect(view.ID).Should(gomega.Equal(variableIds[3]))
		gomega.Expect(view.Name).Should(gomega.Equal("one-3"))
		gomega.Expect(view.Index).ShouldNot(gomega.BeEmpty())
	})
})
