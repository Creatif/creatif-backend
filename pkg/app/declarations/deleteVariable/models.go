package deleteVariable

type Model struct {
	Name string `json:"name"`
}

func NewModel(name string) Model {
	return Model{
		Name: name,
	}
}
