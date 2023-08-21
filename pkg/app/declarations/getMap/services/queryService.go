package services

import (
	"creatif/pkg/lib/sdk"
	"fmt"
	"strings"
)

type QueryStrategy interface {
	GetQuery() string
}

type FullQueryStrategy struct{}

func (f FullQueryStrategy) GetQuery() string {
	return `SELECT 
    n.id, 
    n.name, 
    n.groups, 
    n.behaviour, 
    n.type, 
    n.metadata, 
    n.created_at, 
    n.updated_at
		FROM declarations.map_nodes AS m
		INNER JOIN declarations.map_nodes AS mn ON m.id = mn.map_id
		INNER JOIN declarations.nodes AS n ON n.id = mn.node_id
		WHERE m.id = ?
`
}

type NamesOnlyQueryStrategy struct{}

func (f NamesOnlyQueryStrategy) GetQuery() string {
	return `SELECT 
    n.id,
    n.name,
    n.groups,
    n.behaviour,
    n.type,
    n.metadata,
    n.created_at,
    n.updated_at
		FROM declarations.map_nodes AS m
		INNER JOIN declarations.map_nodes AS mn ON m.id = mn.map_id
		INNER JOIN declarations.nodes AS n ON n.id = mn.node_id
		WHERE m.id = ?
`
}

type CustomFieldsQueryStrategy struct {
	validFields  []string
	chosenFields []string
}

func (c CustomFieldsQueryStrategy) GetQuery() string {
	fields := strings.Join(sdk.Map(c.chosenFields, func(idx int, value string) string {
		return fmt.Sprintf("n.%s", value)
	}), ",")

	fields = ", " + fields

	return fmt.Sprintf(`
		SELECT m.id, m.name AS mapName, m.created_at AS mapCreatedAt, m.updated_at AS mapUpdatedAt %s
		FROM declarations.map_nodes AS m
		INNER JOIN declarations.map_nodes AS mn ON m.id = mn.map_id
		INNER JOIN declarations.nodes AS n ON n.id = mn.node_id
		WHERE m.id = ?
`, fields)
}

func CreateStrategy(returnType string, chosenFields []string) QueryStrategy {
	if returnType == "full" {
		return FullQueryStrategy{}
	} else if returnType == "names" {
		return NamesOnlyQueryStrategy{}
	} else if len(chosenFields) != 0 {
		return CustomFieldsQueryStrategy{
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
