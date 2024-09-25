package createAdmin

import (
	"creatif/pkg/app/domain/app"
	auth2 "creatif/pkg/app/services/auth"
	pkg "creatif/pkg/lib"
	"creatif/pkg/lib/appErrors"
	"creatif/pkg/lib/storage"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type Main struct {
	model Model
}

func (c Main) Validate() error {
	if errs := c.model.Validate(); errs != nil {
		return appErrors.NewValidationError(errs)
	}

	var user app.User
	if res := storage.Gorm().Where("is_admin = ?", true).Select("id").First(&user); res.Error != nil {
		if !errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return appErrors.NewValidationError(map[string]string{
				"admin": fmt.Sprintf("There is already an admin for this site"),
			})
		} else if res.Error != nil {
			// blank
		}
	}

	return nil
}

func (c Main) Authenticate() error {
	return nil
}

func (c Main) Authorize() error {
	return nil
}

func (c Main) Logic() (interface{}, error) {
	pass, err := hashPassword(c.model.Password)
	if err != nil {
		return nil, appErrors.NewUnexpectedError(err)
	}

	user := app.NewUser(c.model.Name, c.model.LastName, c.model.Email, pass, auth2.EmailProvider, true, true)

	if res := storage.Gorm().Create(&user); res.Error != nil {
		return nil, appErrors.NewDatabaseError(res.Error)
	}

	return nil, nil
}

func (c Main) Handle() (interface{}, error) {
	if err := c.Validate(); err != nil {
		return nil, err
	}

	if err := c.Authenticate(); err != nil {
		return nil, err
	}

	if err := c.Authorize(); err != nil {
		return nil, err
	}

	_, err := c.Logic()

	if err != nil {
		return nil, err
	}

	return nil, nil
}

func New(model Model) pkg.Job[Model, interface{}, interface{}] {
	return Main{model: model}
}
