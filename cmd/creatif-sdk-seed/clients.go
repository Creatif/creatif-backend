package main

import (
	"crypto/tls"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"
)

func createAnonymousClient() *http.Client {
	return newClient(newClientParams(&http.Transport{
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
		MaxConnsPerHost:     1024,
		TLSHandshakeTimeout: 0,
	}, nil, nil, 0))
}

func createAuthenticatedClient(authToken string) *http.Client {
	cookieJar, err := cookiejar.New(nil)
	if err != nil {
		handleHttpError(newHttpResult(nil, err, 0, false, Cannot_Continue_Procedure), nil)
	}

	var cookies []*http.Cookie
	cookie := &http.Cookie{
		Name:    "api_authentication",
		Value:   authToken,
		Path:    "/",
		Expires: time.Now().Add(1 * time.Hour),
	}

	cookies = append(cookies, cookie)

	u, _ := url.Parse("http://localhost:3002")
	cookieJar.SetCookies(u, cookies)

	return newClient(newClientParams(&http.Transport{
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
		MaxConnsPerHost:     1024,
		TLSHandshakeTimeout: 0,
	}, nil, cookieJar, 0))
}
