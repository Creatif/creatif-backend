package declarations

import (
	"creatif/pkg/app/domain"
	"creatif/pkg/lib/storage"
	"fmt"
	"github.com/segmentio/ksuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"time"
)

type ListVariable struct {
	ID      string  `gorm:"primarykey;type:text"`
	ShortID string  `gorm:"uniqueIndex:unique_variable;type:text;not null"`
	Index   float64 `gorm:"type:float"`

	Name      string         `gorm:"uniqueIndex:unique_variable;not null"`
	Behaviour string         `gorm:"not null"`
	Metadata  datatypes.JSON `gorm:"type:jsonb"`
	Value     datatypes.JSON `gorm:"type:jsonb"`

	LocaleID string `gorm:"uniqueIndex:unique_variable;type:text"`
	ListID   string `gorm:"uniqueIndex:unique_variable;type:text"`
	List     List   `gorm:"foreignKey:ListID"`

	CreatedAt time.Time `gorm:"<-:create;index"`
	UpdatedAt time.Time
}

func NewListVariable(listId, localeID, name, behaviour string, metadata datatypes.JSON, value datatypes.JSON) ListVariable {
	return ListVariable{
		ID:        ksuid.New().String(),
		ListID:    listId,
		LocaleID:  localeID,
		Name:      name,
		Behaviour: behaviour,
		Metadata:  metadata,
		Value:     value,
	}
}

func (u *ListVariable) BeforeCreate(tx *gorm.DB) (err error) {
	shortId, err := storage.ShortId.Generate()
	if err != nil {
		return err
	}
	u.ShortID = shortId

	return nil
}

func (ListVariable) TableName() string {
	return fmt.Sprintf("%s.%s", "declarations", domain.LIST_VARIABLES_TABLE)
}
