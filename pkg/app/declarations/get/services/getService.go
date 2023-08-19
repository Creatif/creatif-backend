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
		var nodeWithValueQuery Node

		node, err := byId(g.id)
		if err != nil {
			return Node{}, err
		}

		nodeWithValueQuery.ID = node.ID
		nodeWithValueQuery.Name = node.Name
		nodeWithValueQuery.Type = node.Type
		nodeWithValueQuery.Groups = node.Groups
		nodeWithValueQuery.Behaviour = node.Behaviour
		nodeWithValueQuery.Metadata = node.Metadata
		nodeWithValueQuery.CreatedAt = node.CreatedAt
		nodeWithValueQuery.UpdatedAt = node.UpdatedAt

		var textNode assignments.NodeText
		var boolNode assignments.NodeBoolean
		if nodeWithValueQuery.Type == constants.ValueTextType {
			if res := storage.Gorm().Raw(assignTableToQuery("assignments.node_text"), nodeWithValueQuery.ID).
				Scan(&textNode); res.Error != nil {
				return Node{}, appErrors.NewDatabaseError(res.Error).AddError("Node.Get.Logic", nil)
			}

			nodeWithValueQuery.Value = textNode.Value
		} else if nodeWithValueQuery.Type == constants.ValueBooleanType {
			if res := storage.Gorm().Raw(assignTableToQuery("assignments.node_boolean"), nodeWithValueQuery.ID).Scan(&boolNode); res.Error != nil {
				return Node{}, appErrors.NewDatabaseError(res.Error).AddError("Node.Get.Logic", nil)
			}

			nodeWithValueQuery.Value = boolNode.Value
		}

		return nodeWithValueQuery, nil
	}

	var nodeWithValueQuery Node
	node, err := byName(g.id)
	if err != nil {
		return Node{}, err
	}

	nodeWithValueQuery.ID = node.ID
	nodeWithValueQuery.Name = node.Name
	nodeWithValueQuery.Type = node.Type
	nodeWithValueQuery.Groups = node.Groups
	nodeWithValueQuery.Behaviour = node.Behaviour
	nodeWithValueQuery.Metadata = node.Metadata
	nodeWithValueQuery.CreatedAt = node.CreatedAt
	nodeWithValueQuery.UpdatedAt = node.UpdatedAt

	var textNode assignments.NodeText
	var boolNode assignments.NodeBoolean
	if nodeWithValueQuery.Type == constants.ValueTextType {
		if res := storage.Gorm().Raw(assignTableToQuery("assignments.node_text"), nodeWithValueQuery.ID).
			Scan(&textNode); res.Error != nil {
			return Node{}, appErrors.NewDatabaseError(res.Error).AddError("Node.Get.Logic", nil)
		}

		nodeWithValueQuery.Value = textNode.Value
	} else if nodeWithValueQuery.Type == constants.ValueBooleanType {
		if res := storage.Gorm().Raw(assignTableToQuery("assignments.node_boolean"), nodeWithValueQuery.ID).Scan(&boolNode); res.Error != nil {
			return Node{}, appErrors.NewDatabaseError(res.Error).AddError("Node.Get.Logic", nil)
		}

		nodeWithValueQuery.Value = boolNode.Value
	}

	return nodeWithValueQuery, nil
}

func assignTableToQuery(table string) string {
	return fmt.Sprintf(`
SELECT ant.* FROM %s AS ant
	INNER JOIN assignments.nodes AS an ON an.id = ant.assignment_node_id
	INNER JOIN declarations.nodes AS dn ON dn.id = an.declaration_node_id
	WHERE dn.id = ?
`, table)
}
