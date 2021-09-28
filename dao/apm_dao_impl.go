package dao

import (
	"goapm/clickhouse"
	"goapm/logger"
	"context"
	"fmt"
	"go.uber.org/zap"
	model "goapm/domain"
	"strconv"
	"sync"
	"time"
)

var apmDaoOnce sync.Once
var apmDao *ApmDaoImpl

type ApmDaoImpl struct {
	Logger                      *zap.SugaredLogger
	ClickhouseConnectionService clickhouse.ClickhouseConnectionService
}

func NewApmDao(ClickhouseConnectionService clickhouse.ClickhouseConnectionService) *ApmDaoImpl {
	apmDaoOnce.Do(func() {
		apmDao = &ApmDaoImpl{
			Logger:                      logger.LOGGER,
			ClickhouseConnectionService: ClickhouseConnectionService,
		}
	})
	return apmDao
}

func (dao *ApmDaoImpl) GetServices(ctx context.Context, queryParams *model.GetServicesParams) (*[]model.ServiceItem, error) {
	var serviceItems []model.ServiceItem
	query := fmt.Sprintf("SELECT serviceName, quantileMerge(0.99)(quantile) as p99, avgMerge(avg) as avgDuration, sum(count) as numCalls FROM signoz_index_aggregated WHERE timestamp>='%s' AND timestamp<='%s' AND kind='2' GROUP BY serviceName ORDER BY p99 DESC", convertNanosToSeconds(queryParams.Start), convertNanosToSeconds(queryParams.End))

	err := dao.ClickhouseConnectionService.ExecuteSelectFunction(true, &serviceItems, query, nil)
	dao.Logger.Info(query)

	if err != nil {
		dao.Logger.Debug("Error in processing sql query: ", err)
		return nil, fmt.Errorf("error in processing sql query")
	}

	var serviceErrorItems []model.ServiceItem

	query = fmt.Sprintf("SELECT serviceName, sum(count) as numErrors FROM signoz_index_aggregated WHERE timestamp>='%s' AND timestamp<='%s' AND kind='2' AND (statusCode>=500 OR statusCode=2) GROUP BY serviceName", convertNanosToSeconds(queryParams.Start), convertNanosToSeconds(queryParams.End))

	err = dao.ClickhouseConnectionService.ExecuteSelectFunction(true, &serviceErrorItems, query, nil)

	dao.Logger.Info(query)

	if err != nil {
		dao.Logger.Debug("Error in processing sql query: ", err)
		return nil, fmt.Errorf("error in processing sql query")
	}

	m5xx := make(map[string]int)

	for j := range serviceErrorItems {
		m5xx[serviceErrorItems[j].ServiceName] = serviceErrorItems[j].NumErrors
	}
	///////////////////////////////////////////

	//////////////////		Below block gets 4xx of services

	var service4xxItems []model.ServiceItem

	query = fmt.Sprintf("SELECT serviceName, sum(count) as num4xx FROM signoz_index_aggregated WHERE timestamp>='%s' AND timestamp<='%s' AND kind='2' AND statusCode>=400 AND statusCode<500 GROUP BY serviceName", convertNanosToSeconds(queryParams.Start), convertNanosToSeconds(queryParams.End))

	err = dao.ClickhouseConnectionService.ExecuteSelectFunction(true, &service4xxItems, query, nil)

	dao.Logger.Info(query)

	if err != nil {
		dao.Logger.Debug("Error in processing sql query: ", err)
		return nil, fmt.Errorf("error in processing sql query")
	}

	m4xx := make(map[string]int)

	for j := range service4xxItems {
		m5xx[service4xxItems[j].ServiceName] = service4xxItems[j].Num4XX
	}

	for i, _ := range serviceItems {
		if val, ok := m5xx[serviceItems[i].ServiceName]; ok {
			serviceItems[i].NumErrors = val
		}
		if val, ok := m4xx[serviceItems[i].ServiceName]; ok {
			serviceItems[i].Num4XX = val
		}
		serviceItems[i].CallRate = float32(serviceItems[i].NumCalls) / float32(queryParams.Period)
		serviceItems[i].FourXXRate = float32(serviceItems[i].Num4XX) / float32(queryParams.Period)
		serviceItems[i].ErrorRate = float32(serviceItems[i].NumErrors) / float32(queryParams.Period)
	}

	return &serviceItems, nil
}

