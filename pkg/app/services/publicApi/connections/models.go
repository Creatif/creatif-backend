package connections

type Model struct {
	Parent string `gorm:"column:parent"`
	Child  string `gorm:"column:child"`
}
