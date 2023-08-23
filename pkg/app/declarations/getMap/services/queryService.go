package services

import (
	"creatif/pkg/lib/sdk"
	"fmt"
	"strings"
)

type QueryStrategy interface {
	GetQuery() string
	Name() string
}

type FullQueryStrategy struct {
	name string
}

func (f FullQueryStrategy) GetQuery() string {
	return `SELECT 
    n.id, 
    n.name, 
    n.groups, 
    n.behaviour, 
    vn.value,
    n.type,
    n.metadata, 
    n.created_at, 
    n.updated_at
		FROM declarations.map_nodes AS mn
		INNER JOIN declarations.nodes AS n ON n.id = mn.node_id
		INNER JOIN declarations.node_map AS m ON m.id = mn.map_id
		INNER JOIN assignments.nodes AS an ON an.declaration_node_id = n.id
		INNER JOIN assignments.value_node AS vn ON vn.assignment_node_id = an.id
		WHERE m.id = ?
`
}

func (f FullQueryStrategy) Name() string {
	return f.name
}

type NamesOnlyQueryStrategy struct {
	name string
}

func (f NamesOnlyQueryStrategy) GetQuery() string {
	return `SELECT 
    n.id,
    n.name
		FROM declarations.map_nodes AS mn
		INNER JOIN declarations.nodes AS n ON n.id = mn.node_id
		INNER JOIN declarations.node_map AS m ON m.id = mn.map_id
		INNER JOIN assignments.nodes AS an ON an.declaration_node_id = n.id
		INNER JOIN assignments.value_node AS vn ON vn.assignment_node_id = an.id
		WHERE m.id = ?
`
}

func (f NamesOnlyQueryStrategy) Name() string {
	return f.name
}

type CustomFieldsQueryStrategy struct {
	name         string
	validFields  []string
	chosenFields []string
}

func (c CustomFieldsQueryStrategy) GetQuery() string {
	fields := strings.Join(sdk.Map(c.chosenFields, func(idx int, value string) string {
		columnAlias := "n"
		if value == "value" {
			columnAlias = "vn"
		}

		return fmt.Sprintf("%s.%s", columnAlias, value)
	}), ",")

	fields = ", " + fields

	return fmt.Sprintf(`
SELECT 
    n.id, 
    n.name %s
		FROM declarations.map_nodes AS mn
		INNER JOIN declarations.nodes AS n ON n.id = mn.node_id
		INNER JOIN declarations.node_map AS m ON m.id = mn.map_id
		INNER JOIN assignments.nodes AS an ON an.declaration_node_id = n.id
		INNER JOIN assignments.value_node AS vn ON vn.assignment_node_id = an.id
		WHERE m.id = ?
`, fields)
}

func (c CustomFieldsQueryStrategy) Name() string {
	return c.name
}

func CreateStrategy(returnType string, chosenFields []string) QueryStrategy {
	if returnType == "full" {
		return FullQueryStrategy{
			name: "fullQueryStrategy",
		}
	} else if returnType == "names" {
		return NamesOnlyQueryStrategy{
			name: "namesOnlyStrategy",
		}
	} else if len(chosenFields) != 0 {
		return CustomFieldsQueryStrategy{
			name: "customFieldsStrategy",
			validFields: []string{
				"id",
				"name",
				"type",
				"behaviour",
				"metadata",
				"groups",
				"created_at",
				"updated_at",
			},
			chosenFields: chosenFields,
		}
	}

	return NamesOnlyQueryStrategy{}
}
