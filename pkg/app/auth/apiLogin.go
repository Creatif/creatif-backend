package auth

import (
	"creatif/pkg/lib/appErrors"
	"creatif/pkg/lib/logger"
	"encoding/base64"
	"encoding/json"
)

type ApiLogin struct {
	key        [32]byte
	user       AuthenticatedUser
	logBuilder logger.LogBuilder
}

func (e ApiLogin) Login() (string, error) {
	serializedUser, err := json.Marshal(&e.user)
	if err != nil {
		e.logBuilder.Add("apiLogin", err.Error())
		return "", appErrors.NewApplicationError(err)
	}

	encryptedUser, err := encrypt(serializedUser, &e.key)
	if err != nil {
		// TODO: send immediate slack message here
		e.logBuilder.Add("apiLogin.cannotEncryptKey", err.Error())
		return "", appErrors.NewUnexpectedError(err)
	}

	token := base64.StdEncoding.EncodeToString(encryptedUser)
	session := NewAuthenticatedApiSession(e.user.ID, token)

	b, err := json.Marshal(session)
	if err != nil {
		// TODO: send immediate slack message here
		e.logBuilder.Add("apiLogin.cannotEncryptKey", "Key should have length of 32 characters.")
		return "", appErrors.NewUnexpectedError(err)
	}

	return base64.StdEncoding.EncodeToString(b), nil
}

func NewApiLogin(user AuthenticatedUser, key [32]byte, logger logger.LogBuilder) Loginer {
	return ApiLogin{
		user:       user,
		key:        key,
		logBuilder: logger,
	}
}
