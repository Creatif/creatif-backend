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

type Variable struct {
	ID      string `gorm:"primarykey;type:text CHECK(length(id)=26)"`
	ShortID string `gorm:"uniqueIndex:unique_short_id;type:text"`

	Name      string `gorm:"uniqueIndex:unique_variable_per_project"`
	Behaviour string
	Groups    pq.StringArray `gorm:"type:text[]"`
	Metadata  datatypes.JSON `gorm:"type:jsonb"`
	Value     datatypes.JSON `gorm:"type:jsonb"`

	ProjectID string `gorm:"uniqueIndex:unique_variable_per_project;type:text;check:length(id)=26;not null;default: null"`

	CreatedAt time.Time `gorm:"<-:create;index"`
	UpdatedAt time.Time
}

func (u *Variable) BeforeCreate(tx *gorm.DB) (err error) {
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

func (Variable) TableName() string {
	return fmt.Sprintf("%s.%s", "declarations", domain.VARIABLES_TABLE)
}

func NewVariable(projectId, name, behaviour string, groups []string, metadata, value []byte) Variable {
	return Variable{
		Name:      name,
		ProjectID: projectId,
		Groups:    groups,
		Behaviour: behaviour,
		Metadata:  metadata,
		Value:     value,
	}
}
