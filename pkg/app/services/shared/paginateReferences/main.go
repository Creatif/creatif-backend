package paginateReferences

import (
	"creatif/pkg/app/auth"
	"creatif/pkg/app/domain/declarations"
	pkg "creatif/pkg/lib"
	"creatif/pkg/lib/appErrors"
	"creatif/pkg/lib/logger"
	"creatif/pkg/lib/sdk"
	"creatif/pkg/lib/storage"
)

type Main struct {
	model      Model
	logBuilder logger.LogBuilder
	auth       auth.Authentication
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

func (c Main) Logic() (sdk.LogicView[declarations.MapVariable], error) {
	placeholders, countPlaceholders := createPlaceholdersFromModel(c.model)
	tables := getWorkingTables(c.model.StructureType)
	orderBy, direction := createFields(c.model)
	sql := createSql(c.model, tables, orderBy, direction, c.model.RelationshipType)

	var items []declarations.MapVariable
	res := storage.Gorm().Raw(sql, placeholders).Scan(&items)
	if res.Error != nil {
		c.logBuilder.Add("paginateMapVariables", res.Error.Error())
		return sdk.LogicView[declarations.MapVariable]{}, appErrors.NewDatabaseError(res.Error).AddError("Maps.Paginate.Logic", nil)
	}

	countSql := createCountSql(c.model, tables, c.model.RelationshipType)

	var count int64
	res = storage.Gorm().Raw(countSql, countPlaceholders).Scan(&count)
	if res.Error != nil {
		c.logBuilder.Add("paginateMapVariables", res.Error.Error())
		return sdk.LogicView[declarations.MapVariable]{}, appErrors.NewDatabaseError(res.Error).AddError("paginateMapVariable.Logic", nil)
	}

	return sdk.LogicView[declarations.MapVariable]{
		Total: count,
		Data:  items,
	}, nil
}

func (c Main) Handle() (sdk.PaginationView[View], error) {
	if err := c.Validate(); err != nil {
		return sdk.PaginationView[View]{}, err
	}

	if err := c.Authenticate(); err != nil {
		return sdk.PaginationView[View]{}, err
	}

	if err := c.Authorize(); err != nil {
		return sdk.PaginationView[View]{}, err
	}

	model, err := c.Logic()

	if err != nil {
		return sdk.PaginationView[View]{}, err
	}

	items, err := newView(model.Data)
	if err != nil {
		return sdk.PaginationView[View]{}, appErrors.NewApplicationError(err).AddError("ListItems.Paginate.Handle", nil)
	}

	return sdk.PaginationView[View]{
		Total: model.Total,
		Page:  c.model.Page,
		Data:  items,
	}, nil
}

func New(model Model, auth auth.Authentication, logBuilder logger.LogBuilder) pkg.Job[Model, sdk.PaginationView[View], sdk.LogicView[declarations.MapVariable]] {
	logBuilder.Add("paginateMapVariables", "Created")
	return Main{model: model, logBuilder: logBuilder, auth: auth}
}
