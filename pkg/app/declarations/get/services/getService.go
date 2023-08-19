package services

import (
	"creatif/pkg/app/domain/assignments"
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/lib/appErrors"
	"creatif/pkg/lib/constants"
	"creatif/pkg/lib/sdk"
	"creatif/pkg/lib/storage"
	"fmt"
	"github.com/lib/pq"
	"gorm.io/datatypes"
	"time"
)

type Node struct {
	ID string

	Name      string
	Type      string // text,image,file,boolean
	Behaviour string // readonly,modifiable
	Groups    pq.StringArray
	Metadata  datatypes.JSON
	Value     interface{}

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type GetService struct {
	id string
}

func NewGetService(id string) GetService {
	return GetService{id: id}
}

func (g GetService) GetNode(byId func(id string) (declarations.Node, error), byName func(name string) (declarations.Node, error)) (Node, error) {
	if sdk.IsValidUuid(g.id) {
		node, err := byId(g.id)
		if err != nil {
			return Node{}, err
		}

		return queryValue(node)
	}

	node, err := byName(g.id)
	if err != nil {
		return Node{}, err
	}

	return queryValue(node)
}

func assignTableToQuery(table string) string {
	return fmt.Sprintf(`
SELECT ant.* FROM %s AS ant
	INNER JOIN assignments.nodes AS an ON an.id = ant.assignment_node_id
	INNER JOIN declarations.nodes AS dn ON dn.id = an.declaration_node_id
	WHERE dn.id = ?
`, table)
}

func queryValue(node declarations.Node) (Node, error) {
	serviceNode := declarationNodeToServiceNode(node)

	var textNode assignments.NodeText
	var boolNode assignments.NodeBoolean
	if serviceNode.Type == constants.ValueTextType {
		if res := storage.Gorm().Raw(assignTableToQuery("assignments.node_text"), serviceNode.ID).
			Scan(&textNode); res.Error != nil {
			return Node{}, appErrors.NewDatabaseError(res.Error).AddError("Node.Get.Logic", nil)
		}

		serviceNode.Value = textNode.Value
	} else if serviceNode.Type == constants.ValueBooleanType {
		if res := storage.Gorm().Raw(assignTableToQuery("assignments.node_boolean"), serviceNode.ID).Scan(&boolNode); res.Error != nil {
			return Node{}, appErrors.NewDatabaseError(res.Error).AddError("Node.Get.Logic", nil)
		}

		serviceNode.Value = boolNode.Value
	}

	return serviceNode, nil
}

func declarationNodeToServiceNode(node declarations.Node) Node {
	return Node{
		ID:        node.ID,
		Name:      node.Name,
		Type:      node.Type,
		Behaviour: node.Behaviour,
		Groups:    node.Groups,
		Metadata:  node.Metadata,
		Value:     nil,
		CreatedAt: node.CreatedAt,
		UpdatedAt: node.UpdatedAt,
	}
}
