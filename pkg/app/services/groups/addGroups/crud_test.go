package addGroups

import (
	"creatif/pkg/app/auth"
	"creatif/pkg/lib/logger"
	"fmt"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Adding new groups", func() {
	ginkgo.It("Should add new groups", func() {
		projectId := testCreateProject("project")
		groups := make([]GroupModel, 50)
		for i := 0; i < 50; i++ {
			groups[i] = GroupModel{
				ID:     "",
				Name:   fmt.Sprintf("group-%d", i),
				Type:   "new",
				Action: "create",
			}
		}

		l := logger.NewLogBuilder()
		handler := New(NewModel(projectId, groups), auth.NewTestingAuthentication(false, projectId), l)
		model, err := handler.Handle()
		testAssertErrNil(err)

		gomega.Expect(len(model)).Should(gomega.Equal(50))
	})

	ginkgo.It("Should remove previous groups and add new groups", func() {
		projectId := testCreateProject("project")
		groups := make([]GroupModel, 50)
		for i := 0; i < 50; i++ {
			groups[i] = GroupModel{
				ID:     "",
				Name:   fmt.Sprintf("group-%d", i),
				Type:   "new",
				Action: "create",
			}
		}

		l := logger.NewLogBuilder()
		handler := New(NewModel(projectId, groups), auth.NewTestingAuthentication(false, projectId), l)
		model, err := handler.Handle()
		testAssertErrNil(err)

		gomega.Expect(len(model)).Should(gomega.Equal(50))

		updatedGroups := make([]GroupModel, 20)
		for i := 0; i < 50; i++ {
			if i < 20 {
				groups[i] = GroupModel{
					ID:     groups[i].ID,
					Name:   fmt.Sprintf("group-%d", i),
					Type:   "current",
					Action: "remove",
				}
			}
		}

		l = logger.NewLogBuilder()
		handler = New(NewModel(projectId, updatedGroups), auth.NewTestingAuthentication(false, projectId), l)
		model, err = handler.Handle()
		testAssertErrNil(err)

		gomega.Expect(len(model)).Should(gomega.Equal(50))
	})
})
