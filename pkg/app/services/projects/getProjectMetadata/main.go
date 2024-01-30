package getProjectMetadata

import (
	"creatif/pkg/app/auth"
	"creatif/pkg/app/services/locales"
	pkg "creatif/pkg/lib"
	"creatif/pkg/lib/appErrors"
	"creatif/pkg/lib/logger"
	"fmt"
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
	fmt.Println("projectid: ", c.auth.User().ProjectID, "userid: ", c.auth.User().ID)
	logicModels, err := getVariablesMetadata(c.auth.User().ProjectID, c.auth.User().ID)
	if err != nil {
		return PreViewModel{}, err
	}

	preViewModel := PreViewModel{
		ID:        logicModels[0].ID,
		Name:      logicModels[0].Name,
		State:     logicModels[0].State,
		UserID:    logicModels[0].UserID,
		Variables: make(map[string][]PreViewStructure),
		Maps:      make([]PreViewStructure, 0),
		Lists:     make([]PreViewStructure, 0),
	}

	if len(logicModels) == 1 && logicModels[0].VariableName == "" && logicModels[0].Map == "" && logicModels[0].List == "" {
		return preViewModel, nil
	}

	for _, v := range logicModels {
		variableLocale, _ := locales.GetAlphaWithID(v.VariableLocale)

		if _, ok := preViewModel.Variables[v.VariableLocale]; !ok && variableLocale != "" {
			preViewModel.Variables[variableLocale] = make([]PreViewStructure, 0)
		}

		if v.VariableName != "" {
			preViewModel.Variables[variableLocale] = append(preViewModel.Variables[variableLocale], PreViewStructure{
				Name:    v.VariableName,
				ID:      v.VariableID,
				ShortID: v.VariableShortID,
			})
		}

		if v.Map != "" {
			preViewModel.Maps = append(preViewModel.Maps, PreViewStructure{
				Name:    v.Map,
				ID:      v.MapID,
				ShortID: v.MapShortID,
			})
		}

		if v.List != "" {
			preViewModel.Lists = append(preViewModel.Lists, PreViewStructure{
				Name:    v.List,
				ID:      v.ListID,
				ShortID: v.ListShortID,
			})
		}
	}

	if len(preViewModel.Variables) == 0 {
		preViewModel.Variables = nil
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
