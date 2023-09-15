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

type Node struct {
	ID string `gorm:"primarykey;type:text CHECK(length(id)=26)"`

	Name      string `gorm:"index;uniqueIndex:unique_node"`
	Behaviour string
	Groups    pq.StringArray `gorm:"type:text[]"`
	Metadata  datatypes.JSON

	CreatedAt time.Time `gorm:"<-:createNode;index"`
	UpdatedAt time.Time
}

func (u *Node) BeforeCreate(tx *gorm.DB) (err error) {
	id, err := sdk.NewULID()
	if err != nil {
		return err
	}

	u.ID = id

	return nil
}

func (Node) TableName() string {
	return fmt.Sprintf("%s.%s", "declarations", domain.DECLARATION_NODES_TABLE)
}

func NewNode(name, behaviour string, groups []string, metadata []byte) Node {
	return Node{
		Name:      name,
		Groups:    groups,
		Behaviour: behaviour,
		Metadata:  metadata,
	}
}
