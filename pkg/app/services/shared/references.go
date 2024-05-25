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

func CreateDeclarationReferences(refs []Reference, structureId, childId, childType, projectId string) ([]declarations.Reference, error) {
	references := make([]declarations.Reference, 0)

	for _, r := range refs {
		pr, err := getParentReference(r.StructureType, r.VariableID, structureId)
		if err != nil {
			return nil, err
		}

		ref := declarations.NewReference(r.Name, r.StructureType, childType, pr.ID, childId, pr.StructureID, structureId, projectId)
		references = append(references, ref)
	}

	return references, nil
}

func UpdateReferences(refs []UpdateReference, childStructureId, ownerId, projectId string, tx *gorm.DB) error {
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
			fmt.Println(incomingRef.Name, incomingRef.StructureName, incomingRef.StructureType, incomingRef.VariableID, childStructureId)
			if updatableRef.Name == incomingRef.Name {
				pr, err := getParentReference(incomingRef.StructureType, incomingRef.VariableID, childStructureId)
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
			pr, err := getParentReference(incomingRef.StructureType, incomingRef.VariableID, ownerId)
			if err != nil {
				return err
			}

			ref := declarations.NewReference(incomingRef.Name, incomingRef.StructureType, "map", pr.ID, ownerId, pr.StructureID, childStructureId, projectId)
			if res := tx.Create(&ref); res.Error != nil {
				return res.Error
			}
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

func RemoveAsParent(parentId string, tx *gorm.DB) error {
	res := tx.Exec(fmt.Sprintf("DELETE FROM %s WHERE parent_id = ?", (declarations.Reference{}).TableName()), parentId)

	return res.Error
}

func IsParent(variableId string) error {
	var id string
	res := storage.Gorm().Raw(fmt.Sprintf("SELECT id FROM %s WHERE parent_id = ?", (declarations.Reference{}).TableName()), variableId).Scan(&id)
	id = id

	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected != 0 {
		return errors.New("is_parent")
	}

	return nil
}

func RemoveManyAsParent(parentIds []string, tx *gorm.DB) error {
	res := tx.Exec(fmt.Sprintf("DELETE FROM %s WHERE parent_id IN(?)", (declarations.Reference{}).TableName()), parentIds)

	return res.Error
}

func RemoveManyAsChild(childIds []string, tx *gorm.DB) error {
	res := tx.Exec(fmt.Sprintf("DELETE FROM %s WHERE child_id IN(?)", (declarations.Reference{}).TableName()), childIds)

	return res.Error
}

func RemoveAsChild(childId string, tx *gorm.DB) error {
	res := tx.Exec(fmt.Sprintf("DELETE FROM %s WHERE child_id = ?", (declarations.Reference{}).TableName()), childId)

	return res.Error
}

func getParentReference(structureType, variableId, structureId string) (ParentReference, error) {
	if structureType == "list" {
		sql := fmt.Sprintf(`SELECT lv.id AS id, lv.list_id AS structure_id FROM declarations.list_variables AS lv WHERE lv.id = ?`)

		var pr ParentReference
		res := storage.Gorm().Raw(sql, variableId).Scan(&pr)

		if res.Error != nil {
			return ParentReference{}, errors.New(fmt.Sprintf("referenceInvalid:%s", res.Error.Error()))
		}

		if pr.StructureID == structureId {
			return ParentReference{}, errors.New(fmt.Sprintf("referenceInvalid:%s", "Invalid reference. A reference cannot be a parent to itself."))
		}

		if res.RowsAffected == 0 {
			return ParentReference{}, errors.New(fmt.Sprintf("referenceInvalid:%s", "Invalid reference. Parent reference not found"))
		}

		return pr, nil
	}

	if structureType == "map" {
		sql := fmt.Sprintf(`SELECT lv.id AS id, lv.map_id AS structure_id FROM declarations.map_variables AS lv WHERE lv.id = ?`)

		var pr ParentReference
		res := storage.Gorm().Raw(sql, variableId).Scan(&pr)

		if res.Error != nil {
			return ParentReference{}, errors.New(fmt.Sprintf("referenceInvalid:%s", res.Error.Error()))
		}

		if pr.StructureID == structureId {
			return ParentReference{}, errors.New(fmt.Sprintf("referenceInvalid:%s", "Invalid reference. A reference cannot be a parent to itself."))
		}

		if res.RowsAffected == 0 {
			return ParentReference{}, errors.New(fmt.Sprintf("referenceInvalid:%s", "Invalid reference. Parent reference not found"))
		}

		return pr, nil
	}

	sql := fmt.Sprintf(`SELECT id FROM declarations.variables WHERE id = ?`)

	var pr ParentReference
	res := storage.Gorm().Raw(sql, variableId).Scan(&pr)

	if res.Error != nil {
		return ParentReference{}, errors.New(fmt.Sprintf("referenceInvalid:%s", res.Error.Error()))
	}

	if res.RowsAffected == 0 {
		return ParentReference{}, errors.New(fmt.Sprintf("referenceInvalid:%s", "Invalid reference. Parent reference not found"))
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
