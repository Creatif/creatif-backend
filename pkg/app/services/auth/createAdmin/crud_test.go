package createAdmin

import (
	"creatif/pkg/lib/logger"
	"fmt"
	"github.com/google/uuid"
	"github.com/onsi/ginkgo/v2"
)

var _ = ginkgo.Describe("Create admin tests", func() {
	ginkgo.It("should register the admin user", func() {
		handler := New(NewModel(
			"name",
			"lastName",
			fmt.Sprintf("%s@gmail.com", uuid.NewString()),
			"password",
			true,
		),
			logger.NewLogBuilder())

		_, err := handler.Handle()
		testAssertErrNil(err)
	})
})
