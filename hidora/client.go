package hidora

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
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
	HIDORA_DEFAULT_REGION HidoraRegion  = RegionChGen
	HIDORA_DEFAULT_ZONE   HidoraZone    = ZoneChGen1
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
		defaultZone:     HIDORA_DEFAULT_ZONE,
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
func (s *HidoraClient) GetDefaultEnvGroup() HidoraEnvGroup {
	return s.defaultEnvGroup
}

func (s *HidoraClient) GetDefaultRegion() HidoraRegion {
	return s.defaultRegion
}

func (s *HidoraClient) GetDefaultZone() HidoraZone {
	return s.defaultZone
}

// do performs a single HTTP request based on the ScalewayRequest object.
// /////////////////////////////////////////////////
// Finish implementation of function
// /////////////////////////////////////////////////
func (c *HidoraClient) do(req *HidoraRequest, res interface{}) (sdkErr error) {

	if req == nil {
		// Catch with log library
		return errors.New("Request is nil !")
	}

	// build url
	url, sdkErr := req.getURL(c.apiUrl.Host)
	if sdkErr != nil {
		return sdkErr
	}
	log.Printf("creating %s request on %s", req.Method, url.String())

	// build request
	httpRequest, err := http.NewRequest(req.Method, url.String(), req.Body)
	if err != nil {
		return errors.New("Could not create request")
	}

	httpRequest.Header = req.setHeaders(req.auth, c.userAgent, false)

	// execute request
	httpResponse, err := c.httpClient.Do(httpRequest)
	if err != nil {
		return errors.Wrap(err, "error executing request")
	}

	defer func() {
		closeErr := httpResponse.Body.Close()
		if sdkErr == nil && closeErr != nil {
			sdkErr = errors.Wrap(closeErr, "could not close http response")
		}
	}()

	sdkErr = hasResponseError(httpResponse)
	if sdkErr != nil {
		return sdkErr
	}

	if res != nil {
		contentType := httpResponse.Header.Get("Content-Type")

		switch contentType {
		case "application/json":
			err = json.NewDecoder(httpResponse.Body).Decode(&res)
			if err != nil {
				return errors.Wrap(err, "could not parse %s response body", contentType)
			}
		default:
			buffer, isBuffer := res.(io.Writer)
			if !isBuffer {
				return errors.Wrap(err, "could not handle %s response body with %T result type", contentType, buffer)
			}

			_, err := io.Copy(buffer, httpResponse.Body)
			if err != nil {
				return errors.Wrap(err, "could not copy %s response body", contentType)
			}
		}

		// Handle instance API X-Total-Count header
		xTotalCountStr := httpResponse.Header.Get("X-Total-Count")
		if legacyLister, isLegacyLister := res.(legacyLister); isLegacyLister && xTotalCountStr != "" {
			xTotalCount, err := strconv.Atoi(xTotalCountStr)
			if err != nil {
				return errors.Wrap(err, "could not parse X-Total-Count header")
			}
			legacyLister.UnsafeSetTotalCount(xTotalCount)
		}
	}

	return nil
}
