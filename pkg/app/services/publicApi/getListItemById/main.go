package getListItemById

import "C"
import (
	"creatif/pkg/app/auth"
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
	c.logBuilder.Add("getVersions", "Validating...")
	if errs := c.model.Validate(); errs != nil {
		return appErrors.NewValidationError(errs)
	}

	c.logBuilder.Add("getVersions", "Validated")
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
	var mapItem MapItem
	res := storage.Gorm().Raw(getListItemSql(), c.model.ProjectID, c.model.VersionName, c.model.ItemID).Scan(&mapItem)
	if res.Error != nil {
		return LogicModel{}, appErrors.NewApplicationError(res.Error)
	}

	if res.RowsAffected == 0 {
		return LogicModel{}, appErrors.NewNotFoundError(errors.New("This item does not seem to exist"))
	}

	var connections []ConnectionMapItem
	res = storage.Gorm().Raw(getConnectionsMapSql(), c.model.ProjectID, c.model.VersionName, c.model.ItemID, c.model.ProjectID, c.model.VersionName, c.model.ItemID).Scan(&connections)
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
	logBuilder.Add("getVersions", "Created")
	return Main{model: model, logBuilder: logBuilder, auth: auth}
}
