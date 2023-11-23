package auth

import (
	"creatif/pkg/app/domain/app"
	"creatif/pkg/lib/logger"
	"creatif/pkg/lib/storage"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type frontendAuthentication struct {
	session    string
	logBuilder logger.LogBuilder
	user       AuthenticatedUser
	key        [32]byte
	doRefresh  bool
}

func (a *frontendAuthentication) Authenticate() error {
	if a.session == "" {
		return errors.New("Authentication cookie not provided.")
	}

	jsonToken, err := base64.StdEncoding.DecodeString(a.session)
	if err != nil {
		a.logBuilder.Add("authentication.base64DecodeSession", err.Error())
		return errors.New("Unauthenticated")
	}

	var session AuthenticatedFrontendSession
	if err := json.Unmarshal(jsonToken, &session); err != nil {
		a.logBuilder.Add("authentication.sessionDecode", err.Error())
		return errors.New("Unauthenticated")
	}

	if session.Type != "frontend" {
		return errors.New("Unauthenticated")
	}

	var user app.User
	if res := storage.Gorm().Where("id = ?", session.ID).Select("key", "email").First(&user); res.Error != nil {
		a.logBuilder.Add("authentication.userNotFound", res.Error.Error())
		return errors.New("Unauthenticated")
	}

	encrypedUser, err := base64.StdEncoding.DecodeString(session.Token)
	if err != nil {
		a.logBuilder.Add("authentication.base64DecodeToken", err.Error())
		return errors.New("Unauthenticated")
	}

	var key [32]byte
	for i, v := range user.Key {
		key[i] = byte(v)
	}

	jsonAuthenticatedUser, err := decrypt(encrypedUser, &key)
	if err != nil {
		a.logBuilder.Add("authentication.decryptUser", err.Error())
		return errors.New("Unauthenticated")
	}

	var authenticatedUser AuthenticatedUser
	if err := json.Unmarshal(jsonAuthenticatedUser, &authenticatedUser); err != nil {
		a.logBuilder.Add("authentication.tokenDecode", err.Error())
		return errors.New("Unauthenticated")
	}

	if authenticatedUser.Email != user.Email {
		a.logBuilder.Add("authentication.differentEmails", fmt.Sprintf("Provided email %s did not match the user email %s", authenticatedUser.Email, user.Email))
		return errors.New("Unauthenticated")
	}

	refresh := authenticatedUser.Refresh
	if time.Now().After(refresh.Add(1 * time.Hour)) {
		return errors.New("Unauthenticated")
	}

	if time.Now().After(refresh.Add(45 * time.Minute)) {
		authenticatedUser.Refresh = time.Now()
		a.doRefresh = true
	}

	a.key = key
	a.user = authenticatedUser

	return nil
}

func (a *frontendAuthentication) User() AuthenticatedUser {
	return a.user
}

func (a *frontendAuthentication) ShouldRefresh() bool {
	return a.doRefresh
}

func (a *frontendAuthentication) Logout(cb func()) {
	cb()
}

func (a *frontendAuthentication) Refresh() (string, error) {
	if a.doRefresh {
		loginer := NewEmailLogin(a.user, a.key, a.logBuilder)
		return loginer.Login()
	}

	return a.session, nil
}

func NewFrontendAuthentication(session string, builder logger.LogBuilder) Authentication {
	return &frontendAuthentication{
		session:    session,
		doRefresh:  false,
		logBuilder: builder,
	}
}
