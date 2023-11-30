package getListGroups

import (
	"creatif/pkg/lib/sdk"
	"errors"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type LogicModel struct {
	Group string
}

type Model struct {
	ID      string
	Name    string
	ShortID string
}

func NewModel(listID, name, shortId string) Model {
	return Model{
		ID:      listID,
		Name:    name,
		ShortID: shortId,
	}
}

type View struct {
}

func newView(model []LogicModel) []string {
	return sdk.Map(model, func(idx int, value LogicModel) string {
		return value.Group
	})
}

func (a *Model) Validate() map[string]string {
	v := map[string]interface{}{
		"listID":   a.ID,
		"idExists": nil,
	}

	if err := validation.Validate(v,
		validation.Map(
			validation.Key("listID", validation.When(a.ID != "", validation.RuneLength(26, 26))),
			validation.Key("idExists", validation.By(func(value interface{}) error {
				name := a.Name
				shortId := a.ShortID
				id := a.ID

				if id != "" && len(id) != 26 {
					return errors.New("ID must have 26 characters")
				}

				if name == "" && shortId == "" && id == "" {
					return errors.New("At least one of 'id', 'name' or 'shortID' must be supplied in order to identify this list.")
				}
				return nil
			})),
		),
	); err != nil {
		return sdk.ErrorToResponseError(err)
	}

	return nil
}