package shared

import (
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/lib/storage"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type Reference struct {
	Name          string
	StructureName string
	StructureType string
	VariableID    string
}

type UpdateReference struct {
	Name          string
	StructureName string
	StructureType string
	VariableID    string
}

type ParentReference struct {
	StructureID   string
	ParentShortID string
	ID            string `gorm:"column:id"`
}

func CreateDeclarationReferences(refs []Reference, ownerId, childStructureId, projectId string) ([]declarations.Reference, error) {
	references := make([]declarations.Reference, 0)
	for _, r := range refs {
		pr, err := getParentReference(r.StructureName, r.StructureType, r.VariableID, ownerId)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Reference with ID '%s' not found. Structure name: %s; Structure type: %s. Underlying error: %s", r.VariableID, r.StructureName, r.StructureType, err.Error()))
		}

		ref := declarations.NewReference(r.Name, r.StructureType, "map", pr.ID, ownerId, pr.StructureID, childStructureId, projectId)
		references = append(references, ref)
	}

	return references, nil
}

func UpdateReferences(refs []UpdateReference, structureId, ownerId, projectId string, tx *gorm.DB) error {
	// if there are not refs sent, clear the refs since user might not have frontend validation enabled
	if len(refs) == 0 {
		if err := deleteAllRefsByChild(ownerId, tx); err != nil {
			return err
		}

		return nil
	}

	updateableReferences, err := getRefsByChild(ownerId, tx)
	if err != nil {
		return err
	}

	// 1. if the incoming ref is in update refs, update it
	// 2. if the incoming ref is NOT in update refs, create it
	// 3. if the update ref is not in incoming refs, delete it

	// update
	for _, incomingRef := range refs {
		updatePerformed := false

		// to update
		for _, updatableRef := range updateableReferences {
			if updatableRef.Name == incomingRef.Name {
				pr, err := getParentReference(incomingRef.StructureName, incomingRef.StructureType, incomingRef.VariableID, structureId)
				if err != nil {
					return errors.New(fmt.Sprintf("Reference with ID '%s' not found. Structure name: %s; Structure type: %s. Underlying error: %s", incomingRef.VariableID, incomingRef.StructureName, incomingRef.StructureType, err.Error()))
				}

				res := tx.Exec(fmt.Sprintf("UPDATE %s SET parent_type = ?, parent_id = ? WHERE child_id = ? AND name = ?", (declarations.Reference{}).TableName()), incomingRef.StructureType, pr.ID, ownerId, incomingRef.Name)
				if res.Error != nil {
					return res.Error
				}

				if res.RowsAffected == 0 {
					return errors.New("No references found. This must be an error since update references have been sent.")
				}

				updatePerformed = true
			}
		}

		// create
		if !updatePerformed {
			pr, err := getParentReference(incomingRef.StructureName, incomingRef.StructureType, incomingRef.VariableID, ownerId)
			if err != nil {
				return errors.New(fmt.Sprintf("Reference with ID '%s' not found. Structure name: %s; Structure type: %s. Underlying error: %s", incomingRef.VariableID, incomingRef.StructureName, incomingRef.StructureType, err.Error()))
			}

			ref := declarations.NewReference(incomingRef.Name, incomingRef.StructureType, "map", pr.ID, ownerId, structureId, "", projectId)
			tx.Create(&ref)
		}
	}

	// update ref delete
	for _, updateRef := range updateableReferences {
		found := false
		for _, incomingRef := range refs {
			if incomingRef.Name == updateRef.Name {
				found = true
				break
			}
		}

		if !found {
			res := tx.Exec(fmt.Sprintf(`DELETE FROM %s WHERE child_id = ? AND name = ?`, (declarations.Reference{}).TableName()), ownerId, updateRef.Name)

			if res.Error != nil {
				return res.Error
			}

			if res.RowsAffected == 0 {
				return errors.New("Invalid reference. References should have been deleted but it was not.")
			}
		}
	}

	return nil
}

func RemoveAsParent(structureType, parentId string) error {
	res := storage.Gorm().Exec(fmt.Sprintf("DELETE FROM %s WHERE parent_type = ? AND (parent_id = ?)", (declarations.Reference{}).TableName()), structureType, parentId)

	return res.Error
}

func RemoveAsChild(structureType, childId string) error {
	res := storage.Gorm().Exec(fmt.Sprintf("DELETE FROM %s WHERE child_type = ? AND (child_id = ?)", (declarations.Reference{}).TableName()), structureType, childId)

	return res.Error
}

func getParentReference(structureName, structureType, variableId, structureId string) (ParentReference, error) {
	if structureType == "list" {
		sql := fmt.Sprintf(`SELECT lv.id AS id, l.id AS structure_id FROM declarations.list_variables AS lv INNER JOIN declarations.lists AS l ON l.id = lv.list_id AND (l.name = ? OR l.id = ? OR l.short_id = ?) AND lv.id = ?`)

		var pr ParentReference
		if res := storage.Gorm().Raw(sql, structureName, structureName, structureName, variableId).Scan(&pr); res.Error != nil {
			return ParentReference{}, res.Error
		}

		if pr.StructureID == structureId {
			return ParentReference{}, errors.New("Invalid parent reference. A reference can not have itself as the parent")
		}

		return pr, nil
	}

	if structureType == "map" {
		sql := fmt.Sprintf(`SELECT lv.id AS id, l.id AS structure_id FROM declarations.map_variables AS lv INNER JOIN declarations.maps AS l ON l.id = lv.map_id AND (l.name = ? OR l.id = ? OR l.short_id = ?) AND lv.id = ?`)

		var pr ParentReference
		if res := storage.Gorm().Raw(sql, structureName, structureName, structureName, variableId).Scan(&pr); res.Error != nil {
			return ParentReference{}, res.Error
		}

		if pr.StructureID == structureId {
			return ParentReference{}, errors.New("Invalid parent reference. A reference can not have itself as the parent")
		}

		return pr, nil
	}

	sql := fmt.Sprintf(`SELECT id FROM declarations.variables WHERE id = ?`)

	var pr ParentReference
	if res := storage.Gorm().Raw(sql, variableId).Scan(&pr); res.Error != nil {
		return pr, res.Error
	}

	return pr, nil
}

func deleteAllRefsByChild(childId string, tx *gorm.DB) error {
	res := tx.Exec(fmt.Sprintf(`DELETE FROM %s WHERE child_id = ?`, (declarations.Reference{}).TableName()), childId)
	return res.Error
}

func getRefsByChild(childId string, tx *gorm.DB) ([]declarations.Reference, error) {
	sql := fmt.Sprintf(`SELECT * FROM %s WHERE child_id = ?`, (declarations.Reference{}).TableName())

	var refs []declarations.Reference
	res := tx.Raw(sql, childId).Scan(&refs)

	if res.Error != nil {
		return refs, res.Error
	}

	return refs, nil
}
