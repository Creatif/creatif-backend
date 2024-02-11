package declarations

import (
	"creatif/pkg/app/domain"
	"fmt"
	"time"
)

type Group struct {
	Name      string `gorm:"primarykey;type:text"`
	ProjectID string `gorm:"index;type:text"`

	VariableGroups []VariableGroup `gorm:"foreignKey:GroupID;references:Name;constraint:OnDelete:CASCADE;"`

	CreatedAt time.Time `gorm:"<-:create"`
	UpdatedAt time.Time `gorm:"<-:update"`
}

func NewGroup(projectId, name string) Group {
	return Group{
		Name:      name,
		ProjectID: projectId,
	}
}

func (Group) TableName() string {
	return fmt.Sprintf("%s.%s", "declarations", domain.GROUPS_TABLE)
}
