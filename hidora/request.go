package hidora

import (
	"io"
	"net/http"
	"net/url"
)

type HidoraRequest struct {
	Method  string
	Headers http.Header
	Query   url.Values
	Body    io.Reader
}
