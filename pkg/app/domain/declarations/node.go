package declarations

import (
	"creatif/pkg/app/domain"
	"fmt"
	"github.com/google/uuid"
	"github.com/lib/pq"
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
	ID uuid.UUID `gorm:"primarykey;type:uuid"`

	Name       string         `gorm:"index;uniqueIndex:unique_node"`
	Type       string         // text,image,file,boolean
	Behaviour  string         // readonly,modifiable
	Groups     pq.StringArray `gorm:"type:text[]"` // if groups is set, group should be invalidated
	Metadata   datatypes.JSON
	Validation datatypes.JSONType[Validation]

	// TODO: change this to be uuid.UUID when projects and exploration are over, project must exist and be UUID
	/*	ProjectID *string `gorm:"type:uuid;uniqueIndex:unique_node"`
		Project   domain.Project*/

	CreatedAt time.Time `gorm:"<-:create"`
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (u *Node) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New()

	return
}

func (Node) TableName() string {
	return fmt.Sprintf("%s.%s", "declarations", domain.DECLARATION_NODES_TABLE)
}

func NewNode(name, t, behaviour string, groups []string, metadata []byte) Node {
	return Node{
		Name:      name,
		Type:      t,
		Groups:    groups,
		Behaviour: behaviour,
		Metadata:  metadata,
	}
}