func (dao *ApmDaoImpl) GetServicesList(ctx context.Context) (*[]string, error) {
	var services []string

	query := fmt.Sprintf(`SELECT DISTINCT serviceName FROM signoz_index_aggregated WHERE toDate(timestamp) > now() - INTERVAL 1 DAY`)

	err := dao.ClickhouseConnectionService.ExecuteSelectFunction(true, &services, query, nil)

	dao.Logger.Info(query)

	if err != nil {
		dao.Logger.Debug("Error in processing sql query: ", err)
		return nil, fmt.Errorf("error in processing sql query")
	}

	return &services, nil
}

func (dao *ApmDaoImpl) GetServiceOverview(ctx context.Context, queryParams *model.GetServiceOverviewParams) (*[]model.ServiceOverviewItem, error) {
	var serviceOverviewItems []model.ServiceOverviewItem

	query := fmt.Sprintf("SELECT toStartOfInterval(timestamp, INTERVAL %s minute) as time, quantileMerge(0.99)(quantile) as p99, quantileMerge(0.95)(quantile) as p95,quantileMerge(0.50)(quantile) as p50, sum(count) as numCalls FROM signoz_index_aggregated WHERE timestamp>='%s' AND timestamp<='%s' AND kind='2' AND serviceName='%s' GROUP BY time ORDER BY time DESC", strconv.Itoa(queryParams.StepSeconds/60), convertNanosToSeconds(queryParams.Start), convertNanosToSeconds(queryParams.End), queryParams.ServiceName)

	err := dao.ClickhouseConnectionService.ExecuteSelectFunction(true, &serviceOverviewItems, query, nil)

	dao.Logger.Info(query)

	if err != nil {
		dao.Logger.Debug("Error in processing sql query: ", err)
		return nil, fmt.Errorf("error in processing sql query")
	}

	var serviceErrorItems []model.ServiceErrorItem

	query = fmt.Sprintf("SELECT toStartOfInterval(timestamp, INTERVAL %s minute) as time, sum(count) as numErrors FROM signoz_index_aggregated WHERE timestamp>='%s' AND timestamp<='%s' AND kind='2' AND serviceName='%s' AND (statusCode>=500 OR statusCode=2) GROUP BY time ORDER BY time DESC", strconv.Itoa(queryParams.StepSeconds/60), convertNanosToSeconds(queryParams.Start), convertNanosToSeconds(queryParams.End), queryParams.ServiceName)

	err = dao.ClickhouseConnectionService.ExecuteSelectFunction(true, &serviceErrorItems, query, nil)

	dao.Logger.Info(query)

	if err != nil {
		dao.Logger.Debug("Error in processing sql query: ", err)
		return nil, fmt.Errorf("error in processing sql query")
	}

	m := make(map[int64]int)

	for j, _ := range serviceErrorItems {
		timeObj, _ := time.Parse(time.RFC3339Nano, serviceErrorItems[j].Time)
		m[timeObj.UnixNano()] = serviceErrorItems[j].NumErrors
	}

	for i, _ := range serviceOverviewItems {
		timeObj, _ := time.Parse(time.RFC3339Nano, serviceOverviewItems[i].Time)
		serviceOverviewItems[i].Timestamp = timeObj.UnixNano()
		serviceOverviewItems[i].Time = ""

		if val, ok := m[serviceOverviewItems[i].Timestamp]; ok {
			serviceOverviewItems[i].NumErrors = val
		}
		serviceOverviewItems[i].ErrorRate = float32(serviceOverviewItems[i].NumErrors) * 100 / float32(serviceOverviewItems[i].NumCalls)
		serviceOverviewItems[i].CallRate = float32(serviceOverviewItems[i].NumCalls) / float32(queryParams.StepSeconds)
	}

	return &serviceOverviewItems, nil
}

