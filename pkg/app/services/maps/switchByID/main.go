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

func (c Main) Logic() (float64, error) {
	var chosenMap declarations.Map
	res := storage.Gorm().Raw(fmt.Sprintf("SELECT id FROM %s WHERE project_id = ? AND (id = ? OR name = ? OR short_id = ?)", (declarations.Map{}).TableName()), c.model.ProjectID, c.model.Name, c.model.Name, c.model.Name).Scan(&chosenMap)
	if res.Error != nil {
		c.logBuilder.Add("switchByID.Map", fmt.Sprintf("Update: Invalid query: %s", res.Error.Error()))

		return 0, appErrors.NewApplicationError(res.Error)
	}

	if res.RowsAffected == 0 {
		c.logBuilder.Add("switchByID", fmt.Sprintf("Update: Map variable not found: %s", "Rows affected: 0"))
		return 0, appErrors.NewNotFoundError(errors.New("Could not find map"))
	}

	sdVariables, err := getSourceDestinationVariables(c.model.ProjectID, c.model.Name, c.model.Source, c.model.Destination)
	if err != nil {
		return 0, err
	}

	idxRange, err := getHighestLowestIndex(chosenMap.ID, c.model.ProjectID)
	if err != nil && err.Error() == "not_found" {
		return 0, appErrors.NewNotFoundError(err)
	} else if err != nil {
		return 0, appErrors.NewApplicationError(err)
	}

	if idxRange.Highest == sdVariables.destination.Index {
		return idxRange.Highest + 1, updateWithCustomIndex(idxRange.Highest+1, sdVariables.source.ID, chosenMap.ID)
	}

	if idxRange.Lowest == sdVariables.destination.Index {
		return idxRange.Highest + 1, updateWithCustomIndex(idxRange.Lowest-1, sdVariables.source.ID, chosenMap.ID)
	}

	upperIndexOperator := "<"
	if sdVariables.source.Index < sdVariables.destination.Index {
		upperIndexOperator = ">"
	}
	var upperIndexes []float64
	res = storage.Gorm().Raw(fmt.Sprintf(`
SELECT index 
FROM declarations.map_variables 
WHERE map_id = ? AND index %s (SELECT index FROM declarations.map_variables WHERE id = ?) ORDER BY index DESC LIMIT 1`, upperIndexOperator), chosenMap.ID, c.model.Destination).Scan(&upperIndexes)

	if res.Error != nil {
		return 0, appErrors.NewValidationError(map[string]string{
			"invalidSourceDestination": "Incomplete declaration map",
		})
	}

	if res.RowsAffected == 0 {
		return 0, appErrors.NewValidationError(map[string]string{
			"invalidSourceDestination": "Incomplete declaration map",
		})
	}

	var realIndex float64
	if res.RowsAffected != 0 {
		realIndex = upperIndexes[0]
	}

	res = storage.Gorm().Exec(fmt.Sprintf(`
UPDATE %s
SET index = round(((coalesce(?, 1000) + (SELECT index FROM declarations.map_variables WHERE id = ?)) / 2)::numeric, 10)  WHERE id = ? AND map_id = ?
`,
		(declarations.MapVariable{}).TableName(),
	), realIndex, c.model.Destination, c.model.Source, chosenMap.ID)

	if res.Error != nil {
		c.logBuilder.Add("switchByID", fmt.Sprintf("Update: Invalid query: %s", res.Error.Error()))

		return 0, appErrors.NewApplicationError(res.Error)
	}

	if res.RowsAffected == 0 {
		c.logBuilder.Add("switchByID", fmt.Sprintf("Update: Not found: %s", "Rows affected: 0"))
		return 0, appErrors.NewNotFoundError(errors.New("Could not switch map variables."))
	}

	var emptyVariableWithIndex declarations.MapVariable
	res = storage.Gorm().Where("id = ?", c.model.Source).Select("index").First(&emptyVariableWithIndex)
	if res.Error != nil {
		return 0, appErrors.NewNotFoundError(errors.New("Could not switch map variables."))
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

func New(model Model, auth auth.Authentication, logBuilder logger.LogBuilder) pkg.Job[Model, float64, float64] {
	logBuilder.Add("switchByID", "Created")
	return Main{model: model, logBuilder: logBuilder, auth: auth}
}
