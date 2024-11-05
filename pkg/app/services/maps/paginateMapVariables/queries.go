package paginateMapVariables

import (
	"creatif/pkg/lib/storage"
	"fmt"
)

type GroupsQuery struct {
	VariableID string `gorm:"column:variable_id"`
	GroupName  string `gorm:"column:name"`
}

type groupsReturnType = []map[string][]string

func getItemGroups(ids []string) (groupsReturnType, error) {
	sql := fmt.Sprintf("SELECT g.name, vg.variable_id FROM declarations.groups AS g INNER JOIN declarations.variable_groups AS vg ON vg.variable_id IN(?) AND g.id = ANY(vg.groups) GROUP BY vg.variable_id, g.name")

	var m []GroupsQuery
	if res := storage.Gorm().Raw(sql, ids).Scan(&m); res.Error != nil {
		return nil, res.Error
	}

	results := make([]map[string][]string, 0)
	visited := make([]string, 0)
	for _, v := range m {
		id := v.VariableID
		alreadyPopulated := false
		for _, visit := range visited {
			if visit == id {
				alreadyPopulated = true
				break
			}
		}

		visited = append(visited, id)
		if !alreadyPopulated {
			result := make(map[string][]string)
			result[id] = make([]string, 0)

			for _, p := range m {
				if id == p.VariableID {
					result[id] = append(result[id], p.GroupName)
				}
			}

			results = append(results, result)
		}
	}

	return results, nil
}
