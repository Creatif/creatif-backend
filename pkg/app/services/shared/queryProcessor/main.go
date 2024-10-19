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
			if q.Type == "int" {
				v, err := strconv.ParseInt(q.Value, 10, 64)
				if err != nil {
					return "", errors.New("invalid data type. Expected an integer, got something else")
				}

				if sql != "" {
					sql += fmt.Sprintf("AND %s", fmt.Sprintf("CAST((lv.value->>'%s') AS integer) = %d", q.Column, v))
				} else {
					sql += fmt.Sprintf("CAST((lv.value->>'%s') AS integer) = %d", q.Column, v)
				}
			} else if q.Type == "float" {
				v, err := strconv.ParseFloat(q.Value, 64)
				if err != nil {
					return "", errors.New("invalid data type. Expected an integer, got something else")
				}

				if sql != "" {
					sql += fmt.Sprintf("AND %s", fmt.Sprintf("CAST((lv.value->>'%s') AS double precision) = %f", q.Column, v))
				} else {
					sql += fmt.Sprintf("CAST((lv.value->>'%s') AS double precision) = %f", q.Column, v)
				}
			} else if q.Type == "string" {
				if sql != "" {
					sql += fmt.Sprintf("AND %s", fmt.Sprintf("(lv.value->>'%s') = '%s'", q.Column, q.Value))
				} else {
					sql += fmt.Sprintf("(lv.value->>'%s') = '%s' ", q.Column, q.Value)
				}
			}
		}

		if q.Operator == "unequal" {
			if q.Type == "int" {
				v, err := strconv.ParseInt(q.Value, 10, 64)
				if err != nil {
					return "", errors.New("invalid data type. Expected an integer, got something else")
				}

				if sql != "" {
					sql += fmt.Sprintf("AND %s", fmt.Sprintf("CAST((lv.value->>'%s') AS integer) != %d", q.Column, v))
				} else {
					sql += fmt.Sprintf("CAST((lv.value->>'%s') AS integer) != %d", q.Column, v)
				}
			} else if q.Type == "float" {
				v, err := strconv.ParseFloat(q.Value, 64)
				if err != nil {
					return "", errors.New("invalid data type. Expected an integer, got something else")
				}

				if sql != "" {
					sql += fmt.Sprintf("AND %s", fmt.Sprintf("CAST((lv.value->>'%s') AS double precision) != %f", q.Column, v))
				} else {
					sql += fmt.Sprintf("CAST((lv.value->>'%s') AS double precision) != %f", q.Column, v)
				}
			} else if q.Type == "string" {
				if sql != "" {
					sql += fmt.Sprintf("AND %s", fmt.Sprintf("(lv.value->>'%s') != '%s'", q.Column, q.Value))
				} else {
					sql += fmt.Sprintf("(lv.value->>'%s') != '%s' ", q.Column, q.Value)
				}
			}
		}

		/*		if q.Operator == "greaterThan" {
					if q.Type == "int" {
						v, err := strconv.ParseInt(q.Value, 10, 64)
						if err != nil {
							return "", errors.New("invalid data type. Expected an integer, got something else")
						}

						return fmt.Sprintf("(lv.value->>'%s') > %d", q.Column, v), nil
					}

					if q.Type == "float" {
						v, err := strconv.ParseInt(q.Value, 10, 64)
						if err != nil {
							return "", errors.New("invalid data type. Expected an integer, got something else")
						}

						return fmt.Sprintf("(lv.value->>'%s') > %d", q.Column, v), nil
					}

					return fmt.Sprintf("(lv.value->>'%s') > %s", q.Column, q.Value), nil
				}

				if q.Operator == "greaterThanOrEqual" {
					if q.Type == "int" {
						v, err := strconv.ParseInt(q.Value, 10, 64)
						if err != nil {
							return "", errors.New("invalid data type. Expected an integer, got something else")
						}

						return fmt.Sprintf("(lv.value->>'%s') >= %d", q.Column, v), nil
					}

					if q.Type == "float" {
						v, err := strconv.ParseInt(q.Value, 10, 64)
						if err != nil {
							return "", errors.New("invalid data type. Expected an integer, got something else")
						}

						return fmt.Sprintf("(lv.value->>'%s') >= %d", q.Column, v), nil
					}

					return fmt.Sprintf("(lv.value->>'%s') >= %s", q.Column, q.Value), nil
				}

				if q.Operator == "lessThan" {
					if q.Type == "int" {
						v, err := strconv.ParseInt(q.Value, 10, 64)
						if err != nil {
							return "", errors.New("invalid data type. Expected an integer, got something else")
						}

						return fmt.Sprintf("(lv.value->>'%s') < %d", q.Column, v), nil
					}

					if q.Type == "float" {
						v, err := strconv.ParseInt(q.Value, 10, 64)
						if err != nil {
							return "", errors.New("invalid data type. Expected an integer, got something else")
						}

						return fmt.Sprintf("(lv.value->>'%s') < %d", q.Column, v), nil
					}

					return fmt.Sprintf("(lv.value->>'%s') < %s", q.Column, q.Value), nil
				}

				if q.Operator == "lessThanOrEqual" {
					if q.Type == "int" {
						v, err := strconv.ParseInt(q.Value, 10, 64)
						if err != nil {
							return "", errors.New("invalid data type. Expected an integer, got something else")
						}

						return fmt.Sprintf("(lv.value->>'%s') <= %d", q.Column, v), nil
					}

					if q.Type == "float" {
						v, err := strconv.ParseInt(q.Value, 10, 64)
						if err != nil {
							return "", errors.New("invalid data type. Expected an integer, got something else")
						}

						return fmt.Sprintf("(lv.value->>'%s') <= %d", q.Column, v), nil
					}

					return fmt.Sprintf("(lv.value->>'%s') <= %s", q.Column, q.Value), nil
				}*/
	}

	return sql, nil
}
