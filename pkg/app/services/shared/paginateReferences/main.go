package paginateReferences

import (
	"creatif/pkg/app/auth"
	pkg "creatif/pkg/lib"
	"creatif/pkg/lib/appErrors"
	"creatif/pkg/lib/sdk"
	"creatif/pkg/lib/storage"
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

func (c Main) Logic() (sdk.LogicView[QueryVariable], error) {
	placeholders, _ := createPlaceholdersFromModel(c.model)
	tables := getWorkingTables(c.model.StructureType)
	orderBy, direction := createFields(c.model)
	sql := createSql(c.model, tables, orderBy, direction, c.model.RelationshipType)

	var items []QueryVariable
	res := storage.Gorm().Raw(sql, placeholders).Scan(&items)
	if res.Error != nil {
		return sdk.LogicView[QueryVariable]{}, appErrors.NewDatabaseError(res.Error).AddError("Maps.Paginate.Logic", nil)
	}

	ids := sdk.Map(items, func(idx int, value QueryVariable) string {
		return value.ID
	})

	if sdk.Includes(c.model.Fields, "groups") {
		groups, err := getItemGroups(ids)
		if err != nil {
			return sdk.LogicView[QueryVariable]{}, appErrors.NewDatabaseError(err).AddError("ListItems.Paginate.Logic", nil)
		}

		resultsOfPagination := items
		for _, g := range groups {
			for i, p := range resultsOfPagination {
				if grps, ok := g[p.ID]; ok {
					p.Groups = grps
				}

				resultsOfPagination[i] = p
			}
		}
	}

	return sdk.LogicView[QueryVariable]{
		Total: 0,
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

func New(model Model, auth auth.Authentication) pkg.Job[Model, sdk.PaginationView[View], sdk.LogicView[QueryVariable]] {
	return Main{model: model, auth: auth}
}