func (dao *ApmDaoImpl) SearchSpans(ctx context.Context, queryParams *model.SpanSearchParams) (*[]model.SearchSpansResult, error) {
	query := fmt.Sprintf("SELECT timestamp, spanID, traceID, serviceName, name, kind, durationNano, tagsKeys, tagsValues FROM signoz_index_final WHERE timestamp >= ? AND timestamp <= ?")

	args := []interface{}{strconv.FormatInt(queryParams.Start.UnixNano(), 10), strconv.FormatInt(queryParams.End.UnixNano(), 10)}

	if len(queryParams.ServiceName) != 0 {
		query = query + " AND serviceName = ?"
		args = append(args, queryParams.ServiceName)
	}

	if len(queryParams.OperationName) != 0 {

		query = query + " AND name = ?"
		args = append(args, queryParams.OperationName)

	}

	if len(queryParams.Kind) != 0 {
		query = query + " AND kind = ?"
		args = append(args, queryParams.Kind)

	}

	if len(queryParams.MinDuration) != 0 {
		query = query + " AND durationNano >= ?"
		args = append(args, queryParams.MinDuration)
	}
	if len(queryParams.MaxDuration) != 0 {
		query = query + " AND durationNano <= ?"
		args = append(args, queryParams.MaxDuration)
	}

	for _, item := range queryParams.Tags {

		if item.Key == "error" && item.Value == "true" {
			query = query + " AND ( has(tags, 'error:true') OR statusCode>=500 OR statusCode=2)"
			continue
		}

		if item.Operator == "equals" {
			query = query + " AND has(tags, ?)"
			args = append(args, fmt.Sprintf("%s:%s", item.Key, item.Value))
		} else if item.Operator == "contains" {
			query = query + " AND tagsValues[indexOf(tagsKeys, ?)] ILIKE ?"
			args = append(args, item.Key)
			args = append(args, fmt.Sprintf("%%%s%%", item.Value))
		} else if item.Operator == "regex" {
			query = query + " AND match(tagsValues[indexOf(tagsKeys, ?)], ?)"
			args = append(args, item.Key)
			args = append(args, item.Value)
		} else if item.Operator == "isnotnull" {
			query = query + " AND has(tagsKeys, ?)"
			args = append(args, item.Key)
		} else {
			return nil, fmt.Errorf("tag Operator %s not supported", item.Operator)
		}

	}

	query = query + " ORDER BY timestamp DESC LIMIT 100"

	var searchScanResponses []model.SearchSpanReponseItem

	err := dao.ClickhouseConnectionService.ExecuteSelectFunction(true, &searchScanResponses, query, args)

	dao.Logger.Info(query)

	if err != nil {
		dao.Logger.Debug("Error in processing sql query: ", err)
		return nil, fmt.Errorf("error in processing sql query")
	}

	searchSpansResult := []model.SearchSpansResult{
		{
			Columns: []string{"__time", "SpanId", "TraceId", "ServiceName", "Name", "Kind", "DurationNano", "TagsKeys", "TagsValues"},
			Events:  make([][]interface{}, len(searchScanResponses)),
		},
	}

	for i, item := range searchScanResponses {
		spanEvents := item.GetValues()
		searchSpansResult[0].Events[i] = spanEvents
	}

	return &searchSpansResult, nil
}

