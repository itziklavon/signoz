package config

import (
	"goapm/http"
	"goapm/logger"
	"errors"
	"github.com/jarcoal/httpmock"
	jsoniter "github.com/json-iterator/go"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

var classUnderTest ConfServerService

func TestMain(m *testing.M) {
	os.Setenv("confsrvDomain", "http://config.capcap1.com")
	logger.InitLogger()
	classUnderTest = ConfServerService{
		Log:        logger.LOGGER,
		Json:       jsoniter.ConfigCompatibleWithStandardLibrary,
		RestClient: http.NewRestClient(),
	}

	http.RestClientService.Client.SetTimeout(1 * time.Second)
	http.RestClientService.Client.SetRetryCount(0)
	httpmock.ActivateNonDefault(http.RestClientService.Client.GetClient())
	m.Run()
}

func setupTestCase() func(t *testing.T) {
	httpmock.RegisterResponder("GET", "http://config.capcap1.com/tst/tst/master",
		httpmock.NewStringResponder(200, "     {\"name\":\"tst\",\"profiles\":[\"dev\"],\"label\":\"master\",\"version\":\"b8cfe2779e9604804e625135b96b4724ea378736\",\n     \"propertySources\":[\n        {\"name\":\"https://github.com/eriklupander/go-microservice-config.git/accountservice-dev.yml\",\n        \"source\":\n            {\"server_port\":6767,\"server_name\":\"tst\"}\n        }]\n     }"))
	return func(t *testing.T) {
		t.Log("teardown test case")
	}
}

func setupTestCasePanic() func(t *testing.T) {
	httpmock.RegisterResponder("GET", "http://config.capcap1.com/tst/tst/master",
		httpmock.NewErrorResponder(errors.New("error")))
	return func(t *testing.T) {
		t.Log("teardown test case")
	}
}

func TestLoadConfigurationFromBranch(t *testing.T) {
	teardownTestCase := setupTestCase()
	defer teardownTestCase(t)

	classUnderTest.LoadConfigurationFromBranch("http://config.capcap1.com", "tst", "tst", "master")
	assert.Equal(t, "6767", viper.GetString("server_port"))
	assert.Equal(t, "tst", viper.GetString("server_name"))
}

func TestLoadConfigurationFromBranchPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	teardownTestCase := setupTestCasePanic()
	defer teardownTestCase(t)
	classUnderTest.LoadConfigurationFromBranch("http://config.capcap1.com", "tst", "tst", "master")
	assert.Equal(t, "", viper.GetString("server_port"))
	assert.Equal(t, "", viper.GetString("server_name"))
}
