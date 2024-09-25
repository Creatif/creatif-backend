package auth

import (
	"creatif/pkg/app/domain/app"
	"creatif/pkg/lib/storage"
	"encoding/base64"
	"encoding/json"
	"errors"
	"time"
)

type apiAuthentication struct {
	session   string
	user      AuthenticatedUser
	key       [32]byte
	doRefresh bool
}

func (a *apiAuthentication) Authenticate() error {
	if a.session == "" {
		return errors.New("Authentication cookie not provided.")
	}

	jsonToken, err := base64.StdEncoding.DecodeString(a.session)
	if err != nil {
		return errors.New("Failed to decode session cookie")
	}

	var session AuthenticatedApiSession
	if err := json.Unmarshal(jsonToken, &session); err != nil {
		return errors.New("Failed to decode API session object")
	}

	if session.Type != "api" {
		return errors.New("Unauthenticated")
	}

	var user app.User
	if res := storage.Gorm().Where("id = ?", session.ID).Select("key", "email").First(&user); res.Error != nil {
		return errors.New("User not found")
	}

	encrypedUser, err := base64.StdEncoding.DecodeString(session.Token)
	if err != nil {
		return errors.New("Failed to decode token.")
	}

	var key [32]byte
	for i, v := range user.Key {
		key[i] = byte(v)
	}

	jsonAuthenticatedUser, err := decrypt(encrypedUser, &key)
	if err != nil {
		return errors.New("Failed to decrypt decoded user.")
	}

	var authenticatedUser AuthenticatedUser
	if err := json.Unmarshal(jsonAuthenticatedUser, &authenticatedUser); err != nil {
		return errors.New("Failed to decode authenticated user.")
	}

	if authenticatedUser.Email != user.Email {
		return errors.New("Failed email")
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
		loginer := NewApiLogin(a.user, a.key)
		return loginer.Login()
	}

	return a.session, nil
}

func NewApiAuthentication(session string) Authentication {
	return &apiAuthentication{
		session:   session,
		doRefresh: false,
	}
}
