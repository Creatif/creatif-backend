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

type apiAuthentication struct {
	session    string
	logBuilder logger.LogBuilder
	user       AuthenticatedUser
	key        [32]byte
	projectId  string
	apiKey     string
	doRefresh  bool
}

func (a *apiAuthentication) Authenticate() error {
	if a.session == "" {
		return errors.New("Authentication cookie not provided.")
	}

	jsonToken, err := base64.StdEncoding.DecodeString(a.session)
	if err != nil {
		a.logBuilder.Add("authentication.base64DecodeSession", err.Error())
		return errors.New("Failed to decode session cookie")
	}

	var session AuthenticatedApiSession
	if err := json.Unmarshal(jsonToken, &session); err != nil {
		a.logBuilder.Add("authentication.sessionDecode", err.Error())
		return errors.New("Failed to decode API session object")
	}

	if session.Type != "api" {
		return errors.New("Unauthenticated")
	}

	var user app.User
	if res := storage.Gorm().Where("id = ?", session.ID).Select("key", "email").First(&user); res.Error != nil {
		a.logBuilder.Add("authentication.userNotFound", res.Error.Error())
		return errors.New("User not found")
	}

	encrypedUser, err := base64.StdEncoding.DecodeString(session.Token)
	if err != nil {
		a.logBuilder.Add("authentication.base64DecodeToken", err.Error())
		return errors.New("Failed to decode token.")
	}

	var key [32]byte
	for i, v := range user.Key {
		key[i] = byte(v)
	}

	jsonAuthenticatedUser, err := decrypt(encrypedUser, &key)
	if err != nil {
		a.logBuilder.Add("authentication.decryptUser", err.Error())
		return errors.New("Failed to decrypt decoded user.")
	}

	var authenticatedUser AuthenticatedUser
	if err := json.Unmarshal(jsonAuthenticatedUser, &authenticatedUser); err != nil {
		a.logBuilder.Add("authentication.tokenDecode", err.Error())
		return errors.New("Failed to decode authenticated user.")
	}

	if authenticatedUser.Email != user.Email {
		a.logBuilder.Add("authentication.differentEmails", fmt.Sprintf("Provided email %s did not match the user email %s", authenticatedUser.Email, user.Email))
		return errors.New("Failed email")
	}

	if authenticatedUser.ApiKey != a.apiKey || authenticatedUser.ProjectID != a.projectId {
		a.logBuilder.Add("authentication.invalidAuthHeaders", fmt.Sprintf("Provided headers did not match the session headers"))
		return errors.New("Failed auth headers")
	}

	refresh := authenticatedUser.Refresh
	if time.Now().After(refresh.Add(61 * time.Minute)) {
		return errors.New("Users session has expired")
	}

	/**
		- Refresh is the time when the token was created
	    - Now is the time when the request was sent
	    - Now() must be > than refresh time + refresh interval
	*/
	if time.Now().After(refresh.Add(15 * time.Minute)) {
		authenticatedUser.Refresh = time.Now()
		a.doRefresh = true
	}

	a.key = key
	a.user = authenticatedUser

	return nil
}

func (a *apiAuthentication) User() AuthenticatedUser {
	return a.user
}

func (a *apiAuthentication) ShouldRefresh() bool {
	return a.doRefresh
}

func (a *apiAuthentication) Logout(cb func()) {
	cb()
}

func (a *apiAuthentication) Refresh() (string, error) {
	if a.doRefresh {
		loginer := NewApiLogin(a.user, a.key, a.logBuilder)
		return loginer.Login()
	}

	return a.session, nil
}

func NewApiAuthentication(session, projectId, apiKey string, builder logger.LogBuilder) Authentication {
	return &apiAuthentication{
		session:    session,
		doRefresh:  false,
		projectId:  projectId,
		apiKey:     apiKey,
		logBuilder: builder,
	}
}
