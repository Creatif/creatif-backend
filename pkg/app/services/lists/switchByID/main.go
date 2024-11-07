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
	var chosenList declarations.List
	res := storage.Gorm().Raw(fmt.Sprintf("SELECT id FROM %s WHERE project_id = ? AND (id = ? OR name = ? OR short_id = ?)", (declarations.List{}).TableName()), c.model.ProjectID, c.model.Name, c.model.Name, c.model.Name).Scan(&chosenList)
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

	idxRange, err := getHighestLowestIndex(chosenList.ID, c.model.ProjectID)
	if err != nil && err.Error() == "not_found" {
		return 0, appErrors.NewNotFoundError(err)
	} else if err != nil {
		return 0, appErrors.NewApplicationError(err)
	}

	if c.model.OrderDirection == "asc" && sdVariables.destination.Index == idxRange.Lowest {
		return idxRange.Lowest - 1024, updateWithCustomIndex(idxRange.Lowest-1024, sdVariables.source.ID, chosenList.ID)
	} else if c.model.OrderDirection == "asc" && sdVariables.destination.Index == idxRange.Highest {
		return idxRange.Highest + 1024, updateWithCustomIndex(idxRange.Highest+1024, sdVariables.source.ID, chosenList.ID)
	} else if c.model.OrderDirection == "desc" && sdVariables.destination.Index == idxRange.Highest {
		return idxRange.Highest + 1024, updateWithCustomIndex(idxRange.Highest+1024, sdVariables.source.ID, chosenList.ID)
	} else if c.model.OrderDirection == "desc" && sdVariables.destination.Index == idxRange.Lowest {
		return idxRange.Lowest - 1024, updateWithCustomIndex(idxRange.Lowest-1024, sdVariables.source.ID, chosenList.ID)
	}
	var upperIndex float64
	if c.model.OrderDirection == "desc" {
		idx, err := getIndexDesc(chosenList.ID, sdVariables.destination.Index)
		if err != nil {
			return 0, appErrors.NewApplicationError(err)
		}

		upperIndex = idx
	} else if c.model.OrderDirection == "asc" {
		idx, err := getIndexAsc(chosenList.ID, sdVariables.destination.Index)
		if err != nil {
			return 0, appErrors.NewApplicationError(err)
		}

		upperIndex = idx
	}

	newIndex := (sdVariables.destination.Index + upperIndex) / 2
	if err := updateDestinationIndex(chosenList.ID, sdVariables.source.ID, newIndex); err != nil {
		return 0, appErrors.NewApplicationError(err)
	}

	return newIndex, nil
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
