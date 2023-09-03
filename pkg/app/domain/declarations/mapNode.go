package declarations

import (
	"creatif/pkg/app/domain"
	"fmt"
	"github.com/segmentio/ksuid"
	"gorm.io/gorm"
	"time"
)

type MapNode struct {
	ID ksuid.KSUID `gorm:"primarykey;type:text CHECK(length(id)=27)"`

	NodeID ksuid.KSUID `gorm:"type:text CHECK(length(node_id)=27)"`
	Node   Node        `gorm:"foreignKey:NodeID"`

	MapID ksuid.KSUID `gorm:"type:text CHECK(length(map_id)=27)"`
	Map   Map         `gorm:"foreignKey:MapID"`

	CreatedAt time.Time `gorm:"<-:create"`
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func NewMapNode(nodeId ksuid.KSUID) MapNode {
	return MapNode{
		NodeID: nodeId,
	}
}

func (u *MapNode) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = ksuid.New()

	return
}

func (MapNode) TableName() string {
	return fmt.Sprintf("%s.%s", "declarations", domain.NODE_MAP_NODES_TABLE)
}
