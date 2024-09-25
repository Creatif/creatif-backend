package createList

import (
	"creatif/pkg/app/auth"
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/app/services/locales"
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

	if err := validateListWithNameExists(c.model.Name, c.model.ProjectID); err != nil {
		return err
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

func (c Main) Logic() (declarations.List, error) {
	list := declarations.NewList(c.model.ProjectID, c.model.Name)
	if err := storage.Transaction(func(tx *gorm.DB) error {
		if res := tx.Create(&list); res.Error != nil {
			return res.Error
		}

		highestIndex, err := getHighestIndex(list.ID)
		if err != nil {
			return err
		}

		if len(c.model.Variables) > 0 {
			listVariables := make([]declarations.ListVariable, len(c.model.Variables))
			for i := 0; i < len(c.model.Variables); i++ {
				localeID, _ := locales.GetIDWithAlpha(c.model.Variables[i].Locale)
				v := c.model.Variables[i]
				listVariables[i] = declarations.NewListVariable(list.ID, localeID, v.Name, v.Behaviour, v.Metadata, v.Value)
				listVariables[i].Index = float64(highestIndex) + 1000
				highestIndex += 1000
			}

			if res := tx.Create(&listVariables); res.Error != nil {
				return res.Error
			}

			tx.Model(&declarations.List{}).Where("id = ?", list.ID).Update("serial", highestIndex)
		}

		return nil
	}); err != nil {
		return declarations.List{}, appErrors.NewDatabaseError(err).AddError("createList.Logic", nil)
	}

	return list, nil
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

func New(model Model, auth auth.Authentication) pkg.Job[Model, View, declarations.List] {
	return Main{model: model, auth: auth}
}
