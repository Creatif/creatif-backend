package shared

import (
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/lib/appErrors"
	"creatif/pkg/lib/storage"
	"fmt"
)

type QueryReference struct {
	ID                string
	Name              string
	ParentType        string
	ChildType         string
	ParentID          string
	ChildID           string
	ChildStructureID  string
	ParentStructureID string
	StructureName     string
}

func QueryReferences(id, projectId string) ([]QueryReference, error) {
	var parents []QueryReference
	if err := queryParentReferences(id, projectId, &parents); err != nil {
		return nil, err
	}

	var children []QueryReference
	if err := queryChildReferences(id, projectId, &children); err != nil {
		return nil, err
	}

	return append(parents, children...), nil
}

func queryParentReferences(id, projectId string, references *[]QueryReference) error {
	sql := fmt.Sprintf(`
	SELECT DISTINCT ON (structure_name) id, parent_type, child_type, parent_id, name, child_id, child_structure_id, parent_structure_id,
	COALESCE(
  		(SELECT name FROM declarations.maps WHERE id = child_structure_id AND project_id = ?),
  		(SELECT name FROM declarations.lists WHERE id = child_structure_id AND project_id = ?)
 	)AS structure_name
	FROM %s WHERE parent_id = ? AND project_id = ?
`, (declarations.Reference{}).TableName())

	res := storage.Gorm().
		Raw(sql, projectId, projectId, id, projectId).
		Scan(references)

	if res.Error != nil {
		return appErrors.NewDatabaseError(res.Error).AddError("queryMapVariable.Logic", nil)
	}

	return nil
}

func queryChildReferences(id, projectId string, references *[]QueryReference) error {
	sql := fmt.Sprintf(`
	SELECT DISTINCT ON (structure_name) id, parent_type, child_type, parent_id, name, child_id, parent_structure_id, child_structure_id,
	COALESCE(
  		(SELECT name FROM declarations.maps WHERE id = parent_structure_id AND project_id = ?),
  		(SELECT name FROM declarations.lists WHERE id = parent_structure_id AND project_id = ?)
 	) AS structure_name
	FROM %s
	WHERE child_id = ? AND project_id = ?
`, (declarations.Reference{}).TableName())

	res := storage.Gorm().
		Raw(sql, projectId, projectId, id, projectId).
		Scan(references)

	if res.Error != nil {
		return appErrors.NewDatabaseError(res.Error).AddError("queryMapVariable.Logic", nil)
	}

	return nil
}
