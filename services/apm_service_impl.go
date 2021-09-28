package services

import (
	"goapm/logger"
	"context"
	"go.uber.org/zap"
	"goapm/dao"
	model "goapm/domain"
	"sync"
)

var apmServiceOnce sync.Once
var apmService *ApmServiceImpl

type ApmServiceImpl struct {
	Logger *zap.SugaredLogger
	ApmDao dao.ApmDao
}

func NewApmServiceImpl(ApmDao dao.ApmDao) *ApmServiceImpl {
	apmServiceOnce.Do(func() {
		apmService = &ApmServiceImpl{
			Logger: logger.LOGGER,
			ApmDao: ApmDao,
		}
	})
	return apmService
}

func (service *ApmServiceImpl) GetServices(ctx context.Context, query *model.GetServicesParams) (*[]model.ServiceItem, error) {
	return service.ApmDao.GetServices(ctx, query)
}

func (service *ApmServiceImpl) GetServicesList(ctx context.Context) (*[]string, error) {
	return service.ApmDao.GetServicesList(ctx)
}

func (service *ApmServiceImpl) GetServiceOverview(ctx context.Context, query *model.GetServiceOverviewParams) (*[]model.ServiceOverviewItem, error) {
	return service.ApmDao.GetServiceOverview(ctx, query)
}

func (service *ApmServiceImpl) SearchSpans(ctx context.Context, query *model.SpanSearchParams) (*[]model.SearchSpansResult, error) {
	return service.ApmDao.SearchSpans(ctx, query)
}

func (service *ApmServiceImpl) GetServiceDBOverview(ctx context.Context, query *model.GetServiceOverviewParams) (*[]model.ServiceDBOverviewItem, error) {
	return service.ApmDao.GetServiceDBOverview(ctx, query)
}

func (service *ApmServiceImpl) GetServiceExternalAvgDuration(ctx context.Context, query *model.GetServiceOverviewParams) (*[]model.ServiceExternalItem, error) {
	return service.ApmDao.GetServiceExternalAvgDuration(ctx, query)
}

func (service *ApmServiceImpl) GetServiceExternalErrors(ctx context.Context, query *model.GetServiceOverviewParams) (*[]model.ServiceExternalItem, error) {
	return service.ApmDao.GetServiceExternalErrors(ctx, query)
}

func (service *ApmServiceImpl) GetServiceExternal(ctx context.Context, query *model.GetServiceOverviewParams) (*[]model.ServiceExternalItem, error) {
	return service.ApmDao.GetServiceExternal(ctx, query)
}

func (service *ApmServiceImpl) GetTopEndpoints(ctx context.Context, query *model.GetTopEndpointsParams) (*[]model.TopEndpointsItem, error) {
	return service.ApmDao.GetTopEndpoints(ctx, query)
}

func (service *ApmServiceImpl) GetUsage(ctx context.Context, query *model.GetUsageParams) (*[]model.UsageItem, error) {
	return service.ApmDao.GetUsage(ctx, query)
}

func (service *ApmServiceImpl) GetOperations(ctx context.Context, serviceName string) (*[]string, error) {
	return service.ApmDao.GetOperations(ctx, serviceName)
}

func (service *ApmServiceImpl) GetTags(ctx context.Context, serviceName string) (*[]model.TagItem, error) {
	return service.ApmDao.GetTags(ctx, serviceName)
}

func (service *ApmServiceImpl) SearchTraces(ctx context.Context, traceID string) (*[]model.SearchSpansResult, error) {
	return service.ApmDao.SearchTraces(ctx, traceID)
}

func (service *ApmServiceImpl) GetServiceMapDependencies(ctx context.Context, query *model.GetServicesParams) (*[]model.ServiceMapDependencyResponseItem, error) {
	return service.ApmDao.GetServiceMapDependencies(ctx, query)
}

func (service *ApmServiceImpl) SearchSpansAggregate(ctx context.Context, queryParams *model.SpanSearchAggregatesParams) ([]model.SpanSearchAggregatesResponseItem, error) {
	return service.ApmDao.SearchSpansAggregate(ctx, queryParams)
}
