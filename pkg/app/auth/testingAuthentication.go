package auth

import (
	"creatif/pkg/app/domain/app"
	auth2 "creatif/pkg/app/services/auth"
	storage2 "creatif/pkg/lib/storage"
	"fmt"
	"github.com/google/uuid"
	"github.com/onsi/ginkgo/v2"
	"time"
)

var testSessionUser AuthenticatedUser

type testingAuthentication struct {
	shouldCreateUser bool
	projectId        string
	TestSessionUser  *AuthenticatedUser
}

func (a *testingAuthentication) Authenticate() error {
	return nil
}

func (a *testingAuthentication) User() AuthenticatedUser {
	if a.TestSessionUser != nil {
		return *a.TestSessionUser
	}

	user := app.NewUser(uuid.NewString(), uuid.NewString(), fmt.Sprintf("%s@gmail.com", uuid.New().String()), "password", auth2.EmailProvider, true, true)
	res := storage2.Gorm().Create(&user)
	if res.Error != nil {
		ginkgo.Fail(res.Error.Error())
	}

	testSessionUser = AuthenticatedUser{
		ID:        user.ID,
		ProjectID: a.projectId,
		Name:      user.Name,
		LastName:  user.LastName,
		Email:     user.Email,
		Refresh:   time.Time{},
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}

	return testSessionUser
}

func (a *testingAuthentication) Refresh() (string, error) {
	return "", nil
}

func (a *testingAuthentication) Logout(cb func()) {
}

func (a *testingAuthentication) ShouldRefresh() bool {
	return false
}

func NewTestingAuthentication(shouldCreateUser bool, projectId string) Authentication {
	return &testingAuthentication{shouldCreateUser: shouldCreateUser, projectId: projectId}
}
