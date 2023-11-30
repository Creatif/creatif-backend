package getProjectMetadata

type LogicModel struct {
	ID            string
	Name          string
	State         string
	UserID        string   `gorm:"column:user_id"`
	VariableNames []string `gorm:"column:variable_names"`
	MapNames      []string `gorm:"column:map_names"`
	ListNames     []string `gorm:"column:list_names"`
}

type View struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	State  string `json:"state"`
	UserID string `json:"userID"`

	VariableNames []string `json:"variableNames"`
	MapNames      []string `json:"mapNames"`
	ListNames     []string `json:"listNames"`
}

func newView(model LogicModel) View {
	return View{
		ID:            model.ID,
		Name:          model.Name,
		State:         model.State,
		UserID:        model.UserID,
		VariableNames: model.VariableNames,
		MapNames:      model.MapNames,
		ListNames:     model.ListNames,
	}
}
