package convert

import (
	"io/ioutil"
	"main/pkg/model/core"
	"net/http"
	"strings"
)

func ParseHTTPResponse(resp *http.Response) *core.Response {
	parsedResponse := &core.Response{
		ID:      core.GenUID(),
		Code:    resp.StatusCode,
		Message: resp.Status,
		Headers: make(map[string]string),
	}

	// Add headers
	for key, values := range resp.Header {
		parsedResponse.Headers[key] = strings.Join(values, ", ")
	}

	// Parse body
	body, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	parsedResponse.Body = string(body)

	return parsedResponse
}
