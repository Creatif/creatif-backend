package getGroups

import (
	"creatif/pkg/app/auth"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Getting groups", func() {
	ginkgo.It("should return a list of groups", func() {
		projectId := testCreateProject("project")
		testCreateGroups(projectId, 43)

		handler := New(NewModel(projectId), auth.NewTestingAuthentication(false, projectId))
		model, err := handler.Handle()
		testAssertErrNil(err)

		gomega.Expect(len(model)).Should(gomega.Equal(43))
	})
})
