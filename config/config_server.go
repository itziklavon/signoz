package config

import (
	"goapm/http"
	"goapm/logger"
	"context"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"sync"
)

var once sync.Once
var confSrvService *ConfServerService

type ConfServerService struct {
	Log        *zap.SugaredLogger
	Json       jsoniter.API
	RestClient http.RestClientInterface
}

// NewConfSrvService
//creates singleton instance of config server service
//expects http client as constructor parameter
func NewConfSrvService(restClient http.RestClientInterface) *ConfServerService {
	once.Do(func() {
		confSrvService = &ConfServerService{
			Log:        logger.LOGGER,
			Json:       jsoniter.ConfigCompatibleWithStandardLibrary,
			RestClient: restClient,
		}
	})
	return confSrvService
}

// LoadConfigurationFromBranch
// Get response from config server(env variable - confsrvDomain)
// store all in external library - viper
func (confSrv *ConfServerService) LoadConfigurationFromBranch(configServerUrl string, appName string, profile string, branch string) {
	url := fmt.Sprintf("%s/%s/%s/%s", configServerUrl, appName, profile, branch)
	confSrv.Log.Info(fmt.Sprintf("Loading config from %s", url))
	body, err := fetchConfiguration(confSrv, url)
	if err != nil {
		panic("Couldn't load configuration, cannot start. Terminating. Error: " + err.Error())
	}
	parseConfiguration(confSrv, body)
	http.UpdateRestClient()
}

// fetchConfiguration  Make HTTP request to fetch configuration from config server
func fetchConfiguration(confSrv *ConfServerService, url string) ([]byte, error) {
	resp, err := confSrv.RestClient.GetResponse(context.Background(), url, nil)
	if err != nil {
		panic("Couldn't load configuration, cannot start. Terminating. Error: " + err.Error())
	}
	return resp.Body, err
}

// parseConfiguration Pass JSON bytes into struct and then into Viper
func parseConfiguration(confSrv *ConfServerService, body []byte) {
	var cloudConfig springCloudConfig
	err := confSrv.Json.Unmarshal(body, &cloudConfig)
	if err != nil {
		panic("Cannot parse configuration, message: " + err.Error())
	}

	for key, value := range cloudConfig.PropertySources[0].Source {
		viper.Set(key, value)
	}
	if viper.IsSet("server_name") {
		confSrv.Log.Info(fmt.Sprintf("Successfully loaded configuration for service %s", viper.GetString("server_name")))
	}
}

// springCloudConfig Structs having same structure as response from Spring Cloud Config
type springCloudConfig struct {
	Name            string           `json:"name"`
	Profiles        []string         `json:"profiles"`
	Label           string           `json:"label"`
	Version         string           `json:"version"`
	PropertySources []propertySource `json:"propertySources"`
}

type propertySource struct {
	Name   string                 `json:"name"`
	Source map[string]interface{} `json:"source"`
}
