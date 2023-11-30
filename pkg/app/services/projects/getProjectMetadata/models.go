package getProjectMetadata

type LogicModel struct {
	ID             string
	Name           string
	State          string
	UserID         string `gorm:"column:user_id"`
	Map            string `gorm:"column:map_name"`
	List           string `gorm:"column:list_name"`
	VariableName   string `gorm:"column:variable_name"`
	VariableLocale string `gorm:"column:variable_locale"`
	MapLocale      string `gorm:"column:map_locale"`
	ListLocale     string `gorm:"column:list_locale"`
}

type PreViewModel struct {
	ID        string
	Name      string
	State     string
	UserID    string
	Variables map[string][]string
	Maps      map[string][]string
	Lists     map[string][]string
}

type View struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	State  string `json:"state"`
	UserID string `json:"userID"`

	Variables map[string][]string `json:"variables"`
	Maps      map[string][]string `json:"maps"`
	Lists     map[string][]string `json:"lists"`
}

func newView(model PreViewModel) View {
	return View{
		ID:        model.ID,
		Name:      model.Name,
		State:     model.State,
		UserID:    model.UserID,
		Variables: model.Variables,
		Maps:      model.Maps,
		Lists:     model.Lists,
	}
}
