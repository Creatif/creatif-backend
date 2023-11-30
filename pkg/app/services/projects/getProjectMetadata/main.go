package getProjectMetadata

import (
	"creatif/pkg/app/auth"
	"creatif/pkg/app/domain/app"
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

func (c Main) Logic() (LogicModel, error) {
	var logicModel LogicModel
	res := storage.Gorm().Raw(fmt.Sprintf(`
SELECT 
p.id,
p.name,
p.state,
p.user_id,
v.name = ARRAY(SELECT name FROM declarations.variables AS n WHERE n.project_id = p.id AND p.id = ?) AS variable_names,
m.name = ARRAY(SELECT name FROM declarations.maps AS n WHERE n.project_id = p.id AND p.id = ?) AS map_names,
l.name = ARRAY(SELECT name FROM declarations.lists AS n WHERE n.project_id = p.id AND p.id = ?) AS list_names
FROM %s AS p
WHERE p.id = ? AND p.user_id = ?
`,
		(app.Project{}).TableName(),
	), c.auth.User().ProjectID, c.auth.User().ID).Scan(&logicModel)

	if res.Error != nil {
		return LogicModel{}, appErrors.NewNotFoundError(res.Error)
	}

	return LogicModel{}, nil
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

func New(auth auth.Authentication, builder logger.LogBuilder) pkg.Job[interface{}, View, LogicModel] {
	builder.Add("projectService", "Get project")
	return Main{logBuilder: builder, auth: auth}
}
