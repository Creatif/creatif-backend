package loginApi

import (
	"context"
	"creatif/pkg/app/auth"
	"creatif/pkg/app/domain/app"
	pkg "creatif/pkg/lib"
	"creatif/pkg/lib/appErrors"
	"creatif/pkg/lib/cache"
	"creatif/pkg/lib/logger"
	"creatif/pkg/lib/storage"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"strings"
	"time"
)

type Main struct {
	model      Model
	logBuilder logger.LogBuilder
	auth       auth.Authentication
}

func (c Main) Validate() error {
	defer func() {
		delCtx, cancel := context.WithDeadline(context.Background(), time.Now().Add(5*time.Second))
		defer cancel()
		_, err := cache.Cache().Del(delCtx, c.model.Session).Result()
		if err != nil {
			c.logBuilder.Add("loginApi", fmt.Sprintf("Session cache could not be deleted: %s", err.Error()))
		}
	}()

	c.logBuilder.Add("loginApi", "Validating...")
	if errs := c.model.Validate(); errs != nil {
		return appErrors.NewValidationError(errs)
	}

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(5*time.Second))
	defer cancel()
	key, err := cache.Cache().Get(ctx, c.model.Session).Result()
	if err != nil {
		return appErrors.NewAuthenticationError(errors.New("Unauthenticated"))
	}

	split := strings.Split(key, "-")
	if len(split) != 2 {
		return appErrors.NewAuthenticationError(errors.New("Unauthenticated"))
	}

	cacheApiKey := split[0]
	cacheProjectId := split[1]

	if cacheApiKey != c.model.ApiKey || cacheProjectId != c.model.ProjectID {
		return appErrors.NewAuthenticationError(errors.New("Unauthenticated"))
	}

	c.logBuilder.Add("loginApi", "Validated.")
	return nil
}

func (c Main) Authenticate() error {
	return nil
}

func (c Main) Authorize() error {
	return nil
}

func (c Main) Logic() (string, error) {
	var user app.User
	res := storage.Gorm().Raw(fmt.Sprintf(`
SELECT 
	u.id,
	u.key,
	u.confirmed,
	u.name,
	u.password,
	u.last_name,
	u.email,
	u.created_at,
	u.updated_at
FROM %s AS u
INNER JOIN %s AS p
ON p.user_id = u.id AND p.api_key = ? AND p.id = ? AND u.email = ?
`, (app.User{}).TableName(), (app.Project{}).TableName()), c.model.ApiKey, c.model.ProjectID, c.model.Email).Scan(&user)

	if res.Error != nil {
		c.logBuilder.Add("apiLogin.getUser", res.Error.Error())
		return "", appErrors.NewAuthenticationError(errors.New("Unauthenticated"))
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(c.model.Password))
	if err != nil {
		return "", appErrors.NewAuthenticationError(err)
	}

	if !user.Confirmed {
		c.logBuilder.Add("apiLogin.notConfirmed", "User not confirmed")
		return "", appErrors.NewAuthenticationError(errors.New("The user is not confirmed"))
	}

	var key [32]byte
	for i, v := range user.Key {
		key[i] = byte(v)
	}

	var project app.Project
	res = storage.Gorm().Where("user_id = ? AND api_key = ?", user.ID, c.model.ApiKey).Select("ID", "api_key").First(&project)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return "", appErrors.NewAuthenticationError(errors.New("Associated project does not exist"))
	}

	if res.Error != nil {
		return "", appErrors.NewAuthenticationError(res.Error)
	}

	if project.ID != c.model.ProjectID {
		return "", appErrors.NewAuthenticationError(errors.New("Invalid project ID."))
	}

	if project.APIKey != c.model.ApiKey {
		return "", appErrors.NewAuthenticationError(errors.New("Invalid API key."))
	}

	authenticatedUser := auth.NewAuthenticatedUser(user.ID, user.Name, user.LastName, user.Email, user.CreatedAt, user.UpdatedAt, time.Now(), project.ID)
	return auth.NewApiLogin(authenticatedUser, key, c.logBuilder).Login()
}

func (c Main) Handle() (string, error) {
	if err := c.Validate(); err != nil {
		return "", err
	}

	if err := c.Authenticate(); err != nil {
		return "", err
	}

	if err := c.Authorize(); err != nil {
		return "", err
	}

	model, err := c.Logic()

	if err != nil {
		return "", err
	}

	return model, nil
}

func New(model Model, auth auth.Authentication, logBuilder logger.LogBuilder) pkg.Job[Model, string, string] {
	logBuilder.Add("loginApi", "Created")
	return Main{model: model, logBuilder: logBuilder, auth: auth}
}
