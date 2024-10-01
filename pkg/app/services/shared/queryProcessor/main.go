package queryProcessor

import (
	"errors"
	"fmt"
	"strconv"
)

func CreateSql(queries []Query) (string, error) {
	sql := ""
	for _, q := range queries {
		if q.Operator == "equal" {
			return fmt.Sprintf("(lv.value->>'%s') = '%s'", q.Column, q.Value), nil
		}

		if q.Operator == "unequal" {
			return fmt.Sprintf("(lv.value->>'%s') != '%s'", q.Column, q.Value), nil
		}

		if q.Operator == "greaterThan" {
			if q.Type == "int" {
				v, err := strconv.ParseInt(q.Value, 10, 64)
				if err != nil {
					return "", errors.New("invalid data type. Expected an integer, got something else")
				}

				return fmt.Sprintf("(lv.value->>'%s') != %d", q.Column, v), nil
			}

			if q.Type == "float" {
				v, err := strconv.ParseInt(q.Value, 10, 64)
				if err != nil {
					return "", errors.New("invalid data type. Expected an integer, got something else")
				}

				return fmt.Sprintf("(lv.value->>'%s') != %d", q.Column, v), nil
			}

			return fmt.Sprintf("(lv.value->>'%s') != %s", q.Column, q.Value), nil
		}
	}

	return sql, nil
}
