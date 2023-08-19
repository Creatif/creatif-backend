package declarations

import (
	"creatif/pkg/app/domain"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Map struct {
	ID string `gorm:"primarykey"`

	NodeID string
	Node   Node `gorm:"foreignKey:NodeID"`

	CreatedAt time.Time `gorm:"<-:create"`
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func NewMap(nodeId string) Map {
	return Map{
		NodeID: nodeId,
	}
}

func (u *Map) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New().String()

	return
}

func (Map) TableName() string {
	return fmt.Sprintf("%s.%s", "declarations", domain.NODE_MAPS_TABLE)
}
