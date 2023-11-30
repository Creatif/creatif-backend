package getListGroups

import (
	"creatif/pkg/app/auth"
	"creatif/pkg/lib/logger"
	"fmt"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Declaration list variable tests", func() {
	ginkgo.It("should get all distinct groups from a list", func() {
		projectId := testCreateProject("project")
		listName := testCreateList(projectId, "list", 5)

		l := logger.NewLogBuilder()
		handler := New(NewModel("", listName, ""), auth.NewTestingAuthentication(true), l)
		groups, err := handler.Handle()
		testAssertErrNil(err)

		fmt.Println(groups)
		gomega.Expect(len(groups)).To(gomega.Equal(3))
	})
})