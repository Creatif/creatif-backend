package create

import (
	"creatif/pkg/app/domain/declarations"
)

var validFields = []string{
	"id",
	"name",
	"type",
	"behaviour",
	"metadata",
	"groups",
	"created_at",
	"updated_at",
}

type GetMapModel struct {
	// this can be map name or an id of the map
	ID string `json:"id"`
	// this can be, 'full' | names
	Return string
	// this can be individual fields of the node to return, reduces returned data
	// if the user needs only metadata, only metadata will be returned
	// name is always returned
	Fields []string

	validFields []string
	// TODO: Add project ID prop here
}

func NewGetMapModel(id string) GetMapModel {
	return GetMapModel{
		ID:          id,
		validFields: validFields,
	}
}

type View struct {
	ID    string      `json:"id"`
	Name  string      `json:"name"`
	Nodes interface{} `json:"nodes"`
}

func newView(model LogicModel) View {
	return View{
		ID:    model.nodeMap.ID,
		Name:  model.nodeMap.Name,
		Nodes: model.nodes,
	}
}

type LogicModel struct {
	nodeMap declarations.Map
	nodes   []map[string]interface{}
}
