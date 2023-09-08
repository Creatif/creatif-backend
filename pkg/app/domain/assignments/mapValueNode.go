package assignments

import (
	"creatif/pkg/app/domain"
	"creatif/pkg/lib/sdk"
	"fmt"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type MapValueNode struct {
	ID               string `gorm:"primarykey;type:text CHECK(length(id)=26)"`
	AssignmentNodeID string `gorm:"type:text CHECK(length(assignment_node_id)=26)"`

	Value datatypes.JSON
}

func NewMapValueNode(assignmentNodeID string, value datatypes.JSON) MapValueNode {
	return MapValueNode{
		Value:            value,
		AssignmentNodeID: assignmentNodeID,
	}
}

func (u *MapValueNode) BeforeCreate(tx *gorm.DB) (err error) {
	id, err := sdk.NewULID()
	if err != nil {
		return err
	}

	u.ID = id

	return nil
}

func (MapValueNode) TableName() string {
	return fmt.Sprintf("%s.%s", "assignments", domain.ASSIGNMENT_MAP_VALUE_NODE)
}
