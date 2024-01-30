package getProjectMetadata

import "creatif/pkg/lib/sdk"

type LogicModel struct {
	ID              string
	Name            string
	State           string
	UserID          string `gorm:"column:user_id"`
	Map             string `gorm:"column:map_name"`
	MapID           string `gorm:"column:map_id"`
	MapShortID      string `gorm:"column:map_short_id"`
	List            string `gorm:"column:list_name"`
	ListID          string `gorm:"column:list_id"`
	ListShortID     string `gorm:"column:list_short_id"`
	VariableName    string `gorm:"column:variable_name"`
	VariableID      string `gorm:"column:variable_id"`
	VariableShortID string `gorm:"column:variable_short_id"`
	VariableLocale  string `gorm:"column:variable_locale"`
	MapLocale       string `gorm:"column:map_locale"`
	ListLocale      string `gorm:"column:list_locale"`
}

type StructureView struct {
	Name    string `json:"name"`
	ID      string `json:"id"`
	ShortID string `json:"shortId"`
}

type PreViewStructure struct {
	Name    string
	ID      string
	ShortID string
}

type PreViewModel struct {
	ID        string
	Name      string
	State     string
	UserID    string
	Variables map[string][]PreViewStructure
	Maps      []PreViewStructure
	Lists     []PreViewStructure
}

type ViewVariable struct {
	Name    string
	ID      string
	ShortID string
}

type View struct {
	ID         string          `json:"id"`
	Name       string          `json:"name"`
	State      string          `json:"state"`
	UserID     string          `json:"userID"`
	Structures []StructureView `json:"structures"`

	Variables map[string][]StructureView `json:"variables"`
	Maps      []StructureView            `json:"maps"`
	Lists     []StructureView            `json:"lists"`
}

func newView(model PreViewModel) View {
	variables := make(map[string][]StructureView)
	for key, value := range model.Variables {
		variables[key] = sdk.Map(value, func(idx int, value PreViewStructure) StructureView {
			return StructureView{
				Name:    value.Name,
				ID:      value.ID,
				ShortID: value.ShortID,
			}
		})
	}

	return View{
		ID:         model.ID,
		Structures: make([]StructureView, 0),
		Name:       model.Name,
		State:      model.State,
		UserID:     model.UserID,
		Variables:  variables,
		Maps: sdk.Map(model.Maps, func(idx int, value PreViewStructure) StructureView {
			return StructureView{
				Name:    value.Name,
				ID:      value.ID,
				ShortID: value.ShortID,
			}
		}),
		Lists: sdk.Map(model.Lists, func(idx int, value PreViewStructure) StructureView {
			return StructureView{
				Name:    value.Name,
				ID:      value.ID,
				ShortID: value.ShortID,
			}
		}),
	}
}
