package shared

import (
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/lib/storage"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type Reference struct {
	StructureName string
	StructureType string
	VariableID    string
}

type UpdateReference struct {
	ID            string
	StructureName string
	StructureType string
	VariableID    string
}

type ParentReference struct {
	StructureID   string
	ParentShortID string
	ID            string `gorm:"column:id"`
}

func CreateDeclarationReferences(refs []Reference, currentId, currentShortId string) ([]declarations.Reference, error) {
	references := make([]declarations.Reference, 0)
	for _, r := range refs {
		pr, err := getParentReference(r.StructureName, r.StructureType, r.VariableID, currentId)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Reference with ID '%s' not found. Structure name: %s; Structure type: %s. Underlying error: %s", r.VariableID, r.StructureName, r.StructureType, err.Error()))
		}

		ref := declarations.NewReference(r.StructureType, "map", pr.ID, pr.ParentShortID, currentId, currentShortId)
		references = append(references, ref)
	}

	return references, nil
}

func UpdateReferences(refs []UpdateReference, structureId string, tx *gorm.DB) error {
	for _, r := range refs {
		pr, err := getParentReference(r.StructureName, r.StructureType, r.VariableID, structureId)
		if err != nil {
			return errors.New(fmt.Sprintf("Reference with ID '%s' not found. Structure name: %s; Structure type: %s. Underlying error: %s", r.VariableID, r.StructureName, r.StructureType, err.Error()))
		}

		res := storage.Gorm().Exec(fmt.Sprintf("UPDATE %s SET parent_type = ?, parent_id = ?, parent_short_id = ? WHERE id = ?", (declarations.Reference{}).TableName()), r.StructureType, pr.ID, pr.ParentShortID, r.ID)

		if res.Error != nil {
			return res.Error
		}

		if res.RowsAffected == 0 {
			return errors.New("No references found. This must be an error since update references have been sent.")
		}
	}

	return nil
}

func RemoveAsParent(structureType, parentId string) error {
	res := storage.Gorm().Exec(fmt.Sprintf("DELETE FROM %s WHERE parent_type = ? AND (parent_id = ? OR parent_short_id = ?)", (declarations.Reference{}).TableName()), structureType, parentId, parentId)

	return res.Error
}

func RemoveAsChild(structureType, childId string) error {
	res := storage.Gorm().Exec(fmt.Sprintf("DELETE FROM %s WHERE child_type = ? AND (child_id = ? OR child_short_id = ?)", (declarations.Reference{}).TableName()), structureType, childId, childId)

	return res.Error
}

func getParentReference(structureName, structureType, variableId, structureId string) (ParentReference, error) {
	if structureType == "list" {
		sql := fmt.Sprintf(`SELECT lv.id AS id, lv.short_id as parent_short_id, l.id AS structure_id FROM declarations.list_variables AS lv INNER JOIN declarations.lists AS l ON l.id = lv.list_id AND l.name = ? AND lv.id = ?`)

		var pr ParentReference
		if res := storage.Gorm().Raw(sql, structureName, variableId).Scan(&pr); res.Error != nil {
			return ParentReference{}, res.Error
		}

		if pr.StructureID == structureId {
			return ParentReference{}, errors.New("Invalid parent reference. A reference can not have itself as the parent")
		}

		return pr, nil
	}

	if structureType == "map" {
		sql := fmt.Sprintf(`SELECT lv.id AS id, lv.short_id as parent_short_id, l.id AS structure_id FROM declarations.map_variables AS lv INNER JOIN declarations.maps AS l ON l.id = lv.map_id AND l.name = ? AND lv.id = ?`)

		var pr ParentReference
		if res := storage.Gorm().Raw(sql, structureName, variableId).Scan(&pr); res.Error != nil {
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
