package shared

import (
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/lib/appErrors"
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

func CreateDeclarationReferences(refs []Reference, structureId, childId, projectId string) ([]declarations.Reference, error) {
	references := make([]declarations.Reference, 0)
	for _, r := range refs {
		pr, err := getParentReference(r.StructureName, r.StructureType, r.VariableID, structureId)
		if err != nil {
			return nil, err
		}

		ref := declarations.NewReference(r.Name, r.StructureType, "map", pr.ID, childId, pr.StructureID, structureId, projectId)
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
					return err
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
				return err
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

func RemoveAsParent(parentId string) error {
	res := storage.Gorm().Exec(fmt.Sprintf("DELETE FROM %s WHERE parent_id = ?", (declarations.Reference{}).TableName()), parentId)

	return res.Error
}

func RemoveManyAsParent(parentIds []string) error {
	res := storage.Gorm().Exec(fmt.Sprintf("DELETE FROM %s WHERE parent_id IN(?)", (declarations.Reference{}).TableName()), parentIds)

	return res.Error
}

func RemoveManyAsChild(childIds []string) error {
	res := storage.Gorm().Exec(fmt.Sprintf("DELETE FROM %s WHERE child_id IN(?)", (declarations.Reference{}).TableName()), childIds)

	return res.Error
}

func RemoveAsChild(childId string) error {
	res := storage.Gorm().Exec(fmt.Sprintf("DELETE FROM %s WHERE child_id = ?", (declarations.Reference{}).TableName()), childId)

	return res.Error
}

func getParentReference(structureName, structureType, variableId, structureId string) (ParentReference, error) {
	if structureType == "list" {
		sql := fmt.Sprintf(`SELECT lv.id AS id, l.id AS structure_id FROM declarations.list_variables AS lv INNER JOIN declarations.lists AS l ON l.id = lv.list_id AND (l.name = ? OR l.id = ? OR l.short_id = ?) AND lv.id = ?`)

		var pr ParentReference
		res := storage.Gorm().Raw(sql, structureName, structureName, structureName, variableId).Scan(&pr)

		if res.Error != nil {
			return ParentReference{}, appErrors.NewValidationError(map[string]string{
				"referenceInvalid": res.Error.Error(),
			})
		}

		if pr.StructureID == structureId {
			return ParentReference{}, appErrors.NewValidationError(map[string]string{
				"referenceInvalid": "Invalid reference. A reference cannot be a parent to itself.",
			})
		}

		if res.RowsAffected == 0 {
			return ParentReference{}, appErrors.NewValidationError(map[string]string{
				"referenceInvalid": "Invalid reference. Parent reference not found.",
			})
		}

		return pr, nil
	}

	if structureType == "map" {
		sql := fmt.Sprintf(`SELECT lv.id AS id, l.id AS structure_id FROM declarations.map_variables AS lv INNER JOIN declarations.maps AS l ON l.id = lv.map_id AND (l.name = ? OR l.id = ? OR l.short_id = ?) AND lv.id = ?`)

		var pr ParentReference
		res := storage.Gorm().Raw(sql, structureName, structureName, structureName, variableId).Scan(&pr)

		if res.Error != nil {
			return ParentReference{}, appErrors.NewValidationError(map[string]string{
				"referenceInvalid": res.Error.Error(),
			})
		}

		if pr.StructureID == structureId {
			return ParentReference{}, appErrors.NewValidationError(map[string]string{
				"referenceInvalid": "Invalid reference. A reference cannot be a parent to itself.",
			})
		}

		if res.RowsAffected == 0 {
			return ParentReference{}, appErrors.NewValidationError(map[string]string{
				"referenceInvalid": "Invalid reference. Parent reference not found.",
			})
		}

		return pr, nil
	}

	sql := fmt.Sprintf(`SELECT id FROM declarations.variables WHERE id = ?`)

	var pr ParentReference
	res := storage.Gorm().Raw(sql, variableId).Scan(&pr)

	if res.Error != nil {
		return ParentReference{}, appErrors.NewValidationError(map[string]string{
			"referenceInvalid": res.Error.Error(),
		})
	}

	if res.RowsAffected == 0 {
		return ParentReference{}, appErrors.NewValidationError(map[string]string{
			"referenceInvalid": "Invalid reference. Parent reference not found.",
		})
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
