package declarations

import (
	"creatif/pkg/app/domain"
	"creatif/pkg/lib/sdk"
	"fmt"
	"github.com/lib/pq"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"time"
)

type MapNode struct {
	ID string `gorm:"primarykey;type:text CHECK(length(id)=26)"`

	Name      string `gorm:"index;uniqueIndex:unique_node"`
	Behaviour string
	Groups    pq.StringArray `gorm:"type:text[]"`
	Metadata  datatypes.JSON

	MapID string `gorm:"type:text;check:length(id)=26"`
	Map   Map    `gorm:"foreignKey:MapID"`

	// TODO: change this to be string when projects and exploration are over, project must exist and be UUID
	/*	ProjectID *string `gorm:"type:uuid;uniqueIndex:unique_node"`
		Project   domain.Project*/

	CreatedAt time.Time `gorm:"<-:create;index"`
	UpdatedAt time.Time
}

func NewMapNode(name, behaviour string, metadata datatypes.JSON, groups pq.StringArray) MapNode {
	return MapNode{
		Name:      name,
		Behaviour: behaviour,
		Metadata:  metadata,
		Groups:    groups,
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
