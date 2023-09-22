package pagination

type initialModel struct {
	ID string `gorm:"primarykey;type:text CHECK(length(id)=26)"`
}

type Parameters struct {
	Field   string
	OrderBy string
	Groups  []string
	Limit   int
}

type PaginationInfo struct {
	Next       string
	Prev       string
	Parameters Parameters
}
