package loginEmail

import (
	"creatif/pkg/app/auth"
	"creatif/pkg/app/domain/app"
	pkg "creatif/pkg/lib"
	"creatif/pkg/lib/appErrors"
	"creatif/pkg/lib/logger"
	"creatif/pkg/lib/storage"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type Main struct {
	model      Model
	logBuilder logger.LogBuilder
	auth       auth.Authentication
}

func (c Main) Validate() error {
	c.logBuilder.Add("loginEmail", "Validating...")
	if errs := c.model.Validate(); errs != nil {
		return appErrors.NewValidationError(errs)
	}

	c.logBuilder.Add("loginEmail", "Validated.")
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
	if res := storage.Gorm().Where("email = ?", c.model.Email).Select("id", "key", "confirmed", "password", "name", "last_name", "email", "created_at", "updated_at").First(&user); res.Error != nil {
		return "", appErrors.NewAuthenticationError(res.Error)
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(c.model.Password))
	if err != nil {
		return "", appErrors.NewAuthenticationError(err)
	}

	if !user.Confirmed {
		return "", appErrors.NewUserUnconfirmedError(errors.New("The user is not confirmed"))
	}
	var key [32]byte
	for i, v := range user.Key {
		key[i] = byte(v)
	}

	authenticatedUser := auth.NewAuthenticatedUser(user.ID, user.Name, user.LastName, user.Email, user.CreatedAt, user.UpdatedAt, time.Now())
	return auth.NewEmailLogin(authenticatedUser, key, c.logBuilder).Login()
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

func New(model Model, auth auth.Authentication, logBuilder logger.LogBuilder) pkg.Job[Model, string, string] {
	logBuilder.Add("loginEmail", "Created")
	return Main{model: model, logBuilder: logBuilder, auth: auth}
}
