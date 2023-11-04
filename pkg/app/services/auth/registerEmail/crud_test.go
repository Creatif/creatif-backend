package registerEmail

import (
	"creatif/pkg/app/auth"
	"creatif/pkg/lib/logger"
	"fmt"
	"github.com/google/uuid"
	"github.com/onsi/ginkgo/v2"
)

var _ = ginkgo.Describe("Register account tests", func() {
	ginkgo.It("should register a user", func() {
		handler := New(NewModel(
			"name",
			"lastName",
			fmt.Sprintf("%s@gmail.com", uuid.NewString()),
			"password",
			true,
		),
			auth.NewNoopAuthentication(),
			logger.NewLogBuilder())

		_, err := handler.Handle()
		testAssertErrNil(err)
	})
})
