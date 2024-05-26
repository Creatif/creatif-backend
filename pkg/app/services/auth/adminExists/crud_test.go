package adminExists

import (
	"creatif/pkg/app/services/auth/createAdmin"
	"creatif/pkg/lib/logger"
	"fmt"
	"github.com/google/uuid"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Has admin tests", func() {
	ginkgo.It("should check if the admin user exists", func() {
		handler := createAdmin.New(createAdmin.NewModel(
			"name",
			"lastName",
			fmt.Sprintf("%s@gmail.com", uuid.NewString()),
			"password",
		),
			logger.NewLogBuilder())

		_, err := handler.Handle()
		gomega.Expect(err).Should(gomega.BeNil())

		adminExistsHandler := New(logger.NewLogBuilder())
		exists, err := adminExistsHandler.Handle()
		gomega.Expect(err).Should(gomega.BeNil())
		gomega.Expect(exists).Should(gomega.BeTrue())
	})
})
