package declarations

import (
	"creatif/pkg/app/domain"
	"creatif/pkg/lib/sdk"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type MapNode struct {
	ID string `gorm:"primarykey;type:text;check:length(id)=26"`

	NodeID string `gorm:"type:text;check:length(id)=26"`
	Node   Node   `gorm:"foreignKey:NodeID"`

	MapID string `gorm:"type:text;check:length(id)=26"`
	Map   Map    `gorm:"foreignKey:MapID"`

	CreatedAt time.Time `gorm:"<-:create"`
	UpdatedAt time.Time
}

func NewMapNode(nodeId string) MapNode {
	return MapNode{
		NodeID: nodeId,
	}
}

func (u *MapNode) BeforeCreate(tx *gorm.DB) (err error) {
	id, err := sdk.NewULID()
	if err != nil {
		return err
	}

	u.ID = id

	return nil
}

func (MapNode) TableName() string {
	return fmt.Sprintf("%s.%s", "declarations", domain.NODE_MAP_NODES_TABLE)
}
