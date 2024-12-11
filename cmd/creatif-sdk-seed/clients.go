package main

import (
	"creatif-sdk-seed/errorHandler"
	http2 "creatif-sdk-seed/http"
	"crypto/tls"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"
)

func createAnonymousClient() *http.Client {
	return http2.NewClient(http2.NewClientParams(&http.Transport{
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
		MaxConnsPerHost:     1024,
		TLSHandshakeTimeout: 0,
	}, nil, nil, 20*time.Minute))
}

func createAuthenticatedClient(authToken string) *http.Client {
	cookieJar, err := cookiejar.New(nil)
	if err != nil {
		errorHandler.HandleHttpError(http2.NewHttpResult(nil, err, 0, false, Cannot_Continue_Procedure))
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

	return http2.NewClient(http2.NewClientParams(&http.Transport{
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
		MaxConnsPerHost:     1024,
		TLSHandshakeTimeout: 0,
	}, nil, cookieJar, 20*time.Minute))
}
