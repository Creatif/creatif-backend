package adminExists

import (
	"creatif/pkg/app/services/auth/createAdmin"
	"fmt"
	"github.com/google/uuid"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Has admin tests", func() {
	ginkgo.It("should check if the admin user exists", func() {
		handler := createAdmin.New(createAdmin.NewModel(
			"otherName",
			"lastName",
			fmt.Sprintf("%s@gmail.com", uuid.NewString()),
			"sfdsafdsafdsafdsa",
		))

		_, err := handler.Handle()
		gomega.Expect(err).Should(gomega.BeNil())

		adminExistsHandler := New()
		exists, err := adminExistsHandler.Handle()
		gomega.Expect(err).Should(gomega.BeNil())
		gomega.Expect(exists).Should(gomega.BeTrue())
	})
})
