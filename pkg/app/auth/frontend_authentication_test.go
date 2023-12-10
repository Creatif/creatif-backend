package auth

import (
	"creatif/pkg/app/domain/app"
	"creatif/pkg/app/services/auth"
	"creatif/pkg/lib/logger"
	"creatif/pkg/lib/storage"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"time"
)

var _ = ginkgo.Describe("Frontend email login/authentication", func() {
	ginkgo.It("should login a user with email", func() {
		user := app.NewUser("name", "lastName", "marioskrlec222@gmail.com", "digital", auth.EmailProvider, true, true)
		storage.Gorm().Create(&user)

		l := logger.NewLogBuilder()
		var key [32]byte
		for i, v := range user.Key {
			key[i] = byte(v)
		}

		authenticatedUser := NewAuthenticatedUser(user.ID, user.Name, user.LastName, user.Email, user.CreatedAt, user.UpdatedAt, time.Now(), "", "")
		loginer := NewEmailLogin(authenticatedUser, key, l)
		token, err := loginer.Login()
		testAssertErrNil(err)
		gomega.Expect(token).ShouldNot(gomega.BeEmpty())

		l.Flush("")
	})

	ginkgo.It("should login an authenticate a user", func() {
		user := app.NewUser("name", "lastName", "marioskrlec222@gmail.com", "digital", auth.EmailProvider, true, true)
		storage.Gorm().Create(&user)

		var key [32]byte
		for i, v := range user.Key {
			key[i] = byte(v)
		}
		l := logger.NewLogBuilder()
		authenticatedUser := NewAuthenticatedUser(user.ID, user.Name, user.LastName, user.Email, user.CreatedAt, user.UpdatedAt, time.Now(), "", "")
		loginer := NewEmailLogin(authenticatedUser, key, l)
		token, err := loginer.Login()
		testAssertErrNil(err)
		gomega.Expect(token).ShouldNot(gomega.BeEmpty())

		authentication := NewFrontendAuthentication(token, l)

		err = authentication.Authenticate()
		gomega.Expect(err).Should(gomega.BeNil())

		gomega.Expect(user.ID).Should(gomega.Equal(authentication.User().ID))

		l.Flush("")
	})
})
