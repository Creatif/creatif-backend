package createAndDiff

import "C"
import (
	"creatif/pkg/app/auth"
	"creatif/pkg/app/domain/declarations"
	pkg "creatif/pkg/lib"
	"creatif/pkg/lib/appErrors"
	"creatif/pkg/lib/storage"
	"gorm.io/gorm"
)

type Results struct {
	Errors []error
}

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
	/**
	1. Get all lists for this project
	2. Get all maps for this project
	3. Iterate through list structures from model, if any not exists, create it
	    3.1 If the structure exists in db but not in config, save it
	4. Iterate through map structure from model, if any not exists, create it
	    4.1 If structure exists in db but not in config, save it
	5. Create project metadata
	6. Return response with project metadata and saved diff with config and db
	*/

	project, err := getProject(c.model.ID)
	if err != nil {
		return LogicModel{}, err
	}

	var lists []declarations.List
	if err := getStructures(project.ID, (declarations.List{}).TableName(), &lists); err != nil {
		return LogicModel{}, err
	}

	var maps []declarations.Map
	if err := getStructures(project.ID, (declarations.Map{}).TableName(), &maps); err != nil {
		return LogicModel{}, err
	}

	listsToCreate, mapsToCreate, listsNotInConfig, mapsNotInConfig := processListsAndMaps(project.ID, c.model.Structures, lists, maps)

	if len(listsToCreate) != 0 || len(mapsToCreate) != 0 {
		if transactionErr := storage.Transaction(func(tx *gorm.DB) error {
			if len(listsToCreate) != 0 {
				if res := tx.Create(&listsToCreate); res.Error != nil {
					return res.Error
				}
			}

			if len(mapsToCreate) != 0 {
				if res := tx.Create(&mapsToCreate); res.Error != nil {
					return res.Error
				}
			}

			return nil
		}); transactionErr != nil {
			return LogicModel{}, appErrors.NewDatabaseError(transactionErr)
		}
	}

	metadata, err := getProjectMetadata(project.ID)
	if err != nil {
		return LogicModel{}, err
	}

	processedMetadata := processMetadata(metadata)

	var allLists []declarations.List
	if err := getStructures(project.ID, (declarations.List{}).TableName(), &allLists); err != nil {
		return LogicModel{}, err
	}

	var allMaps []declarations.Map
	if err := getStructures(project.ID, (declarations.Map{}).TableName(), &allMaps); err != nil {
		return LogicModel{}, err
	}

	structures := make([]ListOrMap, 0)
	for _, item := range allLists {
		structures = append(structures, ListOrMap{
			ID:            item.ID,
			Name:          item.Name,
			ShortID:       item.ShortID,
			StructureType: "list",
		})
	}

	for _, item := range allMaps {
		structures = append(structures, ListOrMap{
			ID:            item.ID,
			Name:          item.Name,
			ShortID:       item.ShortID,
			StructureType: "map",
		})
	}

	return LogicModel{
		Metadata: processedMetadata,
		Diff: Diff{
			Lists: listsNotInConfig,
			Maps:  mapsNotInConfig,
		},
		Structures: structures,
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
