package switchByIndex

import (
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Declaration list variable tests", func() {
	ginkgo.It("should switch two list variables indexes", func() {
		projectId := testCreateProject("project")
		indexes := testCreateListAndReturnIndexes(projectId, "list", 10)

		source := indexes[0]
		destination := indexes[5]

		handler := New(NewModel(projectId, "list", 0, 5))
		view, err := handler.Handle()
		testAssertErrNil(err)

		gomega.Expect(view.Source.Index).Should(gomega.Equal(destination))
		gomega.Expect(view.Destination.Index).Should(gomega.Equal(source))
	})
})
