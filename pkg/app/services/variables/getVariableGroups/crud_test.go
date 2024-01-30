package getVariableGroups

import (
	"creatif/pkg/app/auth"
	"creatif/pkg/lib/logger"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Declaration list variable tests", func() {
	ginkgo.It("should return empty groups list when there are no groups", func() {
		projectId := testCreateProject("project")
		variable := testCreateDeclarationVariable(projectId, "variable", "modifiable")

		l := logger.NewLogBuilder()
		handler := New(NewModel(variable.Name, projectId), auth.NewTestingAuthentication(true, ""), l)
		groups, err := handler.Handle()
		testAssertErrNil(err)

		gomega.Expect(len(groups)).To(gomega.Equal(0))
	})

	ginkgo.It("should return empty groups list when there are no groups", func() {
		projectId := testCreateProject("project")
		variable := testCreateDetailedVariable(projectId, "variable", "modifiable", []string{"one", "two", "three"}, nil)

		l := logger.NewLogBuilder()
		handler := New(NewModel(variable.Name, projectId), auth.NewTestingAuthentication(true, ""), l)
		groups, err := handler.Handle()
		testAssertErrNil(err)

		gomega.Expect(len(groups)).To(gomega.Equal(3))
	})
})
