package loginApi

import (
	"creatif/pkg/app/auth"
	"creatif/pkg/app/domain/app"
	pkg "creatif/pkg/lib"
	"creatif/pkg/lib/appErrors"
	"creatif/pkg/lib/storage"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type Main struct {
	model Model
	auth  auth.Authentication
}

func (c Main) Validate() error {
	if errs := c.model.Validate(); errs != nil {
		return appErrors.NewValidationError(errs)
	}

	return nil
}

func (c Main) Authenticate() error {
	return nil
}

func (c Main) Authorize() error {
	return nil
}

func (c Main) Logic() (string, error) {
	var user app.User
	res := storage.Gorm().Raw(fmt.Sprintf(`
SELECT 
	u.id,
	u.key,
	u.confirmed,
	u.name,
	u.password,
	u.last_name,
	u.is_admin,
	u.email,
	u.created_at,
	u.updated_at
FROM %s AS u WHERE u.email = ?
`, (app.User{}).TableName()), c.model.Email).Scan(&user)

	if res.Error != nil {
		return "", appErrors.NewAuthenticationError(errors.New("Unauthenticated"))
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(c.model.Password))
	if err != nil {
		return "", appErrors.NewAuthenticationError(errors.New("Email or password are invalid"))
	}

	if !user.Confirmed {
		return "", appErrors.NewAuthenticationError(errors.New("The user is not confirmed"))
	}

	var key [32]byte
	for i, v := range user.Key {
		key[i] = byte(v)
	}

	authenticatedUser := auth.NewAuthenticatedUser(user.ID, user.Name, user.LastName, user.Email, user.CreatedAt, user.UpdatedAt, time.Now())
	return auth.NewApiLogin(authenticatedUser, key).Login()
}

func (c Main) Handle() (string, error) {
	if err := c.Validate(); err != nil {
		return "", err
	}

	if err := c.Authenticate(); err != nil {
		return "", err
	}

	if err := c.Authorize(); err != nil {
		return "", err
	}

	model, err := c.Logic()

	if err != nil {
		return "", err
	}

	return model, nil
}

func New(model Model, auth auth.Authentication) pkg.Job[Model, string, string] {
	return Main{model: model, auth: auth}
}