func (dao *ApmDaoImpl) GetServiceDBOverview(ctx context.Context, queryParams *model.GetServiceOverviewParams) (*[]model.ServiceDBOverviewItem, error) {
	var serviceDBOverviewItems []model.ServiceDBOverviewItem

	query := fmt.Sprintf("SELECT toStartOfInterval(timestamp, INTERVAL %s minute) as time, avgMerge(avg) as avgDuration, sum(count) as numCalls, dbSystem FROM signoz_index_aggregated WHERE serviceName='%s' AND timestamp>='%s' AND timestamp<='%s' AND kind='3' AND dbName IS NOT NULL GROUP BY time, dbSystem ORDER BY time DESC", strconv.Itoa(queryParams.StepSeconds/60), queryParams.ServiceName, convertNanosToSeconds(queryParams.Start), convertNanosToSeconds(queryParams.End))

	err := dao.ClickhouseConnectionService.ExecuteSelectFunction(true, &serviceDBOverviewItems, query, nil)

	dao.Logger.Info(query)

	if err != nil {
		dao.Logger.Debug("Error in processing sql query: ", err)
		return nil, fmt.Errorf("error in processing sql query")
	}

	for i := range serviceDBOverviewItems {
		timeObj, _ := time.Parse(time.RFC3339Nano, serviceDBOverviewItems[i].Time)
		serviceDBOverviewItems[i].Timestamp = timeObj.UnixNano()
		serviceDBOverviewItems[i].Time = ""
		serviceDBOverviewItems[i].CallRate = float32(serviceDBOverviewItems[i].NumCalls) / float32(queryParams.StepSeconds)
	}

	if serviceDBOverviewItems == nil {
		serviceDBOverviewItems = []model.ServiceDBOverviewItem{}
	}

	return &serviceDBOverviewItems, nil
}

func (dao *ApmDaoImpl) GetServiceExternalAvgDuration(ctx context.Context, queryParams *model.GetServiceOverviewParams) (*[]model.ServiceExternalItem, error) {
	var serviceExternalItems []model.ServiceExternalItem

	query := fmt.Sprintf("SELECT toStartOfInterval(timestamp, INTERVAL %s minute) as time, avgMerge(avg) as avgDuration FROM signoz_index_aggregated WHERE serviceName='%s' AND timestamp>='%s' AND timestamp<='%s' AND kind='3' AND externalHttpUrl IS NOT NULL GROUP BY time ORDER BY time DESC", strconv.Itoa(queryParams.StepSeconds/60), queryParams.ServiceName, convertNanosToSeconds(queryParams.Start), convertNanosToSeconds(queryParams.End))

	err := dao.ClickhouseConnectionService.ExecuteSelectFunction(true, &serviceExternalItems, query, nil)

	dao.Logger.Info(query)

	if err != nil {
		dao.Logger.Debug("Error in processing sql query: ", err)
		return nil, fmt.Errorf("error in processing sql query")
	}

	for i := range serviceExternalItems {
		timeObj, _ := time.Parse(time.RFC3339Nano, serviceExternalItems[i].Time)
		serviceExternalItems[i].Timestamp = timeObj.UnixNano()
		serviceExternalItems[i].Time = ""
		serviceExternalItems[i].CallRate = float32(serviceExternalItems[i].NumCalls) / float32(queryParams.StepSeconds)
	}

	if serviceExternalItems == nil {
		serviceExternalItems = []model.ServiceExternalItem{}
	}

	return &serviceExternalItems, nil
}

