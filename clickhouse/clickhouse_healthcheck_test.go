package clickhouse

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func NewClickhouseHealthCheckServiceTest(service ClickhouseConnectionService) *ClickhouseHealthCheckImpl {
	return &ClickhouseHealthCheckImpl{
		ClickhouseConnectionService: service,
	}
}

func TestHealthyClickhouseSuccess(t *testing.T) {
	lastRunTime = time.Now().UTC().Add(-1 * time.Hour)
	mockDB, _, _ := sqlmock.New()
	defer mockDB.Close()
	sqlxDBMock := sqlx.NewDb(mockDB, "sqlmock")
	connectionsMap.Put("1", sqlxDBMock)
	clickhouseConnectionMock := new(MockClickhouseConnectionService)
	clickhouseConnectionMock.On("GetConnectionMap").Return(connectionsMap)
	classUnderTest := NewClickhouseHealthCheckServiceTest(clickhouseConnectionMock)
	serviceDetails := classUnderTest.CheckService(context.Background())
	assert.True(t, serviceDetails.StatusCode == 200)
}

func TestHealthyClickhouseFailed(t *testing.T) {
	lastRunTime = time.Now().UTC().Add(-1 * time.Hour)
	mockDB, _, _ := sqlmock.New()
	defer mockDB.Close()
	sqlxDBMock := sqlx.NewDb(mockDB, "sqlmock")
	connectionsMap.Put("1", sqlxDBMock)
	dataSources.Add("1")
	sqlxDBMock.Close()
	clickhouseConnectionMock := new(MockClickhouseConnectionService)
	clickhouseConnectionMock.On("GetConnectionMap").Return(connectionsMap)
	classUnderTest := NewClickhouseHealthCheckServiceTest(clickhouseConnectionMock)
	serviceDetails := classUnderTest.CheckService(context.Background())
	assert.True(t, serviceDetails.StatusCode == 503)
}

func TestHealthyClickhouseSingleton(t *testing.T) {
	clickhouseConnectionMock := new(MockClickhouseConnectionService)
	_ = NewClickhouseHealthCheckService(clickhouseConnectionMock)
	_ = NewClickhouseHealthCheckService(clickhouseConnectionMock)
}
