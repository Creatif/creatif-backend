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
		groups := make([]string, 50)
		for i := 0; i < 50; i++ {
			groups[i] = fmt.Sprintf("group-%d", i)
		}

		l := logger.NewLogBuilder()
		handler := New(NewModel(projectId, groups), auth.NewTestingAuthentication(false, projectId), l)
		model, err := handler.Handle()
		testAssertErrNil(err)

		gomega.Expect(len(model)).Should(gomega.Equal(50))
	})

	ginkgo.It("Should remove previous groups and add new groups", func() {
		projectId := testCreateProject("project")
		groups := make([]string, 50)
		for i := 0; i < 50; i++ {
			groups[i] = fmt.Sprintf("group-%d", i)
		}

		l := logger.NewLogBuilder()
		handler := New(NewModel(projectId, groups), auth.NewTestingAuthentication(false, projectId), l)
		model, err := handler.Handle()
		testAssertErrNil(err)

		gomega.Expect(len(model)).Should(gomega.Equal(50))

		groups = make([]string, 20)
		for i := 0; i < 20; i++ {
			groups[i] = fmt.Sprintf("group-%d", i)
		}

		l = logger.NewLogBuilder()
		handler = New(NewModel(projectId, groups), auth.NewTestingAuthentication(false, projectId), l)
		model, err = handler.Handle()
		testAssertErrNil(err)

		gomega.Expect(len(model)).Should(gomega.Equal(20))
	})
})