func (dao *ApmDaoImpl) GetServiceExternalErrors(ctx context.Context, queryParams *model.GetServiceOverviewParams) (*[]model.ServiceExternalItem, error) {
	var serviceExternalErrorItems []model.ServiceExternalItem

	query := fmt.Sprintf("SELECT toStartOfInterval(timestamp, INTERVAL %s minute) as time, avgMerge(avg) as avgDuration, sum(count) as numCalls, externalHttpUrl FROM signoz_index_aggregated WHERE serviceName='%s' AND timestamp>='%s' AND timestamp<='%s' AND kind='3' AND externalHttpUrl IS NOT NULL AND (statusCode >= 500 OR statusCode=2) GROUP BY time, externalHttpUrl ORDER BY time DESC", strconv.Itoa(queryParams.StepSeconds/60), queryParams.ServiceName, convertNanosToSeconds(queryParams.Start), convertNanosToSeconds(queryParams.End))

	err := dao.ClickhouseConnectionService.ExecuteSelectFunction(true, &serviceExternalErrorItems, query, nil)

	dao.Logger.Info(query)

	if err != nil {
		dao.Logger.Debug("Error in processing sql query: ", err)
		return nil, fmt.Errorf("error in processing sql query")
	}
	var serviceExternalTotalItems []model.ServiceExternalItem

	queryTotal := fmt.Sprintf("SELECT toStartOfInterval(timestamp, INTERVAL %s minute) as time, avg(durationNano) as avgDuration, count(1) as numCalls, externalHttpUrl FROM signoz_index_final WHERE serviceName='%s' AND timestamp>='%s' AND timestamp<='%s' AND kind='3' AND externalHttpUrl IS NOT NULL GROUP BY time, externalHttpUrl ORDER BY time DESC", strconv.Itoa(queryParams.StepSeconds/60), queryParams.ServiceName, strconv.FormatInt(queryParams.Start.UnixNano(), 10), strconv.FormatInt(queryParams.End.UnixNano(), 10))

	errTotal := dao.ClickhouseConnectionService.ExecuteSelectFunction(true, &serviceExternalTotalItems, queryTotal, nil)

	if errTotal != nil {
		dao.Logger.Debug("Error in processing sql query: ", err)
		return nil, fmt.Errorf("error in processing sql query")
	}

	m := make(map[string]int)

	for j, _ := range serviceExternalErrorItems {
		timeObj, _ := time.Parse(time.RFC3339Nano, serviceExternalErrorItems[j].Time)
		m[strconv.FormatInt(timeObj.UnixNano(), 10)+"-"+serviceExternalErrorItems[j].ExternalHttpUrl] = serviceExternalErrorItems[j].NumCalls
	}

	for i, _ := range serviceExternalTotalItems {
		timeObj, _ := time.Parse(time.RFC3339Nano, serviceExternalTotalItems[i].Time)
		serviceExternalTotalItems[i].Timestamp = int64(timeObj.UnixNano())
		serviceExternalTotalItems[i].Time = ""
		// serviceExternalTotalItems[i].CallRate = float32(serviceExternalTotalItems[i].NumCalls) / float32(queryParams.StepSeconds)

		if val, ok := m[strconv.FormatInt(serviceExternalTotalItems[i].Timestamp, 10)+"-"+serviceExternalTotalItems[i].ExternalHttpUrl]; ok {
			serviceExternalTotalItems[i].NumErrors = val
			serviceExternalTotalItems[i].ErrorRate = float32(serviceExternalTotalItems[i].NumErrors) * 100 / float32(serviceExternalTotalItems[i].NumCalls)
		}
		serviceExternalTotalItems[i].CallRate = 0
		serviceExternalTotalItems[i].NumCalls = 0

	}

	if serviceExternalTotalItems == nil {
		serviceExternalTotalItems = []model.ServiceExternalItem{}
	}

	return &serviceExternalTotalItems, nil
}

