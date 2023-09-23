package declarations

import (
	"creatif/pkg/app/domain"
	"creatif/pkg/lib/sdk"
	"fmt"
	"github.com/lib/pq"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"time"
)

type Variable struct {
	ID string `gorm:"primarykey;type:text CHECK(length(id)=26)"`

	Name      string `gorm:"index;uniqueIndex:unique_variable"`
	Behaviour string
	Groups    pq.StringArray `gorm:"type:text[]"`
	Metadata  datatypes.JSON `gorm:"type:jsonb"`
	Value     datatypes.JSON `gorm:"type:jsonb"`

	ProjectID string `gorm:"index;type:text;check:length(id)=26"`

	CreatedAt time.Time `gorm:"<-:create;index"`
	UpdatedAt time.Time
}

func (u *Variable) BeforeCreate(tx *gorm.DB) (err error) {
	id, err := sdk.NewULID()
	if err != nil {
		return err
	}

	u.ID = id

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
