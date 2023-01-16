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

const (
	HIDORA_HOST_URL       string        = "app.hidora.com"
	HIDORA_CLIENT_TIMEOUT time.Duration = 1800 // 30 minutes
	HIDORA_API_VERSION    string        = "1.0"
	HIDORA_DEFAULT_REGION string        = "ch-gen"
)

func defaultparams() *HidoraClient {
	return &HidoraClient{
		httpClient: &http.Client{
			Timeout: time.Second * HIDORA_CLIENT_TIMEOUT,
		},
		auth: "",
		apiUrl: &url.URL{
			Scheme: "https",
			Host:   HIDORA_HOST_URL,
			Path:   "/" + HIDORA_API_VERSION + "/",
		},
		apiVersion:      HIDORA_API_VERSION,
		userAgent:       "Golang_Hidora_SDK/" + HIDORA_API_VERSION,
		defaultEnvGroup: "Hidora_EnvGroup",
		defaultRegion:   HIDORA_DEFAULT_REGION,
		defaultZone:     HIDORA_DEFAULT_REGION + "-1",
	}
}

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

// Implement "Get" functions to be used by future projects (like Terraform provider)
func (s *HidoraClient) GetDefaultEnvGroup() {}

func (s *HidoraClient) GetDefaultRegion() {}

func (s *HidoraClient) GetDefaultZone() {}
