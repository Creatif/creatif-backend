package pagination

type initialModel struct {
	ID string `gorm:"primarykey;type:text CHECK(length(id)=26)"`
}

type PaginationInfo struct {
	Next    string
	Prev    string
	NextURL string
	PrevURL string
}

func newPaginationInfo(nextCursor, prevCursor, nextURL, prevURL string) PaginationInfo {
	return PaginationInfo{
		Next:    nextCursor,
		Prev:    prevCursor,
		NextURL: nextURL,
		PrevURL: prevURL,
	}
}
