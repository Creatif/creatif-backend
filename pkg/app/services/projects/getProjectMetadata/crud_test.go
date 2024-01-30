package getProjectMetadata

import (
	"creatif/pkg/app/auth"
	"creatif/pkg/lib/logger"
	"fmt"
	"github.com/onsi/ginkgo/v2"
)

var _ = ginkgo.Describe("Get project metadata tests", func() {
	ginkgo.It("Should get the project metadata from various structures", func() {
		projectID := testCreateProject("project")

		for i := 0; i < 10; i++ {
			testCreateMap(projectID, fmt.Sprintf("%d-name", i), 10)
			testCreateList(projectID, fmt.Sprintf("%d-name", i), 10)

			for _, locale := range []string{"aar", "abk", "eng"} {
				testCreateDetailedVariable(projectID, locale, fmt.Sprintf("%d-name", i), "modifiable", []string{}, nil)
			}
		}

		a := auth.NewTestingAuthentication(true, projectID)
		handler := New(a, logger.NewLogBuilder())
		model, err := handler.Handle()
		testAssertErrNil(err)

		fmt.Println(model)
	})
})
