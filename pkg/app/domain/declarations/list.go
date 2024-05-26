package declarations

import (
	"creatif/pkg/app/domain"
	"creatif/pkg/lib/storage"
	"fmt"
	"github.com/segmentio/ksuid"
	"gorm.io/gorm"
	"time"
)

type List struct {
	ID      string `gorm:"primarykey;type:text"`
	ShortID string `gorm:"uniqueIndex:unique_list;type:text;not null"`

	Name          string         `gorm:"uniqueIndex:unique_list_name;not null"`
	ProjectID     string         `gorm:"uniqueIndex:unique_list_name;type:text"`
	ListVariables []ListVariable `gorm:"foreignKey:ListID;constraint:OnDelete:CASCADE;"`

	Serial int64 `gorm:"default: 0"`

	CreatedAt time.Time `gorm:"<-:create"`
	UpdatedAt time.Time
}

func NewList(projectId, name string) List {
	return List{
		ID:        ksuid.New().String(),
		Name:      name,
		ProjectID: projectId,
	}
}

func (u *List) BeforeCreate(tx *gorm.DB) (err error) {
	shortId, err := storage.ShortId.Generate()
	if err != nil {
		return err
	}
	u.ShortID = shortId

	return nil
}

func (List) TableName() string {
	return fmt.Sprintf("%s.%s", "declarations", domain.LIST_TABLE)
}
