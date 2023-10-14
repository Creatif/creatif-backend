package appendToList

import (
	"fmt"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Declaration list append tests", func() {
	ginkgo.It("should append to a list", func() {
		projectId := testCreateProject("project")
		listName := testCreateList(projectId, "list", 5)

		variables := make([]Variable, 5)
		for i := 0; i < 5; i++ {
			variables[i] = Variable{
				Name:      fmt.Sprintf("one-%d", i),
				Metadata:  nil,
				Groups:    nil,
				Behaviour: "readonly",
				Value:     nil,
			}
		}

		handler := New(NewModel(projectId, "eng", listName, variables))

		list, err := handler.Handle()
		testAssertErrNil(err)
		testAssertIDValid(list.ID)

		gomega.Expect(list.Name).Should(gomega.Equal(listName))
		gomega.Expect(list.Locale).Should(gomega.Equal("eng"))
	})
})
