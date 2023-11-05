package loginEmail

import (
	"creatif/pkg/app/auth"
	"creatif/pkg/app/domain/app"
	pkg "creatif/pkg/lib"
	"creatif/pkg/lib/appErrors"
	"creatif/pkg/lib/logger"
	"creatif/pkg/lib/storage"
	"encoding/json"
	"errors"
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

func (c Main) Logic() ([]byte, error) {
	var user app.User
	if res := storage.Gorm().Where("email = ?", c.model.Email).Select("id", "key", "confirmed", "name", "last_name", "email", "created_at", "updated_at").First(&user); res.Error != nil {
		return nil, appErrors.NewAuthenticationError(res.Error)
	}

	if !user.Confirmed {
		return nil, appErrors.NewUserUnconfirmedError(errors.New("The user is not confirmed"))
	}

	authenticatedUser := auth.NewAuthenticatedUser(user.ID, user.Name, user.LastName, user.Email, user.CreatedAt, user.UpdatedAt, time.Now())
	serializedUser, err := json.Marshal(&authenticatedUser)
	if err != nil {
		return nil, appErrors.NewApplicationError(err)
	}

	if len(user.Key) != 32 {
		// TODO: send immediate slack message here
		c.logBuilder.Add("login.invalidUserKey", "Key should have length of 32 characters.")
		return nil, appErrors.NewApplicationError(err)
	}

	var key [32]byte
	for i, k := range user.Key {
		key[i] = byte(k)
	}

	encryptedUser, err := encrypt(serializedUser, &key)
	if err != nil {
		// TODO: send immediate slack message here
		c.logBuilder.Add("login.cannotEncryptKey", "Key should have length of 32 characters.")
		return nil, appErrors.NewApplicationError(err)
	}

	return encryptedUser, nil
}

func (c Main) Handle() ([]byte, error) {
	if err := c.Validate(); err != nil {
		return nil, err
	}

	if err := c.Authenticate(); err != nil {
		return nil, err
	}

	if err := c.Authorize(); err != nil {
		return nil, err
	}

	model, err := c.Logic()

	if err != nil {
		return nil, err
	}

	return model, nil
}

func New(model Model, auth auth.Authentication, logBuilder logger.LogBuilder) pkg.Job[Model, []byte, []byte] {
	logBuilder.Add("loginEmail", "Created")
	return Main{model: model, logBuilder: logBuilder, auth: auth}
}
