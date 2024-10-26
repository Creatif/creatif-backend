package http

import (
	"bytes"
	"net/http"
	"time"
)

type checkRedirectFn = func(req *http.Request, via []*http.Request) error

type clientParams struct {
	Transport     http.RoundTripper
	CheckRedirect checkRedirectFn
	Jar           http.CookieJar
	Timeout       time.Duration
}

func NewClientParams(transport http.RoundTripper, checkRedirect checkRedirectFn, jar http.CookieJar, timeout time.Duration) clientParams {
	return clientParams{
		Transport:     transport,
		CheckRedirect: checkRedirect,
		Jar:           jar,
		Timeout:       timeout,
	}
}

func NewClient(params clientParams) *http.Client {
	return &http.Client{
		Transport:     params.Transport,
		CheckRedirect: params.CheckRedirect,
		Jar:           params.Jar,
		Timeout:       params.Timeout,
	}
}

func NewRequest(request Request) (*http.Request, error) {
	r, err := http.NewRequest(request.Method, request.Url, bytes.NewBuffer(request.Body))

	if err != nil {
		return nil, err
	}

	if len(request.Headers) != 0 {
		for k, v := range request.Headers {
			r.Header.Set(k, v)
		}
	}

	return r, nil
}

func Make(request *http.Request, client *http.Client) (*http.Response, error) {
	return client.Do(request)
}
