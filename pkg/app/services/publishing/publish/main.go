package publish

import "C"
import (
	"context"
	"creatif/pkg/app/auth"
	"creatif/pkg/app/domain/published"
	pkg "creatif/pkg/lib"
	"creatif/pkg/lib/appErrors"
	"creatif/pkg/lib/logger"
	"creatif/pkg/lib/storage"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Results struct {
	Errors []error
}

type Main struct {
	model      Model
	logBuilder logger.LogBuilder
	auth       auth.Authentication
}

func (c Main) Validate() error {
	c.logBuilder.Add("publish", "Validating...")
	if errs := c.model.Validate(); errs != nil {
		return appErrors.NewValidationError(errs)
	}

	if err := validateVersionNameExists(c.model.ProjectID, c.model.Name); err != nil {
		return err
	}

	c.logBuilder.Add("publish", "Validated")
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

func (c Main) Logic() (published.Version, error) {
	name := c.model.Name
	if name == "" {
		name = uuid.NewString()
	}

	version := published.NewVersion(c.model.ProjectID, name, false)
	if transactionError := storage.Transaction(func(tx *gorm.DB) error {
		if res := tx.Exec(fmt.Sprintf("UPDATE %s SET is_production_version = false WHERE project_id = ? AND is_production_version = true", (published.Version{}).TableName()), c.model.ProjectID); res.Error != nil {
			return res.Error
		}

		if res := tx.Create(&version); res.Error != nil {
			return res.Error
		}

		listCtx, listCancel := context.WithTimeout(context.Background(), 45*time.Second)
		mapCtx, mapCancel := context.WithTimeout(context.Background(), 45*time.Second)
		refCtx, refCancel := context.WithTimeout(context.Background(), 45*time.Second)
		defer listCancel()
		defer mapCancel()
		defer refCancel()
		if err := publishLists(tx, c.model.ProjectID, version.ID, listCtx); err != nil {
			return err
		}
		if err := publishMaps(tx, c.model.ProjectID, version.ID, mapCtx); err != nil {
			return err
		}
		if err := publishReferences(tx, c.model.ProjectID, version.ID, refCtx); err != nil {
			return err
		}

		return nil
	}); transactionError != nil {
		return published.Version{}, appErrors.NewApplicationError(transactionError)
	}

	return version, nil
}

func (c Main) Handle() (View, error) {
	if err := c.Validate(); err != nil {
		return View{}, err
	}

	if err := c.Authenticate(); err != nil {
		return View{}, err
	}

	if err := c.Authorize(); err != nil {
		return View{}, err
	}

	model, err := c.Logic()

	if err != nil {
		return View{}, err
	}

	return newView(model), nil
}

func New(model Model, auth auth.Authentication, logBuilder logger.LogBuilder) pkg.Job[Model, View, published.Version] {
	logBuilder.Add("publish", "Created")
	return Main{model: model, logBuilder: logBuilder, auth: auth}
}
