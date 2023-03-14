package hidora

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type HidoraRequest struct {
	Method  string
	Path    string
	Headers http.Header
	Query   url.Values
	Body    io.Reader
}

var DEFAULT_HEADERS map[string]string = map[string]string{
	"Content-Type": "application/x-www-form-urlencoded",
	"Accept":       "application/json",
}

// getAllHeaders constructs a http.Header object and aggregates all headers into the object.
// /////////////////////////////////////////////////
// Finish implementation of function
// /////////////////////////////////////////////////
func (req *HidoraRequest) setHeaders(headers map[string]string, userAgent string) http.Header {
	var allHeaders http.Header

	allHeaders.Set("User-Agent", userAgent)
	if req.Body != nil {
		allHeaders.Set("Content-Type", "application/json")
	}
	for key, value := range req.Headers {
		allHeaders.Del(key)
		for _, v := range value {
			allHeaders.Add(key, v)
		}
	}
	return allHeaders
}

// getURL constructs a URL based on the base url and the client.
func (req *HidoraRequest) getURL(baseURL string) (*url.URL, error) {
	url, err := url.Parse(baseURL + req.Path)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Invalid base URLurl %s: %s", baseURL+req.Path, err))
	}
	url.RawQuery = req.Query.Encode()

	return url, nil
}

// SetBody json marshal the given body and write the json content type
// to the request. It also catches when body is a file.
func (req *HidoraRequest) SetBody(body interface{}) error {
	var contentType string
	var content io.Reader

	switch b := body.(type) {
	case io.Reader:
		contentType = "text/plain"
		content = b
	default:
		buf, err := json.Marshal(body)
		if err != nil {
			return err
		}
		contentType = "application/json"
		content = bytes.NewReader(buf)
	}

	if req.Headers == nil {
		req.Headers = http.Header{}
	}

	req.Headers.Set("Content-Type", contentType)
	req.Body = content

	return nil
}

func (req *HidoraRequest) validate() error {
	// Use with validate go package
	// Implement when validate package ready
	return nil
}
