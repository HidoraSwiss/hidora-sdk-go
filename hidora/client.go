package hidora

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
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
	hidoraHostUrl       string        = "app.hidora.com"
	hidoraClientTimeout time.Duration = 1800 // 30 minutes
	hidoraApiVersion    string        = "1.0"
	hidoraDefaultRegion HidoraRegion  = RegionChGen
	hidoraDefaultZone   HidoraZone    = ZoneChGen1
)

func defaultparams() *HidoraClient {
	return &HidoraClient{
		httpClient: &http.Client{
			Timeout: time.Second * hidoraClientTimeout,
		},
		auth: "",
		apiUrl: &url.URL{
			Scheme: "https",
			Host:   hidoraHostUrl,
			Path:   "/" + hidoraApiVersion + "/",
		},
		apiVersion:      hidoraApiVersion,
		userAgent:       "Golang_Hidora_SDK/" + hidoraApiVersion,
		defaultEnvGroup: "Hidora_EnvGroup",
		defaultRegion:   hidoraDefaultRegion,
		defaultZone:     hidoraDefaultZone,
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
		return fmt.Errorf("Could not create request, error : %w", err)
	}

	httpRequest.Header = req.setHeaders(baseHeaders, c.userAgent)

	// execute request
	httpResponse, err := c.httpClient.Do(httpRequest)
	if err != nil {
		return fmt.Errorf("error executing request, error : %w", err)
	}

	defer func() {
		closeErr := httpResponse.Body.Close()
		if sdkErr == nil && closeErr != nil {
			sdkErr = fmt.Errorf(
				"could not close http response, error : %w",
				closeErr)
		}
	}()

	if !strings.HasPrefix(httpResponse.Status, "200") {
		return fmt.Errorf("Error on response")
	}

	if res != nil {
		contentType := httpResponse.Header.Get("Content-Type")

		switch contentType {
		case "application/json":
			err = json.NewDecoder(httpResponse.Body).Decode(&res)
			if err != nil {
				return fmt.Errorf(
					"could not parse %s response body, error : %w",
					contentType, err)
			}
		default:
			buffer, isBuffer := res.(io.Writer)
			if !isBuffer {
				return fmt.Errorf(
					"could not handle %s response body with %T result type, error : %w",
					contentType, buffer, err)
			}

			_, err := io.Copy(buffer, httpResponse.Body)
			if err != nil {
				return fmt.Errorf(
					"could not copy %s response body, error : %w",
					contentType, err)
			}
		}
	}
	return nil
}
