package pagination

import (
	"errors"
)

func ResolveCursor(direction, previousRequestNextId, previousRequestPrevId string, results []string, limit int) (string, string, error) {
	// first query, have results
	if direction == DIRECTION_FORWARD && previousRequestNextId == "" && previousRequestPrevId == "" {
		if len(results) >= limit {
			return results[len(results)-1], "", nil
		}

		// back to the beginning
		if len(results) < limit {
			return "", "", nil
		}
	}

	if direction == DIRECTION_FORWARD && previousRequestNextId != "" && previousRequestPrevId == "" {
		if len(results) == limit {
			return results[len(results)-1], "", nil
		}

		if len(results) < limit {
			return "", "", nil
		}
	}

	if direction == DIRECTION_FORWARD && previousRequestNextId != "" && previousRequestPrevId != "" {
		if len(results) == limit {
			return results[len(results)-1], results[0], nil
		}

		if len(results) < limit {
			prevId := previousRequestPrevId
			if len(results) != 0 {
				prevId = results[0]
			}

			return "", prevId, nil
		}
	}

	return "", "", errors.New("Cannot determine next cursor. This is a bug that should be fixed!")
}
