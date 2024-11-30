package queryMapVariable

import (
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/lib/appErrors"
	"creatif/pkg/lib/storage"
	"errors"
	"fmt"
	"github.com/lib/pq"
	"gorm.io/datatypes"
	"time"
)

type ChildConnectionStructure struct {
	StructureName string `gorm:"column:structure_name"`
	StructureID   string `gorm:"column:structure_id"`
	StructureType string `gorm:"column:type"`
}

type QueryVariable struct {
	ID      string
	ShortID string
	Index   float64

	Name      string
	Behaviour string
	Metadata  datatypes.JSON
	Value     datatypes.JSON
	Groups    pq.StringArray `gorm:"type:text[];column:groups"`

	MapID    string
	LocaleID string

	CreatedAt time.Time `gorm:"<-:create;index"`
	UpdatedAt time.Time
}

func getVariable(projectId, structureId, variableId string) (QueryVariable, error) {
	sql := fmt.Sprintf(`
			SELECT 
			    lv.id, 
			    lv.name, 
			    lv.index, 
			    lv.behaviour, 
			    lv.short_id, 
			    lv.metadata, 
			       ARRAY((SELECT g.name FROM %s AS vg INNER JOIN %s AS g ON vg.variable_id = lv.id AND g.id = ANY(vg.groups))) AS groups,
				lv.value, 
			    lv.created_at, 
			    lv.updated_at, 
			    lv.locale_id
			FROM %s AS l INNER JOIN %s AS lv
			ON l.project_id = ? AND l.id = ? AND lv.map_id = l.id AND lv.id = ?`,
		(declarations.VariableGroup{}).TableName(), (declarations.Group{}).TableName(), (declarations.Map{}).TableName(), (declarations.MapVariable{}).TableName())

	var variable QueryVariable
	res := storage.Gorm().
		Raw(sql, projectId, structureId, variableId).
		Scan(&variable)

	if res.Error != nil {
		return QueryVariable{}, appErrors.NewDatabaseError(res.Error).AddError("queryMapVariable.Logic", nil)
	}

	if res.RowsAffected == 0 {
		return QueryVariable{}, appErrors.NewNotFoundError(errors.New("No rows found")).AddError("queryMapVariable.Logic", nil)
	}

	return variable, nil
}

func getViewStructuresByVariableFromConnections(variableId string) ([]ChildConnectionStructure, error) {
	var mapChildConnectionStructure []ChildConnectionStructure
	sql := fmt.Sprintf(`
SELECT DISTINCT ON(l.id) l.id AS structure_id, l.name as structure_name, 'map' as "type" FROM %s AS l
    INNER JOIN %s AS lv ON l.id = lv.map_id
    INNER JOIN %s AS conn ON 
        conn.child_variable_id = lv.id AND conn.child_structure_type = 'map' AND 
        conn.parent_variable_id = ?
`,
		(declarations.Map{}).TableName(),
		(declarations.MapVariable{}).TableName(),
		(declarations.Connection{}).TableName())
	res := storage.Gorm().Raw(sql, variableId).Scan(&mapChildConnectionStructure)

	if res.Error != nil {
		return nil, res.Error
	}

	var listChildConnectionStructure []ChildConnectionStructure
	sql = fmt.Sprintf(`
SELECT DISTINCT ON(l.id) l.id AS structure_id, l.name as structure_name, 'list' as "type" FROM %s AS l
    INNER JOIN %s AS lv ON l.id = lv.list_id
    INNER JOIN %s AS conn ON 
        conn.child_variable_id = lv.id AND conn.child_structure_type = 'list' AND 
        conn.parent_variable_id = ?
`,
		(declarations.List{}).TableName(),
		(declarations.ListVariable{}).TableName(),
		(declarations.Connection{}).TableName())
	res = storage.Gorm().Raw(sql, variableId).Scan(&listChildConnectionStructure)

	if res.Error != nil {
		return nil, res.Error
	}

	finalProduct := make([]ChildConnectionStructure, 0, len(mapChildConnectionStructure)+len(listChildConnectionStructure))
	finalProduct = append(finalProduct, mapChildConnectionStructure...)
	finalProduct = append(finalProduct, listChildConnectionStructure...)

	return finalProduct, nil
}
