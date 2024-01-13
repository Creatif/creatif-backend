package appendToList

import (
	"creatif/pkg/app/auth"
	"creatif/pkg/lib/logger"
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
				Locale:    "eng",
				Behaviour: "readonly",
				Value:     nil,
			}
		}

		handler := New(NewModel(projectId, listName, variables), auth.NewTestingAuthentication(false), logger.NewLogBuilder())

		vars, err := handler.Handle()
		testAssertErrNil(err)

		for _, v := range vars {
			testAssertIDValid(v.ID)
		}

		gomega.Expect(len(vars)).Should(gomega.Equal(5))

	})
})
