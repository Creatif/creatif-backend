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

func GetApiAuthenticationCookie(c echo.Context) string {
	cookie, err := c.Cookie("api_authentication")
	if err != nil {
		return ""
	}

	return cookie.Value
}

func EncryptAuthenticationCookie(encryptedUser string) *http.Cookie {
	if os.Getenv("APP_ENV") == "prod" {
		cookie := new(http.Cookie)
		cookie.Name = "authentication"
		cookie.HttpOnly = true
		cookie.Secure = true
		cookie.SameSite = http.SameSiteStrictMode
		cookie.Domain = "https://api.creatif.app"
		cookie.Path = "/api/v1"

		cookie.Value = encryptedUser
		cookie.Expires = time.Now().Add(1 * time.Hour)

		return cookie
	}

	cookie := new(http.Cookie)
	cookie.Name = "authentication"
	cookie.Path = "/"

	cookie.Value = encryptedUser
	cookie.Expires = time.Now().Add(1 * time.Hour)

	return cookie
}

func EncryptApiAuthenticationCookie(encryptedUser string) *http.Cookie {
	if os.Getenv("APP_ENV") == "prod" {
		cookie := new(http.Cookie)
		cookie.Name = "api_authentication"
		cookie.HttpOnly = true
		cookie.Secure = true
		cookie.SameSite = http.SameSiteStrictMode
		cookie.Domain = "https://api.creatif.app"
		cookie.Path = "/api/v1"

		cookie.Value = encryptedUser
		cookie.Expires = time.Now().Add(1 * time.Hour)

		return cookie
	}

	cookie := new(http.Cookie)
	cookie.Name = "api_authentication"
	cookie.Path = "/"

	cookie.Value = encryptedUser
	cookie.Expires = time.Now().Add(1 * time.Hour)

	return cookie
}

func RemoveApiAuthenticationCookie() *http.Cookie {
	if os.Getenv("APP_ENV") == "prod" {
		cookie := new(http.Cookie)
		cookie.Name = "api_authentication"
		cookie.HttpOnly = true
		cookie.Secure = true
		cookie.SameSite = http.SameSiteStrictMode
		cookie.Domain = "https://api.creatif.app"
		cookie.Path = "/api/v1"

		cookie.Value = ""
		cookie.Expires = time.Unix(0, 0)

		return cookie
	}

	cookie := new(http.Cookie)
	cookie.Name = "api_authentication"
	cookie.Path = "/"

	cookie.Value = ""
	cookie.Expires = time.Unix(0, 0)

	return cookie
}
