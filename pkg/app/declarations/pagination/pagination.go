package pagination

import "github.com/pilagod/gorm-cursor-paginator/v2/paginator"

func CreatePaginator(sortField, sortOrder string, limit int, sqlRepr string) *paginator.Paginator {
	order := paginator.DESC
	if sortOrder == "asc" {
		order = paginator.ASC
	}
	p := paginator.New(
		&paginator.Config{
			Rules: []paginator.Rule{
				{
					Key: "ID",
				},
				{
					Key:             sortField,
					Order:           order,
					SQLRepr:         sqlRepr,
					NULLReplacement: "1970-01-01",
				},
			},
			After:  "",
			Before: "",
			Limit:  limit,
			// Order here will apply to keys without order specified.
			// In this example paginator will order by "ID" ASC, "JoinedAt" DESC.
			Order: order,
		},
	)
	// ...
	return p
}
