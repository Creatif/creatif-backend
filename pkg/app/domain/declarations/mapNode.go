package declarations

import (
	"creatif/pkg/app/domain"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type MapNode struct {
	ID uuid.UUID `gorm:"primarykey;type:uuid"`

	NodeID uuid.UUID `gorm:"type:uuid"`
	Node   Node      `gorm:"foreignKey:NodeID"`

	MapID uuid.UUID `gorm:"type:uuid"`
	Map   Map       `gorm:"foreignKey:MapID"`

	CreatedAt time.Time `gorm:"<-:create"`
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func NewMapNode(nodeId uuid.UUID) MapNode {
	return MapNode{
		NodeID: nodeId,
	}
}

func (u *MapNode) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New()

	return
}

func (MapNode) TableName() string {
	return fmt.Sprintf("%s.%s", "declarations", domain.NODE_MAP_NODES_TABLE)
}
