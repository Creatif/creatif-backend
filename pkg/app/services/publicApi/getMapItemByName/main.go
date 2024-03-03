package getMapItemByName

import "C"
import (
	"creatif/pkg/app/auth"
	"creatif/pkg/app/services/locales"
	pkg "creatif/pkg/lib"
	"creatif/pkg/lib/appErrors"
	"creatif/pkg/lib/logger"
	"creatif/pkg/lib/storage"
	"errors"
)

type Main struct {
	model      Model
	logBuilder logger.LogBuilder
	auth       auth.Authentication
}

func (c Main) Validate() error {
	c.logBuilder.Add("getMapItemByName", "Validating...")
	if errs := c.model.Validate(); errs != nil {
		return appErrors.NewValidationError(errs)
	}

	c.logBuilder.Add("getMapItemByName", "Validated")
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
	var locale string
	l, err := locales.GetIDWithAlpha(c.model.Locale)
	if err != nil {
		l, _ := locales.GetIDWithAlpha("eng")
		locale = l
	} else {
		locale = l
	}

	var mapItem Item
	res := storage.Gorm().Raw(getItemSql(), c.model.ProjectID, c.model.VersionName, c.model.Name, locale).Scan(&mapItem)
	if res.Error != nil {
		return LogicModel{}, appErrors.NewApplicationError(res.Error)
	}

	if res.RowsAffected == 0 {
		return LogicModel{}, appErrors.NewNotFoundError(errors.New("This item does not seem to exist"))
	}

	var connections []ConnectionItem
	res = storage.Gorm().Raw(getConnectionsSql(), c.model.ProjectID, c.model.VersionName, mapItem.ItemID, c.model.ProjectID, c.model.VersionName, mapItem.ItemID).Scan(&connections)
	if res.Error != nil {
		return LogicModel{}, appErrors.NewApplicationError(res.Error)
	}

	return LogicModel{
		Item:        mapItem,
		Connections: connections,
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

func New(model Model, auth auth.Authentication, logBuilder logger.LogBuilder) pkg.Job[Model, View, LogicModel] {
	logBuilder.Add("getMapItemByName", "Created")
	return Main{model: model, logBuilder: logBuilder, auth: auth}
}