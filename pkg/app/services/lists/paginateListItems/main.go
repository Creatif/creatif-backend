package paginateListItems

import (
	"creatif/pkg/app/auth"
	pkg "creatif/pkg/lib"
	"creatif/pkg/lib/appErrors"
	"creatif/pkg/lib/sdk"
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
	offset := (c.model.Page - 1) * c.model.Limit
	queryPlaceholders := createQueryPlaceholders(
		c.model.ProjectID,
		c.model.ListName,
		offset,
		c.model.Limit,
		c.model.Groups,
		c.model.Behaviour,
		c.model.Search,
	)

	countPlaceholders := createCountPlaceholders(
		c.model.ProjectID,
		c.model.ListName,
		c.model.Groups,
		c.model.Behaviour,
		c.model.Search,
	)

	defs := createDefaults(c.model.OrderBy, c.model.OrderDirection)
	sq := createSubQueries(
		c.model.Behaviour,
		c.model.Locales,
		c.model.Groups,
		c.model.Search,
	)

	paginationResult, countResult := runQueriesConcurrently(queryPlaceholders, countPlaceholders, sq, defs)

	if paginationResult.error != nil {
		return sdk.LogicView[QueryVariable]{}, appErrors.NewDatabaseError(paginationResult.error).AddError("ListItems.Paginate.Logic", nil)
	}

	if countResult.error != nil {
		return sdk.LogicView[QueryVariable]{}, appErrors.NewDatabaseError(countResult.error).AddError("ListItems.Paginate.Logic", nil)
	}

	ids := sdk.Map(paginationResult.result, func(idx int, value QueryVariable) string {
		return value.ID
	})

	if sdk.Includes(c.model.Fields, "groups") {
		groups, err := getItemGroups(ids)
		if err != nil {
			return sdk.LogicView[QueryVariable]{}, appErrors.NewDatabaseError(err).AddError("ListItems.Paginate.Logic", nil)
		}

		resultsOfPagination := paginationResult.result
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
		Total: countResult.result,
		Data:  paginationResult.result,
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
