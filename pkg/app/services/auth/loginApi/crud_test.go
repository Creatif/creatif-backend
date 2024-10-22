package loginApi

import (
	"creatif/pkg/app/auth"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Login account tests", func() {
	ginkgo.It("password that is less than 8 and more than 20 characters should not be accepted", ginkgo.Label("login"), func() {
		handler := New(NewModel(
			"mario@gmail.com",
			"pass",
		), auth.NewNoopAuthentication())

		_, err := handler.Handle()
		gomega.Expect(err).ShouldNot(gomega.BeNil())

		handler = New(NewModel(
			"mario@gmail.com",
			"passwordthanismorethantwentycharacters",
		), auth.NewNoopAuthentication())

		_, err = handler.Handle()
		gomega.Expect(err).ShouldNot(gomega.BeNil())
	})
})
