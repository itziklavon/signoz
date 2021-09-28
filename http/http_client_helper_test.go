package http

import (
	"context"
	"errors"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestRestClient_GetResponse(t *testing.T) {
	restClient := NewRestClient()
	UpdateRestClient()

	RestClientService.Client.SetTimeout(1 * time.Second)
	RestClientService.Client.SetRetryCount(0)
	httpmock.ActivateNonDefault(RestClientService.Client.GetClient())

	httpmock.RegisterResponder("GET", "http://config.capcap1.com/tst/tst/master",
		httpmock.NewStringResponder(200, "{}"))
	response, err := restClient.GetResponse(context.Background(), "http://config.capcap1.com/tst/tst/master", nil)
	assert.Nil(t, err)
	assert.Equal(t, 200, response.StatusCode)

	httpmock.RegisterResponder("GET", "http://config.capcap1.com/tst/tst/master",
		httpmock.NewErrorResponder(errors.New("Error")))
	response, err = restClient.GetResponse(context.Background(), "http://config.capcap1.com/tst/tst/master", nil)
	assert.NotNil(t, err)
	assert.Equal(t, 500, response.StatusCode)
}

func TestRestClient_PostResponse(t *testing.T) {
	restClient := NewRestClient()
	UpdateRestClient()

	RestClientService.Client.SetTimeout(1 * time.Second)
	RestClientService.Client.SetRetryCount(0)
	httpmock.ActivateNonDefault(RestClientService.Client.GetClient())

	httpmock.RegisterResponder("POST", "http://config.capcap1.com/tst/tst/master",
		httpmock.NewStringResponder(200, "{}"))
	response, err := restClient.PostResponse(context.Background(), "http://config.capcap1.com/tst/tst/master", nil, nil)
	assert.Nil(t, err)
	assert.Equal(t, 200, response.StatusCode)

	httpmock.RegisterResponder("POST", "http://config.capcap1.com/tst/tst/master",
		httpmock.NewErrorResponder(errors.New("Error")))
	response, err = restClient.PostResponse(context.Background(), "http://config.capcap1.com/tst/tst/master", nil, nil)
	assert.NotNil(t, err)
	assert.Equal(t, 500, response.StatusCode)
}

func TestRestClient_PutResponse(t *testing.T) {
	restClient := NewRestClient()
	UpdateRestClient()

	RestClientService.Client.SetTimeout(1 * time.Second)
	RestClientService.Client.SetRetryCount(0)
	httpmock.ActivateNonDefault(RestClientService.Client.GetClient())

	httpmock.RegisterResponder("PUT", "http://config.capcap1.com/tst/tst/master",
		httpmock.NewStringResponder(200, "{}"))
	response, err := restClient.PutResponse(context.Background(), "http://config.capcap1.com/tst/tst/master", nil, nil)
	assert.Nil(t, err)
	assert.Equal(t, 200, response.StatusCode)

	httpmock.RegisterResponder("PUT", "http://config.capcap1.com/tst/tst/master",
		httpmock.NewErrorResponder(errors.New("Error")))
	response, err = restClient.PutResponse(context.Background(), "http://config.capcap1.com/tst/tst/master", nil, nil)
	assert.NotNil(t, err)
	assert.Equal(t, 500, response.StatusCode)
}
