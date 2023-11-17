package createList

import (
	"creatif/pkg/app/auth"
	"creatif/pkg/lib/logger"
	"fmt"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Declaration list create tests", func() {
	ginkgo.It("should create a list", func() {
		projectId := testCreateProject("project")
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

		handler := New(NewModel(projectId, "eng", "list", variables), auth.NewTestingAuthentication(false), logger.NewLogBuilder())

		list, err := handler.Handle()
		testAssertErrNil(err)
		testAssertIDValid(list.ID)

		gomega.Expect(list.Name).Should(gomega.Equal("list"))
		gomega.Expect(list.Locale).Should(gomega.Equal("eng"))
	})
})
