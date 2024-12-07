package pagination

import (
	"context"
	"creatif/pkg/lib/storage"
	"time"
)

type concurrentResult[Result any] struct {
	result Result
	error  error
}

func runQueriesConcurrently(
	structureType string,
	queryPlaceholders map[string]interface{},
	countPlaceholders map[string]interface{},
	sq subQueries,
	defs defaults,
) (concurrentResult[[]QueryVariable], concurrentResult[int64]) {
	paginationCtx, paginationCancel := context.WithTimeout(context.Background(), 30*time.Second)
	countCtx, countCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer paginationCancel()
	defer countCancel()

	paginationChan := make(chan concurrentResult[[]QueryVariable])
	countChan := make(chan concurrentResult[int64])

	go func() {
		paginationSql := createPaginationSql(structureType, sq, defs)
		var items []QueryVariable
		res := storage.Gorm().WithContext(paginationCtx).Raw(paginationSql, queryPlaceholders).Scan(&items)
		if res.Error != nil {
			paginationChan <- concurrentResult[[]QueryVariable]{
				result: nil,
				error:  res.Error,
			}

			countCancel()
			paginationCancel()

			return
		}

		paginationChan <- concurrentResult[[]QueryVariable]{
			result: items,
			error:  nil,
		}
	}()

	go func() {
		countSql := createCountSql(structureType, sq)
		var count int64
		res := storage.Gorm().WithContext(countCtx).Raw(countSql, countPlaceholders).Scan(&count)
		if res.Error != nil {
			countChan <- concurrentResult[int64]{
				result: 0,
				error:  res.Error,
			}

			paginationCancel()
			countCancel()

			return
		}

		countChan <- concurrentResult[int64]{
			result: count,
			error:  nil,
		}
	}()

	paginationResult := <-paginationChan
	countResult := <-countChan

	return paginationResult, countResult
}
