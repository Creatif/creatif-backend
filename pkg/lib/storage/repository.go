package storage

import (
	"creatif/pkg/lib/appErrors"
	"database/sql"
	"gorm.io/gorm"
)

func Create[T any](model T) error {
	res := Gorm().Create(model)
	if res.Error != nil {
		return res.Error
	}

	return nil
}

func Update[T any](model T) error {
	if res := Gorm().Save(model); res.Error != nil {
		return res.Error
	}

	return nil
}

func Get[T any](ID string, model T, sel ...string) error {
	if res := Gorm().
		First(model, "id = ?", ID).
		Select(sel); res.Error != nil {
		return res.Error
	}

	return nil
}

func GetAll[T any](model T) error {
	res := Gorm().Find(model)
	if res.Error != nil {
		return res.Error
	}

	return nil
}

func Find[T any](fn func(db *gorm.DB) (T, error)) (T, error) {
	return fn(Gorm())
}

func Delete(model interface{}) error {
	if res := Gorm().Delete(model); res.Error != nil {
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
