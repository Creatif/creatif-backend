package create

import (
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/lib/sdk"
	"creatif/pkg/lib/storage"
	"errors"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"gorm.io/gorm"
)

type CreateMapModel struct {
	Nodes []string `json:"nodes"`
	Name  string   `json:"name"`
}

func NewCreateMapModel(name string, nodes []string) CreateMapModel {
	return CreateMapModel{
		Nodes: nodes,
		Name:  name,
	}
}

func (a *CreateMapModel) Validate() map[string]string {
	v := map[string]interface{}{
		"uniqueName": a.Name,
		"nodes":      a.Nodes,
		"validNum":   a.Nodes,
	}

	if err := validation.Validate(v,
		validation.Map(
			validation.Key("uniqueName", validation.By(func(value interface{}) error {
				name := value.(string)

				maps := make([]declarations.Map, 0)
				if err := storage.GetBy((&declarations.Map{}).TableName(), "name", name, &maps); err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
					return errors.New("Map with this name already exists")
				}

				if len(maps) != 0 {
					return errors.New("Map with this name already exists")
				}

				return nil
			})),
			validation.Key("nodes", validation.When(len(a.Nodes) != 0, validation.Each(is.UUID))),
			validation.Key("validNum", validation.By(func(value interface{}) error {
				names := value.([]string)
				if len(names) > 100 {
					return errors.New("Number of combined node cannot be higher than 100")
				}

				return nil
			})),
		),
	); err != nil {
		return sdk.ErrorToResponseError(err)
	}

	return nil
}

type View struct {
	Name  string   `json:"name"`
	Nodes []string `json:"nodes"`
}

func newView(name string, nodes []string) View {
	return View{
		Name:  name,
		Nodes: nodes,
	}
}
