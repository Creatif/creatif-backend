package deleteListItemByID

import (
	"creatif/pkg/lib/sdk"
	"errors"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Model struct {
	Name        string
	ID          string
	ShortID     string
	ItemID      string
	ItemShortID string
	ProjectID   string
}

func NewModel(projectId, name, id, shortID, itemId, itemShortID string) Model {
	return Model{
		Name:        name,
		ID:          id,
		ShortID:     shortID,
		ItemID:      itemId,
		ItemShortID: itemShortID,
		ProjectID:   projectId,
	}
}

func (a Model) Validate() map[string]string {
	v := map[string]interface{}{
		"name":         a.Name,
		"id":           a.ID,
		"idExists":     nil,
		"itemIDExists": nil,
		"projectID":    a.ProjectID,
	}

	if err := validation.Validate(v,
		validation.Map(
			validation.Key("name", validation.When(a.Name != "", validation.RuneLength(1, 200))),
			validation.Key("id", validation.When(a.ID != "", validation.RuneLength(26, 26))),
			validation.Key("idExists", validation.By(func(value interface{}) error {
				name := a.Name
				shortId := a.ShortID
				id := a.ID

				if id != "" && len(id) != 26 {
					return errors.New("ItemID must have 26 characters")
				}

				if name == "" && shortId == "" && id == "" {
					return errors.New("At least one of 'id', 'name' or 'shortID' must be supplied in order to identify this list.")
				}
				return nil
			})),
			validation.Key("itemIDExists", validation.By(func(value interface{}) error {
				shortId := a.ItemShortID
				id := a.ItemID

				if id != "" && len(id) != 26 {
					return errors.New("ItemID must have 26 characters")
				}

				if shortId == "" && id == "" {
					return errors.New("At least one of 'id', or 'shortID' must be supplied in order to identify this variable.")
				}
				return nil
			})),

			validation.Key("projectID", validation.Required, validation.RuneLength(26, 26)),
		),
	); err != nil {
		return sdk.ErrorToResponseError(err)
	}

	return nil
}
