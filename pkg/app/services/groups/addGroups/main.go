package addGroups

import (
	"creatif/pkg/app/auth"
	"creatif/pkg/app/domain/declarations"
	pkg "creatif/pkg/lib"
	"creatif/pkg/lib/appErrors"
	"creatif/pkg/lib/logger"
	"creatif/pkg/lib/storage"
	"fmt"
	"gorm.io/gorm"
)

type Main struct {
	model      Model
	logBuilder logger.LogBuilder
	auth       auth.Authentication
}

func (c Main) Validate() error {
	c.logBuilder.Add("addToList", "Validating...")
	if errs := c.model.Validate(); errs != nil {
		return appErrors.NewValidationError(errs)
	}
	c.logBuilder.Add("addToList", "Validated.")

	return nil
}

func (c Main) Authenticate() error {
	if err := c.auth.Authenticate(); err != nil {
		return appErrors.NewAuthenticationError(err)
	}

	return nil
}

func (c Main) Authorize() error {
	return nil
}

func (c Main) Logic() ([]declarations.Group, error) {
	groups := make([]declarations.Group, len(c.model.Groups))

	if transactionErr := storage.Transaction(func(tx *gorm.DB) error {
		if res := tx.Exec(fmt.Sprintf("DELETE FROM %s WHERE project_id = ?", (declarations.Group{}).TableName()), c.model.ProjectID); res.Error != nil {
			return res.Error
		}

		if len(c.model.Groups) > 0 {
			for i, g := range c.model.Groups {
				groups[i] = declarations.NewGroup(c.model.ProjectID, g)
			}

			if res := tx.Create(&groups); res.Error != nil {
				return res.Error
			}
		}

		return nil
	}); transactionErr != nil {
		return []declarations.Group{}, appErrors.NewApplicationError(transactionErr)
	}
	return groups, nil
}

func (c Main) Handle() ([]string, error) {
	if err := c.Validate(); err != nil {
		return []string{}, err
	}

	if err := c.Authenticate(); err != nil {
		return []string{}, err
	}

	if err := c.Authorize(); err != nil {
		return []string{}, err
	}

	model, err := c.Logic()

	if err != nil {
		return []string{}, err
	}

	return newView(model), nil
}

func New(model Model, auth auth.Authentication, logBuilder logger.LogBuilder) pkg.Job[Model, []string, []declarations.Group] {
	logBuilder.Add("addToList", "Created")
	return Main{model: model, logBuilder: logBuilder, auth: auth}
}
