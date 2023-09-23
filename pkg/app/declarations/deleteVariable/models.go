package deleteVariable

type Model struct {
	Name      string `json:"name"`
	ProjectID string `json:"projectID"`
}

func NewModel(projectId, name string) Model {
	return Model{
		Name:      name,
		ProjectID: projectId,
	}
}
