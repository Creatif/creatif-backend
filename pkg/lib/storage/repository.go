package storage

import (
	"creatif/pkg/lib/appErrors"
	"database/sql"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func Create[T any](table string, model T, upsert bool) error {
	g := Gorm().Table(table)
	if upsert {
		g = g.Clauses(clause.OnConflict{DoNothing: true})
	}

	res := g.Create(model)
	if res.Error != nil {
		return res.Error
	}

	return nil
}

func Update[T any](table string, model T) error {
	if res := Gorm().Table(table).Save(model); res.Error != nil {
		return res.Error
	}

	return nil
}

func Get[T any](table, ID string, model T, sel ...string) error {
	if res := Gorm().
		Table(table).
		First(model, "id = ?", ID).
		Select(sel); res.Error != nil {
		return res.Error
	}

	return nil
}

func GetBy[T any](table string, field string, value interface{}, model T, sel ...string) error {
	if res := Gorm().
		Table(table).
		First(model, fmt.Sprintf("%s = ?", field), value).
		Select(sel); res.Error != nil {
		return res.Error
	}

	return nil
}

func GetAll[T any](table string, model T) error {
	res := Gorm().Table(table).Find(model)
	if res.Error != nil {
		return res.Error
	}

	return nil
}

func Find[T any](table string, fn func(db *gorm.DB) (T, error)) (T, error) {
	return fn(Gorm().Table(table))
}

func Delete(table string, model interface{}) error {
	if res := Gorm().Table(table).Delete(model); res.Error != nil {
		return res.Error
	}

	return nil
}

func DeleteBy(table string, field string, value interface{}, model interface{}) error {
	if res := Gorm().Table(table).Where(fmt.Sprintf("%s = ?", field), value).Delete(model); res.Error != nil {
		return res.Error
	}

	return nil
}

func Transaction(fn func(tx *gorm.DB) error) error {
	tx := Gorm().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	if err := fn(tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func Gorm() *gorm.DB {
	return gormHandle
}

func SQLDB() (*sql.DB, error) {
	sql, err := Gorm().DB()
	if err != nil {
		return nil, appErrors.NewDatabaseError(err)
	}

	return sql, nil
}
