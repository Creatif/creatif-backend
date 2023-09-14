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

type Validation struct {
	Required    bool
	Length      string
	ExactValue  string
	ExactValues []string
	IsDate      bool
}

type Node struct {
	ID string `gorm:"primarykey;type:text CHECK(length(id)=26)"`

	Name       string `gorm:"index;uniqueIndex:unique_node"`
	Behaviour  string
	Groups     pq.StringArray `gorm:"type:text[]"`
	Metadata   datatypes.JSON
	Validation datatypes.JSONType[Validation]

	// TODO: change this to be string when projects and exploration are over, project must exist and be UUID
	/*	ProjectID *string `gorm:"type:uuid;uniqueIndex:unique_node"`
		Project   domain.Project*/

	CreatedAt time.Time `gorm:"<-:create;index"`
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
