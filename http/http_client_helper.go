package http

import (
	"goapm/logger"
	"goapm/utils"
	"context"
	"github.com/go-resty/resty/v2"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"net/http"
	"runtime/debug"
	"strings"
	"sync"
	"time"
)

var once sync.Once
var mutex sync.Mutex

var RestClientService *RestClient

// InitResty Initialize http client, set timeout and retry
func InitResty() *resty.Client {
	httpClient := resty.New()
	httpClient.SetTimeout(10 * time.Second)
	httpClient.SetRetryCount(3)
	httpClient.SetTransport(new(RedirectChecker))
	return httpClient
}

// UpdateRestClient update rest client after getting response from config server, to get required timeout and retry count
func UpdateRestClient() {
	mutex.Lock()
	defer mutex.Unlock()
	RestClientService.Client.SetRetryCount(utils.GetOrDefaultInt(viper.GetString("MAX_RETRIES"), 3))
	RestClientService.Client.SetTimeout(time.Duration(utils.GetOrDefaultInt(viper.GetInt("HTTP_READ_TIMEOUT_IN_MILLIS"), 10)) * time.Second)
}

type RestClient struct {
	Log    *zap.SugaredLogger
	Client *resty.Client
}

// NewRestClient create new singleton instance of rest client
func NewRestClient() *RestClient {
	once.Do(func() {
		RestClientService = &RestClient{
			Log:    logger.LOGGER,
			Client: InitResty(),
		}
	})
	return RestClientService
}

// GetResponse http response, expects full uri and headers(can be nil),
//returns an object - GenericHttpResponse, and error if there is one
func (restClient *RestClient) GetResponse(ctx context.Context, url string, headers map[string]string) (GenericHttpResponse, error) {
	if headers == nil {
		headers = make(map[string]string)
	}
	response, err := RestClientService.Client.R().SetContext(ctx).SetHeaders(headers).Get(url)
	if err != nil {
		restClient.Log.Error("unable to get response from uri - "+url, err, string(debug.Stack()))
		return GenericHttpResponse{
			StatusCode: 500,
			Headers:    nil,
			Body:       nil,
		}, err
	}
	return GenericHttpResponse{
		StatusCode: response.StatusCode(),
		Headers:    response.Header(),
		Body:       response.Body(),
	}, nil
}

// PostResponse http response, expects full uri, request body and headers(can be nil),
//returns an object - GenericHttpResponse, and error if there is one
func (restClient *RestClient) PostResponse(ctx context.Context, url string, body interface{}, headers map[string]string) (GenericHttpResponse, error) {
	response, err := restClient.Client.R().SetContext(ctx).SetHeaders(headers).SetBody(body).Post(url)
	if err != nil {
		restClient.Log.Error("unable to get response from uri - "+url, err, string(debug.Stack()))
		return GenericHttpResponse{
			StatusCode: 500,
			Headers:    nil,
			Body:       nil,
		}, err
	}
	return GenericHttpResponse{
		StatusCode: response.StatusCode(),
		Headers:    response.Header(),
		Body:       response.Body(),
	}, nil
}

// PutResponse  http response, expects full uri, request body and headers(can be nil),
//returns an object - GenericHttpResponse, and error if there is one
func (restClient *RestClient) PutResponse(ctx context.Context, url string, body interface{}, headers map[string]string) (GenericHttpResponse, error) {
	response, err := restClient.Client.R().SetContext(ctx).SetHeaders(headers).SetBody(body).Put(url)
	if err != nil {
		restClient.Log.Error("unable to get response from uri - "+url, err, string(debug.Stack()))
		return GenericHttpResponse{
			StatusCode: 500,
			Headers:    nil,
			Body:       nil,
		}, err
	}
	return GenericHttpResponse{
		StatusCode: response.StatusCode(),
		Headers:    response.Header(),
		Body:       response.Body(),
	}, nil
}

type RedirectChecker struct{}

// RoundTrip checks if response is in status 30X does not contain location,
//if it is not contain  location header, add it
func (RedirectChecker) RoundTrip(req *http.Request) (*http.Response, error) {
	// e.g. patch the request before send it

	resp, err := http.DefaultTransport.RoundTrip(req)

	if err == nil && resp != nil {
		switch resp.StatusCode {
		case
			http.StatusMovedPermanently,
			http.StatusSeeOther,
			http.StatusFound,
			http.StatusTemporaryRedirect,
			http.StatusPermanentRedirect:
			if len(resp.Header.Get("Location")) == 0 {
				if strings.Contains(req.URL.Host, "localhost") || strings.Contains(req.URL.Host, "127.0.0.1") {
					resp.StatusCode = 406
				} else {
					resp.Header.Add("Location", req.URL.Host)
				}
			}
		}
	}
	return resp, err
}
