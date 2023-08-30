package declarations

import (
	"creatif/pkg/app/domain"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Map struct {
	ID uuid.UUID `gorm:"primarykey;type:uuid"`

	Name string `gorm:"uniqueIndex"`

	CreatedAt time.Time `gorm:"<-:create"`
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func NewMap(name string) Map {
	return Map{
		Name: name,
	}
}

func (u *Map) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New()

	return
}

func (Map) TableName() string {
	return fmt.Sprintf("%s.%s", "declarations", domain.NODE_MAP_TABLE)
}
