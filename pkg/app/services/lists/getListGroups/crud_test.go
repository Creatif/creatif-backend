package getListGroups

import (
	"creatif/pkg/app/auth"
	"creatif/pkg/app/services/lists/addToList"
	"creatif/pkg/app/services/shared"
	"creatif/pkg/lib/logger"
	"fmt"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Declaration list variable tests", func() {
	ginkgo.It("should get all distinct groups from a list", func() {
		projectId := testCreateProject("project")
		groups := testCreateGroups(projectId, 5)
		view := testCreateList(projectId, "list")

		listVariables := make([]addToList.View, 0)
		for i := 0; i < 10; i++ {
			listVariables = append(listVariables, testAddToList(projectId, view.ID, fmt.Sprintf("name-%d", i), []shared.Reference{}, groups))
		}

		l := logger.NewLogBuilder()
		handler := New(NewModel(view.ID, listVariables[0].ID, projectId), auth.NewTestingAuthentication(true, ""), l)
		groups, err := handler.Handle()
		testAssertErrNil(err)

		gomega.Expect(len(groups)).To(gomega.Equal(5))
	})
})
