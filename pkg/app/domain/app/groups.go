package app

import (
	"creatif/pkg/app/domain"
	"fmt"
	"time"
)

type Group struct {
	ID string `gorm:"primarykey;type:text;default:gen_ulid()"`

	Name      string `gorm:"uniqueIndex:group_name;not null"`
	ProjectID string `gorm:"uniqueIndex:group_name;type:text"`

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
	return fmt.Sprintf("%s.%s", "app", domain.GROUPS_TABLE)
}