func (dao *ApmDaoImpl) GetServiceExternal(ctx context.Context, queryParams *model.GetServiceOverviewParams) (*[]model.ServiceExternalItem, error) {
	var serviceExternalItems []model.ServiceExternalItem

	query := fmt.Sprintf("SELECT toStartOfInterval(timestamp, INTERVAL %s minute) as time, avgMerge(avg) as avgDuration, sum(count) as numCalls, externalHttpUrl FROM signoz_index_aggregated WHERE serviceName='%s' AND timestamp>='%s' AND timestamp<='%s' AND kind='3' AND externalHttpUrl IS NOT NULL GROUP BY time, externalHttpUrl ORDER BY time DESC", strconv.Itoa(queryParams.StepSeconds/60), queryParams.ServiceName, convertNanosToSeconds(queryParams.Start), convertNanosToSeconds(queryParams.End))

	err := dao.ClickhouseConnectionService.ExecuteSelectFunction(true, &serviceExternalItems, query, nil)

	dao.Logger.Info(query)

	if err != nil {
		dao.Logger.Debug("Error in processing sql query: ", err)
		return nil, fmt.Errorf("error in processing sql query")
	}

	for i, _ := range serviceExternalItems {
		timeObj, _ := time.Parse(time.RFC3339Nano, serviceExternalItems[i].Time)
		serviceExternalItems[i].Timestamp = int64(timeObj.UnixNano())
		serviceExternalItems[i].Time = ""
		serviceExternalItems[i].CallRate = float32(serviceExternalItems[i].NumCalls) / float32(queryParams.StepSeconds)
	}

	if serviceExternalItems == nil {
		serviceExternalItems = []model.ServiceExternalItem{}
	}

	return &serviceExternalItems, nil
}

func (dao *ApmDaoImpl) GetTopEndpoints(ctx context.Context, queryParams *model.GetTopEndpointsParams) (*[]model.TopEndpointsItem, error) {
	var topEndpointsItems []model.TopEndpointsItem

	query := fmt.Sprintf("SELECT quantileMerge(0.5)(quantile) as p50, quantileMerge(0.95)(quantile) as p95, quantileMerge(0.99)(quantile) as p99, sum(count) as numCalls, name  FROM signoz_index_aggregated WHERE  timestamp >= '%s' AND timestamp <= '%s' AND  kind='2' and serviceName='%s' GROUP BY name", convertNanosToSeconds(queryParams.Start), convertNanosToSeconds(queryParams.End), queryParams.ServiceName)

	err := dao.ClickhouseConnectionService.ExecuteSelectFunction(true, &topEndpointsItems, query, nil)

	dao.Logger.Info(query)

	if err != nil {
		dao.Logger.Debug("Error in processing sql query: ", err)
		return nil, fmt.Errorf("error in processing sql query")
	}

	if topEndpointsItems == nil {
		topEndpointsItems = []model.TopEndpointsItem{}
	}

	return &topEndpointsItems, nil
}

func (dao *ApmDaoImpl) GetUsage(ctx context.Context, queryParams *model.GetUsageParams) (*[]model.UsageItem, error) {
	var usageItems []model.UsageItem

	var query string
	if len(queryParams.ServiceName) != 0 {
		query = fmt.Sprintf("SELECT toStartOfInterval(timestamp, INTERVAL %d HOUR) as time, count(1) as count FROM signoz_index_final WHERE serviceName='%s' AND timestamp>='%s' AND timestamp<='%s' GROUP BY time ORDER BY time ASC", queryParams.StepHour, queryParams.ServiceName, strconv.FormatInt(queryParams.Start.UnixNano(), 10), strconv.FormatInt(queryParams.End.UnixNano(), 10))
	} else {
		query = fmt.Sprintf("SELECT toStartOfInterval(timestamp, INTERVAL %d HOUR) as time, count(1) as count FROM signoz_index_final WHERE timestamp>='%s' AND timestamp<='%s' GROUP BY time ORDER BY time ASC", queryParams.StepHour, strconv.FormatInt(queryParams.Start.UnixNano(), 10), strconv.FormatInt(queryParams.End.UnixNano(), 10))
	}

	err := dao.ClickhouseConnectionService.ExecuteSelectFunction(true, &usageItems, query, nil)

	dao.Logger.Info(query)

	if err != nil {
		dao.Logger.Debug("Error in processing sql query: ", err)
		return nil, fmt.Errorf("error in processing sql query")
	}

	for i := range usageItems {
		timeObj, _ := time.Parse(time.RFC3339Nano, usageItems[i].Time)
		usageItems[i].Timestamp = int64(timeObj.UnixNano())
		usageItems[i].Time = ""
	}

	if usageItems == nil {
		usageItems = []model.UsageItem{}
	}

	return &usageItems, nil
}

