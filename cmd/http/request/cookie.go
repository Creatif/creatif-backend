package request

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"os"
	"time"
)

func GetAuthenticationCookie(c echo.Context) string {
	cookie, err := c.Cookie("authentication")
	if err != nil {
		return ""
	}

	return cookie.Value
}

func EncryptAuthenticationCookie(encryptedUser string) *http.Cookie {
	cookie := new(http.Cookie)
	cookie.Name = "authentication"
	cookie.HttpOnly = true
	cookie.Secure = true
	cookie.Domain = "https://api.creatif.app"
	cookie.Path = "/api/v1"
	if os.Getenv("APP_ENV") != "prod" {
		cookie.HttpOnly = true
		cookie.Secure = true
		cookie.Domain = "http://localhost"
		cookie.Path = "/"
	}

	cookie.Value = encryptedUser
	cookie.Expires = time.Now().Add(1 * time.Hour)

	return cookie
}
