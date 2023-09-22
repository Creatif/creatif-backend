package pagination

type initialModel struct {
	ID string `gorm:"primarykey;type:text CHECK(length(id)=26)"`
}

type Parameters struct {
	NextID    string
	PrevID    string
	Field     string
	OrderBy   string
	Direction string
	Groups    []string
	Limit     int
}

type PaginationInfo struct {
	Next           string
	Prev           string
	NextParameters Parameters
	PrevParameters Parameters
}
