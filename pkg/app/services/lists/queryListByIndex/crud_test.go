package queryListByIndex

import (
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Declaration list variable tests", func() {
	ginkgo.It("should query a list variable by index 0 (zero)", func() {
		ginkgo.Skip("Querying by index should be revisited later. Currently impossible!")
		projectId := testCreateProject("project")
		listName := testCreateList(projectId, "name", 6)

		handler := New(NewModel(projectId, "eng", listName, 0))
		view, err := handler.Handle()
		testAssertErrNil(err)
		testAssertIDValid(view.ID)

		gomega.Expect(view.Name).Should(gomega.Equal("one-0"))
		gomega.Expect(view.Locale).Should(gomega.Equal("eng"))
	})

	ginkgo.It("should query a list variable by index 3 (zero) - middle", func() {
		ginkgo.Skip("Querying by index should be revisited later. Currently impossible!")

		projectId := testCreateProject("project")
		listName := testCreateList(projectId, "name", 6)

		handler := New(NewModel(projectId, "eng", listName, 3))
		view, err := handler.Handle()
		testAssertErrNil(err)
		testAssertIDValid(view.ID)

		gomega.Expect(view.Name).Should(gomega.Equal("one-3"))
		gomega.Expect(view.Locale).Should(gomega.Equal("eng"))
	})

	ginkgo.It("should query a list variable by index 5 (five) - last element", func() {
		ginkgo.Skip("Querying by index should be revisited later. Currently impossible!")
		projectId := testCreateProject("project")
		listName := testCreateList(projectId, "name", 6)

		handler := New(NewModel(projectId, "eng", listName, 5))
		view, err := handler.Handle()
		testAssertErrNil(err)
		testAssertIDValid(view.ID)

		gomega.Expect(view.Name).Should(gomega.Equal("one-5"))
		gomega.Expect(view.Locale).Should(gomega.Equal("eng"))
	})
})
