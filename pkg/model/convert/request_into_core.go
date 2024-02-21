package convert

import (
	"io/ioutil"
	"main/pkg/model/core"
	"net/http"
	"net/url"
	"strings"
)

const (
	contentTypeParsePost = "application/x-www-form-urlencoded"
)

func ParseHTTPRequest(r *http.Request) *core.Request {
	parsedRequest := &core.Request{
		ID:        core.GenUID(),
		Method:    r.Method,
		Path:      r.RequestURI,
		Headers:   make(map[string]string),
		Cookies:   make(map[string]string),
		GetParams: make(map[string]string),
	}

	// Add GET parameters
	for key, values := range r.URL.Query() {
		parsedRequest.GetParams[key] = strings.Join(values, ", ")
	}

	// Add headers
	for key, values := range r.Header {
		parsedRequest.Headers[key] = strings.Join(values, ", ")
	}

	// Add cookies
	for _, cookie := range r.Cookies() {
		parsedRequest.Cookies[cookie.Name] = cookie.Value
	}

	// Add POST parameters
	if r.Method == http.MethodPost &&
		strings.Contains(r.Header.Get("Content-Type"), contentTypeParsePost) {
		body, _ := ioutil.ReadAll(r.Body)
		defer r.Body.Close()

		formData, _ := url.ParseQuery(string(body))

		parsedRequest.PostParams = make(map[string]string)
		for key, values := range formData {
			parsedRequest.PostParams[key] = strings.Join(values, ", ")
		}
	}

	return parsedRequest
}
