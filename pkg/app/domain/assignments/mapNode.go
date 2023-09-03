package assignments

import (
	"creatif/pkg/app/domain"
	"creatif/pkg/app/domain/declarations"
	"fmt"
	"github.com/segmentio/ksuid"
	"gorm.io/gorm"
	"time"
)

// text, image, file, json, code, boolean
type MapNode struct {
	ID ksuid.KSUID `gorm:"primarykey;type:text CHECK(length(id)=27)"`

	Name string `gorm:"index;uniqueIndex:unique_node"`

	// TODO: change this to be ksuid.KSUID when projects and exploration are over, project must exist and be UUID
	// TODO: remove pointer since now it allows null
	/*	ProjectID *string `gorm:"uniqueIndex:unique_node"`
		Project   domain.Project*/

	MapNodeID ksuid.KSUID      `gorm:"type:uuid"`
	MapNode   declarations.Map `gorm:"foreignKey:MapNodeID"`

	CreatedAt time.Time `gorm:"<-:create"`
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func NewMapNode(name string, mapNodeID ksuid.KSUID) MapNode {
	return MapNode{
		Name:      name,
		MapNodeID: mapNodeID,
	}
}

func (u *MapNode) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = ksuid.New()

	return
}

func (MapNode) TableName() string {
	return fmt.Sprintf("%s.%s", "assignments", domain.ASSIGNMENT_MAP_NODES_TABLE)
}
