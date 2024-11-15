package connections

import (
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/lib/storage"
	"errors"
	"fmt"
)

type Connection struct {
	Path          string
	StructureType string
	VariableID    string
}

func CheckConnectionsIntegrity(connections []Connection) error {
	for _, c := range connections {
		var id string
		if c.StructureType == "map" {
			res := storage.Gorm().Raw(fmt.Sprintf("SELECT id FROM %s WHERE id = ?", (declarations.MapVariable{}).TableName()), c.VariableID).Scan(&id)
			if res.Error != nil {
				return res.Error
			}

			if res.RowsAffected == 0 {
				return errors.New(fmt.Sprintf("Connection does not exist. ID: %s, path: %s, structure type: %s", c.VariableID, c.Path, c.StructureType))
			}
		}

		if c.StructureType == "list" {
			res := storage.Gorm().Raw(fmt.Sprintf("SELECT id FROM %s WHERE id = ?", (declarations.ListVariable{}).TableName()), c.VariableID).Scan(&id)
			if res.Error != nil {
				return res.Error
			}

			if res.RowsAffected == 0 {
				return errors.New(fmt.Sprintf("Connection does not exist. ID: %s, path: %s, structure type: %s", c.VariableID, c.Path, c.StructureType))
			}
		}
	}

	return nil
}
