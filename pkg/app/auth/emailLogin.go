package auth

import (
	"creatif/pkg/lib/appErrors"
	"creatif/pkg/lib/logger"
	"encoding/base64"
	"encoding/json"
)

type EmailLogin struct {
	key        [32]byte
	user       AuthenticatedUser
	logBuilder logger.LogBuilder
}

func (e EmailLogin) Login() (string, error) {
	serializedUser, err := json.Marshal(&e.user)
	if err != nil {
		return "", appErrors.NewApplicationError(err)
	}

	encryptedUser, err := encrypt(serializedUser, &e.key)
	if err != nil {
		// TODO: send immediate slack message here
		e.logBuilder.Add("login.cannotEncryptKey", "Key should have length of 32 characters.")
		return "", appErrors.NewUnexpectedError(err)
	}

	token := base64.StdEncoding.EncodeToString(encryptedUser)
	session := NewAuthenticatedFrontendSession(e.user.ID, token)

	b, err := json.Marshal(session)
	if err != nil {
		// TODO: send immediate slack message here
		e.logBuilder.Add("login.cannotEncryptKey", "Key should have length of 32 characters.")
		return "", appErrors.NewUnexpectedError(err)
	}

	return base64.StdEncoding.EncodeToString(b), nil
}

func NewEmailLogin(user AuthenticatedUser, key [32]byte, logger logger.LogBuilder) Loginer {
	return EmailLogin{
		user:       user,
		key:        key,
		logBuilder: logger,
	}
}
