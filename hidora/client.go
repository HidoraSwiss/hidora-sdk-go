package hidora

import (
	"net/http"
	"net/url"
)

type Client struct {
	BaseUrl    *url.URL
	HTTPClient *http.Client
	Token      string
}
