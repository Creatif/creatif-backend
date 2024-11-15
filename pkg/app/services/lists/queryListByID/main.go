package queryListByID

import (
	"creatif/pkg/app/auth"
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/app/services/shared"
	"creatif/pkg/app/services/shared/connections"
	pkg "creatif/pkg/lib"
	"creatif/pkg/lib/appErrors"
	"creatif/pkg/lib/storage"
	"errors"
	"fmt"
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
	sql := fmt.Sprintf(`
			SELECT 
			    lv.id, 
			    lv.name, 
			    lv.index, 
			    lv.behaviour, 
			    lv.short_id, 
			    lv.metadata, 
			       ARRAY((SELECT g.name FROM %s AS vg INNER JOIN %s AS g ON vg.variable_id = lv.id AND g.id = ANY(vg.groups))) AS groups,
				lv.value, 
			    lv.created_at, 
			    lv.updated_at, 
			    lv.locale_id
			FROM %s AS l INNER JOIN %s AS lv
			ON l.project_id = ? AND l.id = ? AND lv.list_id = l.id AND lv.id = ?`,
		(declarations.VariableGroup{}).TableName(), (declarations.Group{}).TableName(), (declarations.List{}).TableName(), (declarations.ListVariable{}).TableName())

	var variable QueryVariable
	res := storage.Gorm().
		Raw(sql, c.model.ProjectID, c.model.Name, c.model.ItemID).
		Scan(&variable)

	if res.Error != nil {
		return LogicModel{}, appErrors.NewDatabaseError(res.Error).AddError("queryMapVariable.Logic", nil)
	}

	if res.RowsAffected == 0 {
		return LogicModel{}, appErrors.NewNotFoundError(errors.New("No rows found")).AddError("queryMapVariable.Logic", nil)
	}

	replacedValue, err := connections.ReplaceJson(variable.Value, variable.ID)
	if err != nil {
		return LogicModel{}, appErrors.NewApplicationError(err)
	}
	variable.Value = replacedValue

	references, err := shared.QueryReferences(variable.ID, c.model.ProjectID)
	if err != nil {
		return LogicModel{}, err
	}

	return LogicModel{
		Variable:  variable,
		Reference: references,
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
