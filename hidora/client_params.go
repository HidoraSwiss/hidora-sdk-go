package hidora

import (
	"net/http"
	"net/url"
)

// Use to handle parameters send to go client
type settings struct {
	httpClient      *http.Client
	auth            HidoraAuth
	apiUrl          *url.URL
	apiVersion      string
	userAgent       string
	defaultEnvGroup HidoraEnvGroup
	defaultRegion   HidoraRegion
	defaultZone     HidoraZone
}

// Implementation of a new set of settings
func emptySettings() *settings {
	return &settings{}
}

// Implementation to apply parameters set by user
func (s *settings) apply() {}

// Implementation to validate each parameter
func (s *settings) validate() {}
