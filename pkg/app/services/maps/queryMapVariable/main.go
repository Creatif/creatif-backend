package queryMapVariable

import (
	"creatif/pkg/app/auth"
	"creatif/pkg/app/services/shared/connections"
	pkg "creatif/pkg/lib"
	"creatif/pkg/lib/appErrors"
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
	variable, err := getVariable(c.model.ProjectID, c.model.Name, c.model.ItemID)
	if err != nil {
		return LogicModel{}, err
	}

	conns, err := getChildConnectionFromParent(variable.ID)
	if err != nil {
		return LogicModel{}, appErrors.NewApplicationError(err)
	}

	// replace the jsonb connections with actual variables.
	// this directly modifies the jsonb array and replaces the variable.Value.
	replacedValue, err := connections.ReplaceJson(conns, variable.Value, c.model.ConnectionReplaceMethod)
	if err != nil {
		return LogicModel{}, appErrors.NewApplicationError(err)
	}
	variable.Value = replacedValue

	structures, err := getViewStructuresByVariableFromConnections(c.model.ItemID)
	if err != nil {
		return LogicModel{}, appErrors.NewApplicationError(err)

	}

	return LogicModel{
		Variable:                  variable,
		ChildConnectionStructures: structures,
		Connections:               conns,
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
