package domain

import (
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Node struct {
	ID string `gorm:"primarykey"`

	Name      string `gorm:"index:unique_node"`
	Type      string // text,image,file,boolean
	Group     string `gorm:"index"`
	Behaviour string // readonly,modifiable

	ProjectID *string `gorm:"index:unique_node"`
	Project   Project

	CreatedAt time.Time `gorm:"<-:create"`
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (u *Node) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New().String()

	return
}

func (Node) TableName() string {
	return fmt.Sprintf("%s.%s", "declarations", NODES_TABLE)
}

func NewNode(name, t, group, behaviour string) Node {
	return Node{
		Name:      name,
		Type:      t,
		Group:     group,
		Behaviour: behaviour,
	}
}
