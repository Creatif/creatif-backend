package declarations

import (
	"creatif/pkg/app/domain"
	"creatif/pkg/lib/storage"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type List struct {
	ID      string `gorm:"primarykey;type:text;default:gen_ulid()"`
	ShortID string `gorm:"uniqueIndex:unique_list;type:text"`

	Name          string         `gorm:"uniqueIndex:unique_list_name"`
	ProjectID     string         `gorm:"uniqueIndex:unique_list_name;type:text"`
	LocaleID      string         `gorm:"uniqueIndex:unique_list_name;type:text"`
	ListVariables []ListVariable `gorm:"foreignKey:ListID;constraint:OnDelete:CASCADE;"`

	CreatedAt time.Time `gorm:"<-:create"`
	UpdatedAt time.Time
}

func NewList(projectId, name, localeID string) List {
	return List{
		Name:      name,
		LocaleID:  localeID,
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
