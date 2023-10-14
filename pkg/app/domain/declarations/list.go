package declarations

import (
	"creatif/pkg/app/domain"
	"creatif/pkg/lib/sdk"
	"creatif/pkg/lib/storage"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type List struct {
	ID      string `gorm:"primarykey;type:text CHECK(length(id)=26)"`
	ShortID string `gorm:"uniqueIndex:unique_list;type:text"`

	Name          string         `gorm:"uniqueIndex:unique_list_name"`
	ProjectID     string         `gorm:"uniqueIndex:unique_list_name;type:text;check:length(id)=26"`
	LocaleID      string         `gorm:"uniqueIndex:unique_list_name;type:text;check:length(id)=26;not null"`
	ListVariables []ListVariable `gorm:"foreignKey:ListID;constraint:OnDelete:CASCADE;"`

	CreatedAt time.Time `gorm:"<-:create"`
	UpdatedAt time.Time
}

func NewList(projectId, name string) List {
	return List{
		Name:      name,
		ProjectID: projectId,
	}
}

func (u *List) BeforeCreate(tx *gorm.DB) (err error) {
	id, err := sdk.NewULID()
	if err != nil {
		return err
	}

	u.ID = id
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
