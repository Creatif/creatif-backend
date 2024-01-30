package getProjectMetadata

import (
	"creatif/pkg/lib/logger"
	"fmt"
	"github.com/onsi/ginkgo/v2"
)

var _ = ginkgo.Describe("Get project metadata tests", func() {
	ginkgo.It("Should get the project metadata from various structures", func() {
		authentication := testCreateProject("project")

		for i := 0; i < 10; i++ {
			testCreateMap(authentication, fmt.Sprintf("%d-name", i), 10)
			testCreateList(authentication, fmt.Sprintf("%d-name", i), 10)

			for _, locale := range []string{"aar", "abk", "eng"} {
				testCreateDetailedVariable(authentication, locale, fmt.Sprintf("%d-name", i), "modifiable", []string{}, nil)
			}
		}

		handler := New(authentication, logger.NewLogBuilder())
		_, err := handler.Handle()
		testAssertErrNil(err)
	})
})
