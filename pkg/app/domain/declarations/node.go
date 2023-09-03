package declarations

import (
	"creatif/pkg/app/domain"
	"fmt"
	"github.com/lib/pq"
	"github.com/segmentio/ksuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"time"
)

type Validation struct {
	Required    bool
	Length      string
	ExactValue  string
	ExactValues string
	IsDate      bool
}

type Node struct {
	ID ksuid.KSUID `gorm:"primarykey;type:text CHECK(length(id)=27)"`

	Name       string         `gorm:"index;uniqueIndex:unique_node"`
	Behaviour  string         // readonly,modifiable
	Groups     pq.StringArray `gorm:"type:text[]"` // if groups is set, group should be invalidated
	Metadata   datatypes.JSON
	Validation datatypes.JSONType[Validation]

	// TODO: change this to be ksuid.KSUID when projects and exploration are over, project must exist and be UUID
	/*	ProjectID *string `gorm:"type:uuid;uniqueIndex:unique_node"`
		Project   domain.Project*/

	CreatedAt time.Time `gorm:"<-:create;index"`
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (u *Node) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = ksuid.New()

	return
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
