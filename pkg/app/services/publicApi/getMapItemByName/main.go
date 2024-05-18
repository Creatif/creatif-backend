package getMapItemByName

import (
	"creatif/pkg/app/auth"
	"creatif/pkg/app/domain/published"
	"creatif/pkg/app/services/locales"
	"creatif/pkg/app/services/publicApi/publicApiError"
	pkg "creatif/pkg/lib"
	"creatif/pkg/lib/appErrors"
	"creatif/pkg/lib/logger"
	"creatif/pkg/lib/storage"
	"fmt"
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
	var version published.Version
	res := storage.Gorm().Raw(fmt.Sprintf("SELECT * FROM %s WHERE project_id = ? AND is_production_version = true", (published.Version{}).TableName()), c.model.ProjectID).Scan(&version)
	if res.Error != nil {
		return LogicModel{}, publicApiError.NewError("getMapItemByName", map[string]string{
			"error": res.Error.Error(),
		}, publicApiError.DatabaseError)
	}

	if res.RowsAffected == 0 {
		return LogicModel{}, publicApiError.NewError("getMapItemByName", map[string]string{
			"versionNotFound": "Production version has not been found.",
		}, publicApiError.NotFoundError)
	}

	placeholders := make(map[string]interface{})
	placeholders["projectId"] = c.model.ProjectID
	placeholders["versionName"] = version.Name
	placeholders["structureName"] = c.model.StructureName
	placeholders["variableName"] = c.model.Name

	var locale string
	if c.model.Locale != "" {
		l, _ := locales.GetIDWithAlpha(c.model.Locale)
		locale = l
		placeholders["localeId"] = l
	}

	var mapItem Item
	res = storage.Gorm().Raw(getItemSql(locale), placeholders).Scan(&mapItem)
	if res.Error != nil {
		return LogicModel{}, publicApiError.NewError("getMapItemByName", map[string]string{
			"error": res.Error.Error(),
		}, publicApiError.DatabaseError)
	}

	if res.RowsAffected == 0 {
		return LogicModel{}, publicApiError.NewError("getMapItemByName", map[string]string{
			"notFound": "Item has not been found.",
		}, publicApiError.NotFoundError)
	}

	var connections []ConnectionItem
	res = storage.Gorm().Raw(getConnectionsSql(), c.model.ProjectID, version.Name, mapItem.ItemID, c.model.ProjectID, version.Name, mapItem.ItemID).Scan(&connections)
	if res.Error != nil {
		return LogicModel{}, publicApiError.NewError("getMapItemByName", map[string]string{
			"error": res.Error.Error(),
		}, publicApiError.DatabaseError)
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
