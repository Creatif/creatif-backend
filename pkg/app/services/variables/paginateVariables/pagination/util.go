package pagination

import (
	"errors"
)

func ResolveCursor(previousPaginationId, direction, orderBy string, results []string, limit int) (string, error) {
	if direction == DIRECTION_FORWARD {
		if len(results) == 0 {
			return previousPaginationId, nil
		}

		if len(results) < limit {
			return previousPaginationId, nil
		}

		return results[len(results)-1], nil
	}

	if direction == DIRECTION_BACKWARDS {
		if len(results) == 0 {
			return "", nil
		}

		return results[0], nil
	}

	return "", errors.New("Cannot determine next cursor. This is a bug that should be fixed!")
}
