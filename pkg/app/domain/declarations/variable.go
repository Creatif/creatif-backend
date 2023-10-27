package declarations

import (
	"creatif/pkg/app/domain"
	"creatif/pkg/lib/storage"
	"fmt"
	"github.com/lib/pq"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"time"
)

type Variable struct {
	ID      string `gorm:"primarykey;type:text;default:gen_ulid()"`
	ShortID string `gorm:"uniqueIndex:unique_short_id;type:text"`

	Name      string `gorm:"uniqueIndex:unique_variable_per_project"`
	Behaviour string
	Groups    pq.StringArray `gorm:"type:text[]"`
	Metadata  datatypes.JSON `gorm:"type:jsonb"`
	Value     datatypes.JSON `gorm:"type:jsonb"`

	ProjectID string `gorm:"uniqueIndex:unique_variable_per_project;type:text"`
	LocaleID  string `gorm:"uniqueIndex:unique_variable_per_project;type:text"`

	CreatedAt time.Time `gorm:"<-:create;index"`
	UpdatedAt time.Time
}

func (u *Variable) BeforeCreate(tx *gorm.DB) (err error) {
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

func NewVariable(projectId, localeID, name, behaviour string, groups []string, metadata, value []byte) Variable {
	if groups == nil || len(groups) == 0 {
		groups = make(pq.StringArray, 0)
	}
	
	return Variable{
		Name:      name,
		LocaleID:  localeID,
		ProjectID: projectId,
		Groups:    groups,
		Behaviour: behaviour,
		Metadata:  metadata,
		Value:     value,
	}
}
