package getProjectMetadata

import (
	"creatif/pkg/app/auth"
	"creatif/pkg/app/domain/app"
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/app/services/locales"
	pkg "creatif/pkg/lib"
	"creatif/pkg/lib/appErrors"
	"creatif/pkg/lib/logger"
	"creatif/pkg/lib/storage"
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
	var logicModels []LogicModel
	res := storage.Gorm().Raw(fmt.Sprintf(`
SELECT 
p.id,
p.name,
p.state,
p.user_id,
v.name AS variable_name,
m.name AS map_name,
l.name AS list_name,
v.locale_id AS variable_locale,
m.locale_id AS map_locale
FROM %s AS p
LEFT JOIN %s AS v ON p.id = ? AND p.user_id = ? AND v.project_id = p.id AND v.project_id = ?
LEFT JOIN %s AS m ON m.project_id = p.id AND m.project_id = ?
LEFT JOIN %s AS l ON l.project_id = p.id AND l.project_id = ?
WHERE p.id = ? AND p.user_id = ?
`,
		(app.Project{}).TableName(),
		(declarations.Variable{}).TableName(),
		(declarations.Map{}).TableName(),
		(declarations.List{}).TableName(),
	),
		c.auth.User().ProjectID, c.auth.User().ID, c.auth.User().ProjectID,
		c.auth.User().ProjectID,
		c.auth.User().ProjectID,
		c.auth.User().ProjectID, c.auth.User().ID,
	).Scan(&logicModels)

	if res.Error != nil {
		return PreViewModel{}, appErrors.NewNotFoundError(res.Error)
	}

	preViewModel := PreViewModel{
		ID:        logicModels[0].ID,
		Name:      logicModels[0].Name,
		State:     logicModels[0].State,
		UserID:    logicModels[0].UserID,
		Variables: make(map[string][]string),
		Maps:      make(map[string][]string),
		Lists:     make([]string, 0),
	}

	if len(logicModels) == 1 && logicModels[0].VariableName == "" && logicModels[0].Map == "" && logicModels[0].List == "" {
		return preViewModel, nil
	}

	for _, v := range logicModels {
		variableLocale, _ := locales.GetAlphaWithID(v.VariableLocale)
		mapLocale, _ := locales.GetAlphaWithID(v.MapLocale)

		if _, ok := preViewModel.Variables[v.VariableLocale]; !ok && variableLocale != "" {
			preViewModel.Variables[variableLocale] = make([]string, 0)
		}

		if _, ok := preViewModel.Maps[mapLocale]; !ok && mapLocale != "" {
			preViewModel.Maps[mapLocale] = make([]string, 0)
		}

		if v.VariableName != "" {
			preViewModel.Variables[variableLocale] = append(preViewModel.Variables[variableLocale], v.VariableName)
		}

		if v.Map != "" {
			preViewModel.Maps[mapLocale] = append(preViewModel.Maps[mapLocale], v.Map)
		}

		if v.List != "" {
			preViewModel.Lists = append(preViewModel.Lists, v.List)
		}
	}

	if len(preViewModel.Variables) == 0 {
		preViewModel.Variables = nil
	}

	if len(preViewModel.Maps) == 0 {
		preViewModel.Maps = nil
	}

	if len(preViewModel.Lists) == 0 {
		preViewModel.Lists = nil
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
