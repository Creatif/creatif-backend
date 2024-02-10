package getProjectMetadata

import (
	"creatif/pkg/app/auth"
	pkg "creatif/pkg/lib"
	"creatif/pkg/lib/appErrors"
	"creatif/pkg/lib/logger"
	"creatif/pkg/lib/sdk"
)

type Main struct {
	logBuilder logger.LogBuilder
	auth       auth.Authentication
}

func (c Main) Validate() error {
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

func (c Main) Logic() (PreViewModel, error) {
	logicModels, err := getVariablesMetadata(c.auth.User().ProjectID, c.auth.User().ID)
	if err != nil {
		return PreViewModel{}, err
	}

	preViewModel := PreViewModel{
		ID:     logicModels[0].ID,
		Name:   logicModels[0].Name,
		State:  logicModels[0].State,
		UserID: logicModels[0].UserID,
		Maps:   make([]PreViewStructure, 0),
		Lists:  make([]PreViewStructure, 0),
	}

	if len(logicModels) == 1 && logicModels[0].Map == "" && logicModels[0].List == "" {
		return preViewModel, nil
	}

	for _, v := range logicModels {
		if v.Map != "" {
			// if an entry exists, skip
			for _, l := range preViewModel.Lists {
				if l.Name == v.Map {
					continue
				}
			}

			preViewModel.Maps = append(preViewModel.Maps, PreViewStructure{
				Name:    v.Map,
				ID:      v.MapID,
				ShortID: v.MapShortID,
			})
		}

		if v.List != "" {
			if !sdk.IncludesFn(preViewModel.Lists, func(item PreViewStructure) bool {
				return item.Name == v.List
			}) {
				preViewModel.Lists = append(preViewModel.Lists, PreViewStructure{
					Name:    v.List,
					ID:      v.ListID,
					ShortID: v.ListShortID,
				})
			}
		}
	}

	return preViewModel, nil
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

func New(auth auth.Authentication, builder logger.LogBuilder) pkg.Job[interface{}, View, PreViewModel] {
	builder.Add("projectService", "Get project")
	return Main{logBuilder: builder, auth: auth}
}
