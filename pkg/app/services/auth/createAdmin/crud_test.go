package createAdmin

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Create admin tests", func() {
	ginkgo.It("should register the admin user", ginkgo.Label("admin"), func() {
		handler := New(NewModel(
			"name",
			"lastName",
			fmt.Sprintf("%s@gmail.com", uuid.NewString()),
			"password",
		))

		_, err := handler.Handle()
		gomega.Expect(err).Should(gomega.BeNil())
	})

	ginkgo.It("should not be able to register multiple admin users", ginkgo.Label("admin"), func() {
		handler := New(NewModel(
			"name",
			"lastName",
			fmt.Sprintf("%s@gmail.com", uuid.NewString()),
			"password",
		))

		_, err := handler.Handle()
		gomega.Expect(err).Should(gomega.BeNil())

		handler = New(NewModel(
			"name",
			"lastName",
			fmt.Sprintf("%s@gmail.com", uuid.NewString()),
			"someotherpassword",
		))

		_, err = handler.Handle()
		gomega.Expect(err).ShouldNot(gomega.BeNil())
	})

	ginkgo.It("password that is less than 8 and more than 20 characters should not be accepted", ginkgo.Label("admin"), func() {
		handler := New(NewModel(
			"name",
			"lastName",
			fmt.Sprintf("%s@gmail.com", uuid.NewString()),
			"pass",
		))

		_, err := handler.Handle()
		gomega.Expect(err).ShouldNot(gomega.BeNil())

		handler = New(NewModel(
			"name",
			"lastName",
			fmt.Sprintf("%s@gmail.com", uuid.NewString()),
			"passwordthanismorethantwentycharacters",
		))

		_, err = handler.Handle()
		gomega.Expect(err).ShouldNot(gomega.BeNil())
	})
})
