package clickhouse

import (
	"goapm/ds_utils"
	"github.com/stretchr/testify/mock"
)

type ClickhouseConnectionService interface {
	GetConnectionMap() *ds_utils.ConcurrentHashMap
	ExecuteInsertFunction(query string, statement stmt) error
	ExecuteSelectFunction(isArray bool, dest interface{}, query string, args []interface{}) error
}

type MockClickhouseConnectionService struct {
	mock.Mock
}

func (service *MockClickhouseConnectionService) GetConnectionMap() *ds_utils.ConcurrentHashMap {
	args := service.Called(mock.Anything)
	return args.Get(0).(*ds_utils.ConcurrentHashMap)
}

func (service *MockClickhouseConnectionService) ExecuteInsertFunction(query string, statement stmt) error {
	args := service.Called(mock.Anything, mock.Anything)
	return args.Error(0)
}

func (service *MockClickhouseConnectionService) ExecuteSelectFunction(isArray bool, dest interface{}, query string, args []interface{}) (error) {
	argsMock := service.Called(mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	return argsMock.Error(0)
}
