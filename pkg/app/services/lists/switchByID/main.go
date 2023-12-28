package switchByID

import (
	"creatif/pkg/app/auth"
	"creatif/pkg/app/domain/declarations"
	pkg "creatif/pkg/lib"
	"creatif/pkg/lib/appErrors"
	"creatif/pkg/lib/logger"
	"creatif/pkg/lib/storage"
	"errors"
	"fmt"
)

type Main struct {
	model      Model
	logBuilder logger.LogBuilder
	auth       auth.Authentication
}

// 1. Check if source and destination exist here
func (c Main) Validate() error {
	c.logBuilder.Add("switchByID", "Validating...")
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

// 1. Get the list item behind the destination sorted by index
// 2. Update the source item with the new index => [destination-1].index + destination.index / 2
func (c Main) Logic() (interface{}, error) {
	/*	source, destination, err := tryUpdates(c.model.ProjectID, c.model.Name, c.model.ID, c.model.ShortID, c.model.Source, c.model.Destination, 0, 10)
		if err != nil {
			c.logBuilder.Add("switchByID", err.Error())
			return LogicResult{}, appErrors.NewDatabaseError(err)
		}*/

	var list declarations.List
	res := storage.Gorm().Raw(fmt.Sprintf("SELECT id FROM %s WHERE project_id = ? AND (id = ? OR name = ? OR short_id = ?)", (declarations.List{}).TableName()), c.model.ProjectID, c.model.Name, c.model.Name, c.model.Name).Scan(&list)
	if res.Error != nil {
		c.logBuilder.Add("switchByID", fmt.Sprintf("Update: Invalid query: %s", res.Error.Error()))

		return nil, appErrors.NewApplicationError(res.Error)
	}

	if res.RowsAffected == 0 {
		c.logBuilder.Add("switchByID", fmt.Sprintf("Update: List not found: %s", "Rows affected: 0"))
		return nil, appErrors.NewNotFoundError(errors.New("Could not find list"))
	}

	sdVariables, err := getSourceDestinationVariables(c.model.ProjectID, c.model.Name, c.model.Source, c.model.Destination)
	if err != nil {
		return nil, err
	}

	idxRange, err := getHighestLowestIndex(list.ID, c.model.ProjectID)
	if err != nil && err.Error() == "not_found" {
		return nil, appErrors.NewNotFoundError(err)
	} else if err != nil {
		return nil, appErrors.NewApplicationError(err)
	}

	if idxRange.Highest == sdVariables.destination.Index {
		fmt.Println("lowest to highest")
		return nil, updateWithCustomIndex(idxRange.Highest+1, sdVariables.source.ID, list.ID)
	}

	if idxRange.Lowest == sdVariables.destination.Index {
		fmt.Println("highest to lowest")
		return nil, updateWithCustomIndex(idxRange.Lowest-1, sdVariables.source.ID, list.ID)
	}

	upperIndexOperator := "<"
	if sdVariables.source.Index < sdVariables.destination.Index {
		upperIndexOperator = ">"
	}
	var upperIndexes []float64
	res = storage.Gorm().Raw(fmt.Sprintf(`
SELECT index 
FROM declarations.list_variables 
WHERE list_id = ? AND index %s (SELECT index FROM declarations.list_variables WHERE id = ?) ORDER BY index DESC LIMIT 2`, upperIndexOperator), list.ID, c.model.Destination).Scan(&upperIndexes)

	if res.Error != nil {
		return nil, appErrors.NewValidationError(map[string]string{
			"invalidSourceDestination": "Incomplete declaration list",
		})
	}

	if res.RowsAffected == 0 {
		return nil, appErrors.NewValidationError(map[string]string{
			"invalidSourceDestination": "Incomplete declaration list",
		})
	}

	var realIndex float64
	if res.RowsAffected == 1 {
		realIndex = upperIndexes[0]
	} else {
		realIndex = upperIndexes[1]
	}

	res = storage.Gorm().Exec(fmt.Sprintf(`
UPDATE %s
SET index = round(((coalesce(?, 1000) + (SELECT index FROM declarations.list_variables WHERE id = ?)) / 2)::numeric, 10)  WHERE id = ? AND list_id = ?
`,
		(declarations.ListVariable{}).TableName(),
	), realIndex, c.model.Destination, c.model.Source, list.ID)

	if res.Error != nil {
		c.logBuilder.Add("switchByID", fmt.Sprintf("Update: Invalid query: %s", res.Error.Error()))

		return nil, appErrors.NewApplicationError(res.Error)
	}

	if res.RowsAffected == 0 {
		c.logBuilder.Add("switchByID", fmt.Sprintf("Update: Not found: %s", "Rows affected: 0"))
		return nil, appErrors.NewNotFoundError(errors.New("Could not switch list variables."))
	}

	return nil, nil
}

func (c Main) Handle() (interface{}, error) {
	if err := c.Validate(); err != nil {
		return View{}, err
	}

	if err := c.Authenticate(); err != nil {
		return View{}, err
	}

	if err := c.Authorize(); err != nil {
		return View{}, err
	}

	_, err := c.Logic()

	if err != nil {
		return View{}, err
	}

	return nil, nil
}

func New(model Model, auth auth.Authentication, logBuilder logger.LogBuilder) pkg.Job[Model, interface{}, interface{}] {
	logBuilder.Add("switchByID", "Created")
	return Main{model: model, logBuilder: logBuilder, auth: auth}
}
