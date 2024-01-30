package updateList

import (
	"creatif/pkg/app/auth"
	"creatif/pkg/lib/logger"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Declaration list create tests", func() {
	ginkgo.It("should create a list", func() {
		projectId := testCreateProject("project")
		listName := testCreateList(projectId, "list", 1)

		handler := New(NewModel(projectId, []string{"name"}, listName, "newNameList"), auth.NewTestingAuthentication(false, ""), logger.NewLogBuilder())

		view, err := handler.Handle()
		testAssertErrNil(err)
		testAssertIDValid(view.ID)

		gomega.Expect(view.Name).Should(gomega.Equal("newNameList"))
	})
})
