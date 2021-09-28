package utils

import (
	"context"
	"github.com/stretchr/testify/mock"
)

type HealthCheck struct {
	Status int             `json:"status"`
	Checks []ServiceHealth `json:"health_checks"`
}

type ServiceHealth struct {
	Name          string                 `json:"name"`
	Status        string                 `json:"status"`
	StatusCode    int                    `json:"status_code"`
	ErrorMessages map[string]interface{} `json:"error_messages,omitempty"`
}

type HealthCheckService interface {
	CheckService(ctx context.Context) ServiceHealth
}

type MockHealthCheckService struct {
	mock.Mock
}

func (health *MockHealthCheckService) CheckService(ctx context.Context) ServiceHealth {
	args := health.Called(mock.Anything)
	return args.Get(0).(ServiceHealth)
}

func ConstructHealthCheckResponse(serviceDetails ...ServiceHealth) *HealthCheck {
	var healthCheckResponses []ServiceHealth
	healthCheckResponse := &HealthCheck{
		Status: 200,
		Checks: healthCheckResponses,
	}
	for _, check := range serviceDetails {
		if check.StatusCode != 200 {
			healthCheckResponse.Status = check.StatusCode
		}

		healthCheckResponses = append(healthCheckResponses, check)
	}
	healthCheckResponse.Checks = healthCheckResponses
	return healthCheckResponse
}
