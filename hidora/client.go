package hidora

import (
	"net/http"
	"net/url"
	"time"
)

type HidoraAuth string
type HidoraEnvGroup string

// type HidoraRegion string
// type HidoraZone string

type HidoraClient struct {
	httpClient      *http.Client
	auth            HidoraAuth
	apiUrl          *url.URL
	apiVersion      string
	userAgent       string
	defaultEnvGroup HidoraEnvGroup
	defaultRegion   HidoraRegion
	defaultZone     HidoraZone
}

const HIDORA_HOST_URL string = "app.hidora.com"
const HIDORA_CLIENT_TIMEOUT time.Duration = 1800 // 30 minutes

func NewClient(auth HidoraAuth,
	host string,
	apiVersion string,
	defaultEnvGroup HidoraEnvGroup,
	defaultRegion HidoraRegion,
	defaultZone HidoraZone) (*HidoraClient, error) {
	var hidoraClient = defaultparams()
	if host != "" {
		hidoraClient.apiUrl.Host = host
	}
	if apiVersion != "" {
		hidoraClient.apiVersion = apiVersion
		hidoraClient.apiUrl.Path = "/" + apiVersion + "/"
	}
	if defaultEnvGroup != "" {
		hidoraClient.defaultEnvGroup = defaultEnvGroup
	}
	if defaultRegion != "" {
		HidoraClient.defaultRegion = defaultRegion
	}
	if defaultZone != "" {
		HidoraClient.defaultZone = defaultZone
	}
	return hidoraClient, nil
}

func defaultparams() *HidoraClient {
	return &HidoraClient{
		httpClient: &http.Client{
			Timeout: time.Second * HIDORA_CLIENT_TIMEOUT,
		},
		auth: "",
		apiUrl: &url.URL{
			Scheme: "https",
			Host:   HIDORA_HOST_URL,
			Path:   "/1.0/",
		},
		apiVersion:      "1.0",
		userAgent:       "Golang_Hidora_SDK/1.0",
		defaultEnvGroup: "Hidora_EnvGroup",
		defaultRegion:   "ch-gen",
		defaultZone:     "ch-gen-1",
	}
}
