package queryListByIndex

import (
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Declaration list variable tests", func() {
	ginkgo.It("should query a list variable by index 0 (zero)", func() {
		projectId := testCreateProject("project")
		listName := testCreateList(projectId, "name", 6)

		handler := New(NewModel(projectId, listName, 0))
		view, err := handler.Handle()
		testAssertErrNil(err)
		testAssertIDValid(view.ID)

		gomega.Expect(view.Name).Should(gomega.Equal("one-0"))
		gomega.Expect(view.Index).ShouldNot(gomega.BeEmpty())
	})

	ginkgo.It("should query a list variable by index 3 (zero) - middle", func() {
		projectId := testCreateProject("project")
		listName := testCreateList(projectId, "name", 6)

		handler := New(NewModel(projectId, listName, 3))
		view, err := handler.Handle()
		testAssertErrNil(err)
		testAssertIDValid(view.ID)

		gomega.Expect(view.Name).Should(gomega.Equal("one-3"))
		gomega.Expect(view.Index).ShouldNot(gomega.BeEmpty())
	})

	ginkgo.It("should query a list variable by index 5 (five) - last element", func() {
		projectId := testCreateProject("project")
		listName := testCreateList(projectId, "name", 6)

		handler := New(NewModel(projectId, listName, 5))
		view, err := handler.Handle()
		testAssertErrNil(err)
		testAssertIDValid(view.ID)

		gomega.Expect(view.Name).Should(gomega.Equal("one-5"))
		gomega.Expect(view.Index).ShouldNot(gomega.BeEmpty())
	})
})
