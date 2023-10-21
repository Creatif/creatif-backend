package storage

import (
	"fmt"
	"gorm.io/gorm"
)

func WithIN(field string, values []string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if len(values) == 0 {
			return db
		}

		return db.Where(fmt.Sprintf("%s IN ?", field), values)
	}
}

func WithOffset(offset int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(offset)
	}
}

func WithLimit(limit int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Limit(limit)
	}
}

func OrderBy(field, direction string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Order(fmt.Sprintf("%s %s", field, direction))
	}
}

func WithRegex(field string, value string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if value == "" {
			return db
		}

		return db.Or(db.Where(
			fmt.Sprintf("%s ILIKE ? OR %s ILIKE ? OR %s ILIKE ?", field, field, field),
			fmt.Sprintf("%%%s%%", value),
			fmt.Sprintf("%s%%", value),
			fmt.Sprintf("%%%s", value),
		))
	}
}

func WithGenericString(tp string, val string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("%s = ?", tp), val)
	}
}
