package appendToList

import (
	"creatif/pkg/app/auth"
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/app/services/locales"
	pkg "creatif/pkg/lib"
	"creatif/pkg/lib/appErrors"
	"creatif/pkg/lib/logger"
	"creatif/pkg/lib/storage"
	"errors"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Main struct {
	model      Model
	logBuilder logger.LogBuilder
	auth       auth.Authentication
}

func (c Main) Validate() error {
	c.logBuilder.Add("appendToList", "Validating...")
	if errs := c.model.Validate(); errs != nil {
		return appErrors.NewValidationError(errs)
	}

	c.logBuilder.Add("appendToList", "Validated.")
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
	var list declarations.List
	if res := storage.Gorm().Where("id = ? OR short_id = ? OR name = ?", c.model.Name, c.model.Name, c.model.Name).Select("id", "serial", "project_id", "name", "created_at", "updated_at").First(&list); res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return declarations.List{}, appErrors.NewNotFoundError(res.Error).AddError("appendToList.Logic", nil)
		}

		return declarations.List{}, appErrors.NewDatabaseError(res.Error).AddError("appendToList.Logic", nil)
	}

	for _, v := range c.model.Variables {
		pqGroups := pq.StringArray{}
		if v.Groups == nil {
			v.Groups = pq.StringArray{}
		}

		for _, k := range v.Groups {
			pqGroups = append(pqGroups, k)
		}

		v.Groups = pqGroups
	}

	highestIndex, err := getHighestIndex(list.ID)
	if err != nil {
		return declarations.List{}, err
	}

	listVariables := make([]declarations.ListVariable, len(c.model.Variables))
	for i := 0; i < len(c.model.Variables); i++ {
		if c.model.Variables[i].Locale == "" {
			c.model.Variables[i].Locale = "eng"
		}

		localeID, _ := locales.GetIDWithAlpha(c.model.Variables[i].Locale)
		v := c.model.Variables[i]
		listVariables[i] = declarations.NewListVariable(list.ID, localeID, v.Name, v.Behaviour, v.Metadata, v.Groups, v.Value)
		listVariables[i].Index = float64(highestIndex + 1000)
		highestIndex += 1000
	}

	if err := storage.Transaction(func(tx *gorm.DB) error {
		if res := tx.Create(&listVariables); res.Error != nil {
			c.logBuilder.Add("appendToList", res.Error.Error())
			return res.Error
		}

		return nil
	}); err != nil {
		return declarations.List{}, appErrors.NewApplicationError(err)
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

func New(model Model, auth auth.Authentication, logBuilder logger.LogBuilder) pkg.Job[Model, View, declarations.List] {
	logBuilder.Add("appendToList", "Created")
	return Main{model: model, logBuilder: logBuilder, auth: auth}
}
