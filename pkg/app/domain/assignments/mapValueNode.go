package assignments

import (
	"creatif/pkg/app/domain"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type MapValueNode struct {
	ID               uuid.UUID `gorm:"primarykey;type:uuid"`
	AssignmentNodeID uuid.UUID `gorm:"type:uuid"`

	Value datatypes.JSON
}

func NewMapValueNode(assignmentNodeID uuid.UUID, value datatypes.JSON) MapValueNode {
	return MapValueNode{
		Value:            value,
		AssignmentNodeID: assignmentNodeID,
	}
}

func (u *MapValueNode) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New()

	return
}

func (MapValueNode) TableName() string {
	return fmt.Sprintf("%s.%s", "assignments", domain.ASSIGNMENT_MAP_VALUE_NODE)
}
