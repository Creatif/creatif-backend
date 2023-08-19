package create

import (
	"creatif/pkg/lib/sdk"
	"errors"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type CreateMapModel struct {
	Names []string `json:"name"`
}

func NewCreateMapModel(names []string) CreateMapModel {
	return CreateMapModel{
		Names: names,
	}
}

func (a *CreateMapModel) Validate() map[string]string {
	v := map[string]interface{}{
		"names":    a.Names,
		"validNum": a.Names,
	}

	if err := validation.Validate(v,
		validation.Map(
			validation.Key("groups", validation.When(len(a.Names) != 0, validation.Each(is.UUID))),
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
	Names []string `json:"names"`
}

func newView(names []string) View {
	return View{
		Names: names,
	}
}
