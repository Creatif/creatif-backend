package main

import (
	"creatif-sdk-seed/errorHandler"
	http2 "creatif-sdk-seed/http"
	"encoding/json"
	"fmt"
	"github.com/schollz/progressbar/v3"
	"io"
	"net/http"
	"regexp"
	"sync"
)

func publish(client *http.Client, projectId, name string) http2.HttpResult {
	body := map[string]string{
		"name": name,
	}

	b, err := json.Marshal(body)
	if err != nil {
		return http2.NewHttpResult(nil, err, 0, false, Cannot_Continue_Procedure)
	}

	url := fmt.Sprintf("%s%s%s", URL, "/publishing/publish/", projectId)
	req, err := http2.NewRequest(http2.Request{
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Url:    url,
		Method: "PUT",
		Body:   b,
	})
	if err != nil {
		return http2.NewHttpResult(nil, err, 0, false, Cannot_Continue_Procedure)
	}

	response, err := http2.Make(req, client)
	if err != nil {
		// This is a hack because of a problem explained in the thread below.
		// io.EOF is not the same error since the error is created by the http client,
		// not in the response.
		re := regexp.MustCompile(`\bEOF\b`)
		if re.MatchString(err.Error()) {
			return http2.NewHttpResult(nil, nil, 200, true, Cannot_Continue_Procedure)
			// safe to ignore this per this so thread https://stackoverflow.com/questions/17714494/golang-http-request-results-in-eof-errors-when-making-multiple-requests-successi
		} else {
			return http2.NewHttpResult(nil, err, 0, false, Cannot_Continue_Procedure)
		}
	}

	if response != nil && response.Body != nil {
		defer response.Body.Close()
	}

	return http2.NewHttpResult(response, err, response.StatusCode, response.StatusCode >= 200 && response.StatusCode <= 299, Cannot_Continue_Procedure)
}

func publishProjects(client *http.Client, projectProducts []projectProduct) {
	progressbar.NewOptions(-1, progressbar.OptionSetDescription(fmt.Sprintf("Publishing project version %s")))
	wg := sync.WaitGroup{}
	wg.Add(len(projectProducts))
	/**
	No matter which project it is, it will always have a single published version and that versions name
	will be v1.
	*/
	for i, projectListener := range projectProducts {
		go func(product projectProduct, versionIdx int) {
			defer wg.Done()
			fmt.Println("Publishing project version v1")

			errorHandler.HandleHttpError(publish(client, product.projectId, "v1"))
		}(projectListener, i)
	}

	wg.Wait()
}
