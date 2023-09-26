package deleteVariable

type Model struct {
	Name      string `json:"name"`
	EntryName string `json:"entryName"`
	ProjectID string `json:"projectID"`
}

func NewModel(projectId, name, entryName string) Model {
	return Model{
		Name:      name,
		ProjectID: projectId,
		EntryName: entryName,
	}
}
