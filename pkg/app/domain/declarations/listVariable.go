package declarations

import (
	"creatif/pkg/app/domain"
	"creatif/pkg/lib/sdk"
	"creatif/pkg/lib/storage"
	"fmt"
	"github.com/lib/pq"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"time"
)

type ListVariable struct {
	ID      string `gorm:"primarykey;type:text;default:gen_ulid()"`
	ShortID string `gorm:"uniqueIndex:unique_variable;type:text"`
	Index   string `gorm:"type:text;uniqueIndex:unique_list_variable"`

	Name      string
	Behaviour string
	Groups    pq.StringArray `gorm:"type:text[]"`
	Metadata  datatypes.JSON `gorm:"type:jsonb"`
	Value     datatypes.JSON `gorm:"type:jsonb"`

	LocaleID string `gorm:"uniqueIndex:unique_list_variable;type:text"`
	ListID   string `gorm:"uniqueIndex:unique_list_variable;type:text"`
	List     List   `gorm:"foreignKey:ListID"`

	CreatedAt time.Time `gorm:"<-:create;index"`
	UpdatedAt time.Time
}

func NewListVariable(listId, localeID, name, behaviour string, metadata datatypes.JSON, groups pq.StringArray, value datatypes.JSON) ListVariable {
	return ListVariable{
		ListID:    listId,
		LocaleID:  localeID,
		Name:      name,
		Behaviour: behaviour,
		Metadata:  metadata,
		Groups:    groups,
		Value:     value,
	}
}

func (u *ListVariable) BeforeCreate(tx *gorm.DB) (err error) {
	shortId, err := storage.ShortId.Generate()
	if err != nil {
		return err
	}
	u.ShortID = shortId

	if u.Index == "" {
		idx, err := sdk.NewULID()
		if err != nil {
			return err
		}
		u.Index = idx
	}

	return nil
}

func (ListVariable) TableName() string {
	return fmt.Sprintf("%s.%s", "declarations", domain.LIST_VARIABLES_TABLE)
}
