package switchByID

import (
	"creatif/pkg/app/auth"
	"creatif/pkg/app/domain/declarations"
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

func (c Main) Logic() (float64, error) {
	var list declarations.List
	res := storage.Gorm().Raw(fmt.Sprintf("SELECT id FROM %s WHERE project_id = ? AND (id = ? OR name = ? OR short_id = ?)", (declarations.List{}).TableName()), c.model.ProjectID, c.model.Name, c.model.Name, c.model.Name).Scan(&list)
	if res.Error != nil {
		return 0, appErrors.NewApplicationError(res.Error)
	}

	if res.RowsAffected == 0 {
		return 0, appErrors.NewNotFoundError(errors.New("Could not find list"))
	}

	sdVariables, err := getSourceDestinationVariables(c.model.ProjectID, c.model.Name, c.model.Source, c.model.Destination)
	if err != nil {
		return 0, err
	}

	idxRange, err := getHighestLowestIndex(list.ID, c.model.ProjectID)
	if err != nil && err.Error() == "not_found" {
		return 0, appErrors.NewNotFoundError(err)
	} else if err != nil {
		return 0, appErrors.NewApplicationError(err)
	}

	if idxRange.Highest == sdVariables.destination.Index {
		return idxRange.Highest + 1, updateWithCustomIndex(idxRange.Highest+1, sdVariables.source.ID, list.ID)
	}

	if idxRange.Lowest == sdVariables.destination.Index {
		return idxRange.Highest + 1, updateWithCustomIndex(idxRange.Lowest-1, sdVariables.source.ID, list.ID)
	}

	upperIndexOperator := "<"
	if sdVariables.source.Index < sdVariables.destination.Index {
		upperIndexOperator = ">"
	}
	var upperIndexes []float64
	res = storage.Gorm().Raw(fmt.Sprintf(`
SELECT index 
FROM declarations.list_variables 
WHERE list_id = ? AND index %s (SELECT index FROM declarations.list_variables WHERE id = ?) ORDER BY index DESC LIMIT 1`, upperIndexOperator), list.ID, c.model.Destination).Scan(&upperIndexes)

	if res.Error != nil {
		return 0, appErrors.NewValidationError(map[string]string{
			"invalidSourceDestination": "Incomplete declaration list",
		})
	}

	if res.RowsAffected == 0 {
		return 0, appErrors.NewValidationError(map[string]string{
			"invalidSourceDestination": "Incomplete declaration list",
		})
	}

	var realIndex float64
	if res.RowsAffected != 0 {
		realIndex = upperIndexes[0]
	}

	res = storage.Gorm().Exec(fmt.Sprintf(`
UPDATE %s
SET index = round(((coalesce(?, 1000) + (SELECT index FROM declarations.list_variables WHERE id = ?)) / 2)::numeric, 10)  WHERE id = ? AND list_id = ?
`,
		(declarations.ListVariable{}).TableName(),
	), realIndex, c.model.Destination, c.model.Source, list.ID)

	if res.Error != nil {
		return 0, appErrors.NewApplicationError(res.Error)
	}

	if res.RowsAffected == 0 {
		return 0, appErrors.NewNotFoundError(errors.New("Could not switch list variables."))
	}

	var emptyVariableWithIndex declarations.ListVariable
	res = storage.Gorm().Where("id = ?", c.model.Source).Select("index").First(&emptyVariableWithIndex)
	if res.Error != nil {
		return 0, appErrors.NewNotFoundError(errors.New("Could not switch list variables."))
	}

	return emptyVariableWithIndex.Index, nil
}

func (c Main) Handle() (float64, error) {
	if err := c.Validate(); err != nil {
		return 0, err
	}

	if err := c.Authenticate(); err != nil {
		return 0, err
	}

	if err := c.Authorize(); err != nil {
		return 0, err
	}

	changedIndex, err := c.Logic()

	if err != nil {
		return 0, err
	}

	return changedIndex, nil
}

func New(model Model, auth auth.Authentication) pkg.Job[Model, float64, float64] {
	return Main{model: model, auth: auth}
}
