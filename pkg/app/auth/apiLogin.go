package auth

import (
	"creatif/pkg/lib/appErrors"
	"encoding/base64"
	"encoding/json"
)

type ApiLogin struct {
	key  [32]byte
	user AuthenticatedUser
}

func (e ApiLogin) Login() (string, error) {
	serializedUser, err := json.Marshal(&e.user)
	if err != nil {
		return "", appErrors.NewApplicationError(err)
	}

	encryptedUser, err := encrypt(serializedUser, &e.key)
	if err != nil {
		// TODO: send immediate slack message here
		return "", appErrors.NewUnexpectedError(err)
	}

	token := base64.StdEncoding.EncodeToString(encryptedUser)
	session := NewAuthenticatedApiSession(e.user.ID, token)

	b, err := json.Marshal(session)
	if err != nil {
		// TODO: send immediate slack message here
		return "", appErrors.NewUnexpectedError(err)
	}

	return base64.StdEncoding.EncodeToString(b), nil
}

func NewApiLogin(user AuthenticatedUser, key [32]byte) Loginer {
	return ApiLogin{
		user: user,
		key:  key,
	}
}