func (dao *ApmDaoImpl) GetOperations(ctx context.Context, serviceName string) (*[]string, error) {
	var operations []string

	query := fmt.Sprintf(`SELECT DISTINCT(name) FROM signoz_index_aggregated WHERE serviceName='%s'  AND toDate(timestamp) > now() - INTERVAL 1 DAY`, serviceName)

	err := dao.ClickhouseConnectionService.ExecuteSelectFunction(true, &operations, query, nil)

	dao.Logger.Info(query)

	if err != nil {
		dao.Logger.Debug("Error in processing sql query: ", err)
		return nil, fmt.Errorf("error in processing sql query")
	}
	return &operations, nil
}

func (dao *ApmDaoImpl) GetTags(ctx context.Context, serviceName string) (*[]model.TagItem, error) {
	var tagItems []model.TagItem

	query := fmt.Sprintf(`SELECT DISTINCT arrayJoin(tagsKeys) as tagKeys FROM signoz_index_aggregated WHERE serviceName='%s'  AND toDate(timestamp) > now() - INTERVAL 1 DAY`, serviceName)

	err := dao.ClickhouseConnectionService.ExecuteSelectFunction(true, &tagItems, query, nil)

	dao.Logger.Info(query)

	if err != nil {
		dao.Logger.Debug("Error in processing sql query: ", err)
		return nil, fmt.Errorf("error in processing sql query")
	}

	return &tagItems, nil
}

func (dao *ApmDaoImpl) SearchTraces(ctx context.Context, traceID string) (*[]model.SearchSpansResult, error) {
	var searchScanResponses []model.SearchSpanReponseItem

	query := fmt.Sprintf("SELECT timestamp, spanID, traceID, serviceName, name, kind, durationNano, tagsKeys, tagsValues, references FROM signoz_index_final WHERE traceID='%s'", traceID)

	err := dao.ClickhouseConnectionService.ExecuteSelectFunction(true, &searchScanResponses, query, nil)

	dao.Logger.Info(query)

	if err != nil {
		dao.Logger.Debug("Error in processing sql query: ", err)
		return nil, fmt.Errorf("error in processing sql query")
	}

	searchSpansResult := []model.SearchSpansResult{
		{
			Columns: []string{"__time", "SpanId", "TraceId", "ServiceName", "Name", "Kind", "DurationNano", "TagsKeys", "TagsValues", "References"},
			Events:  make([][]interface{}, len(searchScanResponses)),
		},
	}

	for i, item := range searchScanResponses {
		spanEvents := item.GetValues()
		searchSpansResult[0].Events[i] = spanEvents
	}

	return &searchSpansResult, nil
}

func (dao *ApmDaoImpl) GetServiceMapDependencies(ctx context.Context, queryParams *model.GetServicesParams) (*[]model.ServiceMapDependencyResponseItem, error) {
	var serviceMapDependencyItems []model.ServiceMapDependencyItem

	query := fmt.Sprintf(`SELECT spanID, parentSpanID, serviceName FROM signoz_index_final WHERE timestamp>='%s' AND timestamp<='%s'`, strconv.FormatInt(queryParams.Start.UnixNano(), 10), strconv.FormatInt(queryParams.End.UnixNano(), 10))

	err := dao.ClickhouseConnectionService.ExecuteSelectFunction(true, &serviceMapDependencyItems, query, nil)

	zap.S().Info(query)

	if err != nil {
		zap.S().Debug("Error in processing sql query: ", err)
		return nil, fmt.Errorf("error in processing sql query")
	}

	serviceMap := make(map[string]*model.ServiceMapDependencyResponseItem)

	spanId2ServiceNameMap := make(map[string]string)
	for i, _ := range serviceMapDependencyItems {
		spanId2ServiceNameMap[serviceMapDependencyItems[i].SpanId] = serviceMapDependencyItems[i].ServiceName
	}
	for i, _ := range serviceMapDependencyItems {
		parent2childServiceName := spanId2ServiceNameMap[serviceMapDependencyItems[i].ParentSpanId] + "-" + spanId2ServiceNameMap[serviceMapDependencyItems[i].SpanId]
		if _, ok := serviceMap[parent2childServiceName]; !ok {
			serviceMap[parent2childServiceName] = &model.ServiceMapDependencyResponseItem{
				Parent:    spanId2ServiceNameMap[serviceMapDependencyItems[i].ParentSpanId],
				Child:     spanId2ServiceNameMap[serviceMapDependencyItems[i].SpanId],
				CallCount: 1,
			}
		} else {
			serviceMap[parent2childServiceName].CallCount++
		}
	}

	retMe := make([]model.ServiceMapDependencyResponseItem, 0, len(serviceMap))
	for _, dependency := range serviceMap {
		if dependency.Parent == "" {
			continue
		}
		retMe = append(retMe, *dependency)
	}

	return &retMe, nil
}

