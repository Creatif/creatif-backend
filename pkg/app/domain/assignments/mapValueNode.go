package assignments

import (
	"creatif/pkg/app/domain"
	"fmt"
	"github.com/segmentio/ksuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type MapValueNode struct {
	ID               ksuid.KSUID `gorm:"primarykey;type:text CHECK(length(id)=27)"`
	AssignmentNodeID ksuid.KSUID `gorm:"type:text CHECK(length(assignment_node_id)=27)"`

	Value datatypes.JSON
}

func NewMapValueNode(assignmentNodeID ksuid.KSUID, value datatypes.JSON) MapValueNode {
	return MapValueNode{
		Value:            value,
		AssignmentNodeID: assignmentNodeID,
	}
}

func (u *MapValueNode) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = ksuid.New()

	return
}

func (MapValueNode) TableName() string {
	return fmt.Sprintf("%s.%s", "assignments", domain.ASSIGNMENT_MAP_VALUE_NODE)
}
