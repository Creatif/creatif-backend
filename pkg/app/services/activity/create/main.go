package create

import "C"
import (
	"creatif/pkg/app/auth"
	"creatif/pkg/app/domain/app"
	pkg "creatif/pkg/lib"
	"creatif/pkg/lib/appErrors"
	"creatif/pkg/lib/storage"
	"gorm.io/gorm"
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
	if err := c.auth.Authenticate(); err != nil {
		return appErrors.NewAuthenticationError(err)
	}

	return nil
}

func (c Main) Authorize() error {
	return nil
}

func (c Main) Logic() (LogicModel, error) {
	var activity app.Activity
	if transactionError := storage.Gorm().Transaction(func(tx *gorm.DB) error {
		count, err := getActivityCount(c.model.ProjectID)
		if err != nil {
			return err
		}

		dataQuery, err := getLastActivityDataQuery(c.model.ProjectID)
		if err != nil {
			return err
		}

		if dataQuery.ID != "" {
			shouldWriteNewActivity, err := decideToCreateNewActivity(dataQuery.Data, c.model.Data)
			if err != nil {
				return appErrors.NewApplicationError(err)
			}

			if !shouldWriteNewActivity {
				return nil
			}
		}

		if count >= 10 {
			if err := deleteActivity(c.model.ProjectID, dataQuery.ID); err != nil {
				return err
			}
		}

		activity = app.NewActivity(c.model.ProjectID, c.model.Data)
		if res := storage.Gorm().Create(&activity); res.Error != nil {
			return res.Error
		}

		return nil
	}); transactionError != nil {
		return LogicModel{}, appErrors.NewDatabaseError(transactionError)
	}

	return LogicModel{
		ID:            activity.ID,
		Data:          activity.Data,
		CreatedAt:     activity.CreatedAt,
		ActivityAdded: activity.ID != "",
	}, nil
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

func New(model Model, auth auth.Authentication) pkg.Job[Model, View, LogicModel] {
	return Main{model: model, auth: auth}
}
