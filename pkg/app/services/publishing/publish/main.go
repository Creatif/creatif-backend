package publish

import "C"
import (
	"context"
	"creatif/pkg/app/auth"
	"creatif/pkg/app/domain/published"
	pkg "creatif/pkg/lib"
	"creatif/pkg/lib/appErrors"
	"creatif/pkg/lib/logger"
	"creatif/pkg/lib/storage"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Results struct {
	Errors []error
}

type Main struct {
	model      Model
	logBuilder logger.LogBuilder
	auth       auth.Authentication
}

func (c Main) Validate() error {
	c.logBuilder.Add("publish", "Validating...")
	if errs := c.model.Validate(); errs != nil {
		return appErrors.NewValidationError(errs)
	}

	if err := validateVersionNameExists(c.model.ProjectID, c.model.Name); err != nil {
		return err
	}

	c.logBuilder.Add("publish", "Validated")
	return nil
}

func (c Main) Authenticate() error {
	if err := c.auth.Authenticate(); err != nil {
		return appErrors.NewAuthenticationError(err)
	}

	return nil
}

func (c Main) Authorize() error {
	return nil
}

func (c Main) Logic() (published.Version, error) {
	name := c.model.Name
	if name == "" {
		name = uuid.NewString()
	}

	version := published.NewVersion(c.model.ProjectID, name)
	if transactionError := storage.Transaction(func(tx *gorm.DB) error {
		if res := tx.Exec(fmt.Sprintf("DELETE FROM %s WHERE project_id = ?", (published.Version{}).TableName()), c.model.ProjectID); res.Error != nil {
			return res.Error
		}
		fmt.Println("Previous version deleted...")

		if res := tx.Create(&version); res.Error != nil {
			return res.Error
		}
		fmt.Println("Version created...")

		ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
		defer cancel()
		if res := tx.WithContext(ctx).Exec(getMergeSql(version.ID, (published.PublishedList{}).TableName(), getSelectListSql()), c.model.ProjectID); res.Error != nil {
			return res.Error
		}
		fmt.Println("Lists published...")

		ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		var maps []SingleItem
		if res := storage.Gorm().WithContext(ctx).Raw(getSelectMapSql(), c.model.ProjectID).Scan(&maps); res.Error != nil {
			return res.Error
		}

		publishedMaps := make([]published.PublishedMap, 0)
		batches := make([][]published.PublishedMap, 0)
		currentBatchNum := 0
		for _, m := range maps {
			if currentBatchNum == 4500 {
				currentBatchNum = 0
				batches = append(batches, publishedMaps)
				publishedMaps = make([]published.PublishedMap, 0)
			}

			publishedMaps = append(publishedMaps, published.NewPublishedMap(
				m.ID,
				m.ShortID,
				version.ID,
				m.Name,
				m.VariableName,
				m.VariableID,
				m.VariableShortID,
				m.Behaviour,
				m.Locale,
				m.Value,
				m.Groups,
				m.Index,
			))

			currentBatchNum++
		}

		if len(publishedMaps) > 0 {
			batches = append(batches, publishedMaps)
			publishedMaps = nil
		}

		for _, batch := range batches {
			if res := tx.Create(&batch); res.Error != nil {
				return res.Error
			}
		}

		/*		if res := tx.WithContext(ctx).Exec(getReferenceMergeSql(version.ID, getReferencesSql()), c.model.ProjectID); res.Error != nil {
				return res.Error
			}*/
		fmt.Println("References published...")

		return nil
	}); transactionError != nil {
		return published.Version{}, appErrors.NewApplicationError(transactionError)
	}

	return version, nil
}

func (c Main) Handle() (View, error) {
	if err := c.Validate(); err != nil {
		return View{}, err
	}

	if err := c.Authenticate(); err != nil {
		return View{}, err
	}

	if err := c.Authorize(); err != nil {
		return View{}, err
	}

	model, err := c.Logic()

	if err != nil {
		return View{}, err
	}

	return newView(model), nil
}

func New(model Model, auth auth.Authentication, logBuilder logger.LogBuilder) pkg.Job[Model, View, published.Version] {
	logBuilder.Add("publish", "Created")
	return Main{model: model, logBuilder: logBuilder, auth: auth}
}
