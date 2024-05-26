package createAndDiff

import (
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/lib/sdk"
	"errors"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type PreViewStructure struct {
	Name    string
	ID      string
	ShortID string
}

type ListOrMap struct {
	ID            string
	Name          string
	ShortID       string
	StructureType string
}

type PreViewModel struct {
	ID    string
	Name  string
	State string
	Maps  []PreViewStructure
	Lists []PreViewStructure
}

type Diff struct {
	Lists []declarations.List
	Maps  []declarations.Map
}

type LogicModel struct {
	Metadata   PreViewModel
	Diff       Diff
	Structures []ListOrMap
}

type MetadataModel struct {
	ID          string
	Name        string
	State       string
	UserID      string `gorm:"column:user_id"`
	Map         string `gorm:"column:map_name"`
	MapID       string `gorm:"column:map_id"`
	MapShortID  string `gorm:"column:map_short_id"`
	List        string `gorm:"column:list_name"`
	ListID      string `gorm:"column:list_id"`
	ListShortID string `gorm:"column:list_short_id"`
	MapLocale   string `gorm:"column:map_locale"`
	ListLocale  string `gorm:"column:list_locale"`
}

type Structure struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type Model struct {
	Structures []Structure `json:"structures"`
	ID         string      `json:"projectName"`
}

func NewModel(ID string, structures []Structure) Model {
	return Model{
		ID:         ID,
		Structures: structures,
	}
}

func (a Model) Validate() map[string]string {
	v := map[string]interface{}{
		"id":         a.ID,
		"structures": nil,
	}

	if err := validation.Validate(v,
		validation.Map(
			validation.Key("id", validation.Required, validation.RuneLength(27, 27)),
			validation.Key("structures", validation.By(func(value interface{}) error {
				items := a.Structures

				for _, item := range items {
					if len(item.Name) < 1 && len(item.Name) > 200 {
						return errors.New("Structure name invalid. Name is required and must have length between 1 and 200 characters")
					}

					if item.Type != "map" && item.Type != "list" {
						return errors.New("Structure type invalid. Structure type can be either 'map' or 'list'")
					}
				}

				return nil
			})),
		),
	); err != nil {
		return sdk.ErrorToResponseError(err)
	}

	return nil
}