func (dao *ApmDaoImpl) SearchSpansAggregate(ctx context.Context, queryParams *model.SpanSearchAggregatesParams) ([]model.SpanSearchAggregatesResponseItem, error) {
	var spanSearchAggregatesResponseItems []model.SpanSearchAggregatesResponseItem

	aggregationQuery := ""
	if queryParams.Dimension == "duration" {
		switch queryParams.AggregationOption {
		case "p50":
			aggregationQuery = " quantileMerge(0.50)(quantile) as value "
			break

		case "p95":
			aggregationQuery = " quantileMerge(0.95)(quantile) as value "
			break

		case "p99":
			aggregationQuery = " quantileMerge(0.99)(quantile) as value "
			break
		}
	} else if queryParams.Dimension == "calls" {
		aggregationQuery = " sum(count) as value "
	}

	query := fmt.Sprintf("SELECT toStartOfInterval(timestamp, INTERVAL %d minute) as time, %s FROM signoz_index_aggregated WHERE timestamp >= ? AND timestamp <= ?", queryParams.StepSeconds/60, aggregationQuery)

	args := []interface{}{convertNanosToSeconds(queryParams.Start), convertNanosToSeconds(queryParams.End)}

	if len(queryParams.ServiceName) != 0 {
		query = query + " AND serviceName = ?"
		args = append(args, queryParams.ServiceName)
	}

	if len(queryParams.OperationName) != 0 {

		query = query + " AND name = ?"
		args = append(args, queryParams.OperationName)

	}

	if len(queryParams.Kind) != 0 {
		query = query + " AND kind = ?"
		args = append(args, queryParams.Kind)

	}

	query = query + " GROUP BY time ORDER BY time"

	err := dao.ClickhouseConnectionService.ExecuteSelectFunction(true, &spanSearchAggregatesResponseItems, query, args)

	dao.Logger.Info(query)

	if err != nil {
		dao.Logger.Debug("Error in processing sql query: ", err)
		return nil, fmt.Errorf("error in processing sql query")
	}

	for i := range spanSearchAggregatesResponseItems {

		timeObj, _ := time.Parse(time.RFC3339Nano, spanSearchAggregatesResponseItems[i].Time)
		spanSearchAggregatesResponseItems[i].Timestamp = int64(timeObj.UnixNano())
		spanSearchAggregatesResponseItems[i].Time = ""
		if queryParams.AggregationOption == "rate_per_sec" {
			spanSearchAggregatesResponseItems[i].Value = float32(spanSearchAggregatesResponseItems[i].Value) / float32(queryParams.StepSeconds)
		}
	}

	return spanSearchAggregatesResponseItems, nil
}

func convertNanosToSeconds(timeToParse *time.Time) string {
	seconds := int64(time.Second) / int64(time.Nanosecond)
	return strconv.Itoa(int(timeToParse.UnixNano() / seconds))
}
