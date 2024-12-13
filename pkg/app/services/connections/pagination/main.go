package pagination

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

func (c Main) Logic() ([]QueryVariable, error) {
	queryPlaceholders := createQueryPlaceholders(
		c.model.ProjectID,
		c.model.StructureID,
		c.model.ParentVariableID,
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

	conns, err := getConnections(c.model.StructureType, queryPlaceholders, sq, defs)
	if err != nil {
		return nil, appErrors.NewDatabaseError(err).AddError("ListItems.Paginate.Logic", nil)
	}

	ids := sdk.Map(conns, func(idx int, value QueryVariable) string {
		return value.ID
	})

	if sdk.Includes(c.model.Fields, "groups") {
		groups, err := getItemGroups(ids)
		if err != nil {
			return nil, appErrors.NewDatabaseError(err).AddError("ListItems.Paginate.Logic", nil)
		}

		resultsOfPagination := conns
		for _, g := range groups {
			for i, p := range resultsOfPagination {
				if grps, ok := g[p.ID]; ok {
					p.Groups = grps
				}

				resultsOfPagination[i] = p
			}
		}
	}

	return conns, nil
}

func (c Main) Handle() ([]View, error) {
	if err := c.Validate(); err != nil {
		return nil, err
	}

	if err := c.Authenticate(); err != nil {
		return nil, err
	}

	if err := c.Authorize(); err != nil {
		return nil, err
	}

	models, err := c.Logic()

	if err != nil {
		return nil, err
	}

	items, err := newView(models)
	if err != nil {
		return nil, appErrors.NewApplicationError(err).AddError("ListItems.Paginate.Handle", nil)
	}

	return items, nil
}

func New(model Model, auth auth.Authentication) pkg.Job[Model, []View, []QueryVariable] {
	return Main{model: model, auth: auth}
}
