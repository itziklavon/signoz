package http

import (
	"context"
	"github.com/stretchr/testify/mock"
	"net/http"
)

type GenericHttpResponse struct {
	StatusCode int
	Headers    http.Header
	Body       []byte
}

type RestClientInterface interface {
	GetResponse(ctx context.Context, url string, headers map[string]string) (GenericHttpResponse, error)
	PostResponse(ctx context.Context, url string, body interface{}, headers map[string]string) (GenericHttpResponse, error)
	PutResponse(ctx context.Context, url string, body interface{}, headers map[string]string) (GenericHttpResponse, error)
}

type MockRestClient struct {
	mock.Mock
}

func (restClient *MockRestClient) GetResponse(ctx context.Context, url string, headers map[string]string) (GenericHttpResponse, error) {
	args := restClient.Called(mock.Anything, mock.Anything, mock.Anything)
	return args.Get(0).(GenericHttpResponse), args.Error(1)
}

func (restClient *MockRestClient) PostResponse(ctx context.Context, url string, body interface{}, headers map[string]string) (GenericHttpResponse, error) {
	args := restClient.Called(mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	return args.Get(0).(GenericHttpResponse), args.Error(1)
}

func (restClient *MockRestClient) PutResponse(ctx context.Context, url string, body interface{}, headers map[string]string) (GenericHttpResponse, error) {
	args := restClient.Called(mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	return args.Get(0).(GenericHttpResponse), args.Error(1)
}
