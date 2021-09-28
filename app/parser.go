package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	model "goapm/domain"
	"strconv"
	"time"
)

func parseGetServiceOverviewRequest(ctx *fiber.Ctx) (*model.GetServiceOverviewParams, error) {
	startTime, err := parseTime(ctx.Query("start"))
	if err != nil {
		return nil, err
	}
	endTime, err := parseTime(ctx.Query("end"))
	if err != nil {
		return nil, err
	}

	stepStr := ctx.Query("step")
	if len(stepStr) == 0 {
		return nil, errors.New("step param missing in query")
	}
	stepInt, err := strconv.Atoi(stepStr)
	if err != nil {
		return nil, errors.New("step param is not in correct format")
	}

	serviceName := ctx.Query("service")
	if len(serviceName) == 0 {
		return nil, errors.New("serviceName param missing in query")
	}

	getServiceOverviewParams := model.GetServiceOverviewParams{
		Start:       startTime,
		StartTime:   startTime.Format(time.RFC3339Nano),
		End:         endTime,
		EndTime:     endTime.Format(time.RFC3339Nano),
		ServiceName: serviceName,
		Period:      fmt.Sprintf("PT%dM", stepInt/60),
		StepSeconds: stepInt,
	}

	return &getServiceOverviewParams, nil

}

func parseGetServicesRequest(ctx *fiber.Ctx) (*model.GetServicesParams, error) {

	startTime, err := parseTime(ctx.Query("start"))
	if err != nil {
		return nil, err
	}
	endTime, err := parseTime(ctx.Query("end"))
	if err != nil {
		return nil, err
	}

	getServicesParams := model.GetServicesParams{
		Start:     startTime,
		StartTime: startTime.Format(time.RFC3339Nano),
		End:       endTime,
		EndTime:   endTime.Format(time.RFC3339Nano),
		Period:    int(endTime.Unix() - startTime.Unix()),
	}
	return &getServicesParams, nil

}

func parseGetTopEndpointsRequest(ctx *fiber.Ctx) (*model.GetTopEndpointsParams, error) {
	startTime, err := parseTime(ctx.Query("start"))
	if err != nil {
		return nil, err
	}
	endTime, err := parseTime(ctx.Query("end"))
	if err != nil {
		return nil, err
	}

	serviceName := ctx.Query("service")
	if len(serviceName) == 0 {
		return nil, errors.New("serviceName param missing in query")
	}

	getTopEndpointsParams := model.GetTopEndpointsParams{
		StartTime:   startTime.Format(time.RFC3339Nano),
		EndTime:     endTime.Format(time.RFC3339Nano),
		ServiceName: serviceName,
		Start:       startTime,
		End:         endTime,
	}

	return &getTopEndpointsParams, nil

}

func parseSpanSearchRequest(ctx *fiber.Ctx) (*model.SpanSearchParams, error) {

	startTime, err := parseTime(ctx.Query("start"))
	if err != nil {
		return nil, err
	}
	endTime, err := parseTimeMinusBuffer(ctx.Query("end"))
	if err != nil {
		return nil, err
	}

	startTimeStr := startTime.Format(time.RFC3339Nano)
	endTimeStr := endTime.Format(time.RFC3339Nano)
	// fmt.Println(startTimeStr)
	params := &model.SpanSearchParams{
		Intervals: fmt.Sprintf("%s/%s", startTimeStr, endTimeStr),
		Start:     startTime,
		End:       endTime,
		Limit:     100,
		Order:     "descending",
	}

	serviceName := ctx.Query("service")
	if len(serviceName) != 0 {
		// return nil, errors.New("serviceName param missing in query")
		params.ServiceName = serviceName
	}
	operationName := ctx.Query("operation")
	if len(operationName) != 0 {
		params.OperationName = operationName
		zap.S().Debug("Operation Name: ", operationName)
	}

	kind := ctx.Query("kind")
	if len(kind) != 0 {
		params.Kind = kind
		zap.S().Debug("Kind: ", kind)
	}

	minDuration, err := parseTimestamp(ctx.Query("minDuration"))
	if err == nil {
		params.MinDuration = *minDuration
	}
	maxDuration, err := parseTimestamp(ctx.Query("maxDuration"))
	if err == nil {
		params.MaxDuration = *maxDuration
	}

	limitStr := ctx.Query("limit")
	if len(limitStr) != 0 {
		limit, err := strconv.ParseInt(limitStr, 10, 64)
		if err != nil {
			return nil, errors.New("Limit param is not in correct format")
		}
		params.Limit = limit
	} else {
		params.Limit = 100
	}

	offsetStr := ctx.Query("offset")
	if len(offsetStr) != 0 {
		offset, err := strconv.ParseInt(offsetStr, 10, 64)
		if err != nil {
			return nil, errors.New("Offset param is not in correct format")
		}
		params.Offset = offset
	}

	tags, err := parseTags(ctx.Query("tags"))
	if err != nil {
		return nil, err
	}
	if len(*tags) != 0 {
		params.Tags = *tags
	}

	return params, nil
}

func DoesExistInSlice(item string, list []string) bool {
	for _, element := range list {
		if item == element {
			return true
		}
	}
	return false
}

func parseSearchSpanAggregatesRequest(ctx *fiber.Ctx) (*model.SpanSearchAggregatesParams, error) {

	var allowedDimensions = []string{"calls", "duration"}

	var allowedAggregations = map[string][]string{
		"calls":    {"count", "rate_per_sec"},
		"duration": {"avg", "p50", "p95", "p99"},
	}

	startTime, err := parseTime(ctx.Query("start"))
	if err != nil {
		return nil, err
	}
	endTime, err := parseTime(ctx.Query("end"))
	if err != nil {
		return nil, err
	}

	startTimeStr := startTime.Format(time.RFC3339Nano)
	endTimeStr := endTime.Format(time.RFC3339Nano)
	// fmt.Println(startTimeStr)

	stepStr := ctx.Query("step")
	if len(stepStr) == 0 {
		return nil, errors.New("step param missing in query")
	}

	stepInt, err := strconv.Atoi(stepStr)
	if err != nil {
		return nil, errors.New("step param is not in correct format")
	}

	granPeriod := fmt.Sprintf("PT%dM", stepInt/60)
	dimension := ctx.Query("dimension")
	if len(dimension) == 0 {
		return nil, errors.New("dimension param missing in query")
	} else {
		if !DoesExistInSlice(dimension, allowedDimensions) {
			return nil, errors.New(fmt.Sprintf("given dimension: %s is not allowed in query", dimension))
		}
	}

	aggregationOption := ctx.Query("aggregation_option")
	if len(aggregationOption) == 0 {
		return nil, errors.New("Aggregation Option missing in query params")
	} else {
		if !DoesExistInSlice(aggregationOption, allowedAggregations[dimension]) {
			return nil, errors.New(fmt.Sprintf("given aggregation option: %s is not allowed with dimension: %s", aggregationOption, dimension))
		}
	}

	params := &model.SpanSearchAggregatesParams{
		Start:             startTime,
		End:               endTime,
		Intervals:         fmt.Sprintf("%s/%s", startTimeStr, endTimeStr),
		GranOrigin:        startTimeStr,
		GranPeriod:        granPeriod,
		StepSeconds:       stepInt,
		Dimension:         dimension,
		AggregationOption: aggregationOption,
	}

	serviceName := ctx.Query("service")
	if len(serviceName) != 0 {
		// return nil, errors.New("serviceName param missing in query")
		params.ServiceName = serviceName
	}
	operationName := ctx.Query("operation")
	if len(operationName) != 0 {
		params.OperationName = operationName
		zap.S().Debug("Operation Name: ", operationName)
	}

	kind := ctx.Query("kind")
	if len(kind) != 0 {
		params.Kind = kind
		zap.S().Debug("Kind: ", kind)
	}

	minDuration, err := parseTimestamp(ctx.Query("minDuration"))
	if err == nil {
		params.MinDuration = *minDuration
	}
	maxDuration, err := parseTimestamp(ctx.Query("maxDuration"))
	if err == nil {
		params.MaxDuration = *maxDuration
	}

	tags, err := parseTags(ctx.Query("tags"))
	if err != nil {
		return nil, err
	}
	if len(*tags) != 0 {
		params.Tags = *tags
	}

	return params, nil
}

func parseTime(timeStr string) (*time.Time, error) {

	if len(timeStr) == 0 {
		return nil, fmt.Errorf("%s param missing in query", timeStr)
	}

	timeUnix, err := strconv.ParseInt(timeStr, 10, 64)
	if err != nil || len(timeStr) == 0 {
		return nil, fmt.Errorf("%s param is not in correct timestamp format", timeStr)
	}

	timeFmt := time.Unix(0, timeUnix).Add(-5 * time.Minute)

	return &timeFmt, nil

}

func parseGetUsageRequest(ctx *fiber.Ctx) (*model.GetUsageParams, error) {
	startTime, err := parseTime(ctx.Query("start"))
	if err != nil {
		return nil, err
	}
	endTime, err := parseTime(ctx.Query("end"))
	if err != nil {
		return nil, err
	}

	stepStr := ctx.Query("step")
	if len(stepStr) == 0 {
		return nil, errors.New("step param missing in query")
	}
	stepInt, err := strconv.Atoi(stepStr)
	if err != nil {
		return nil, errors.New("step param is not in correct format")
	}

	serviceName := ctx.Query("service")
	stepHour := stepInt / 3600

	getUsageParams := model.GetUsageParams{
		StartTime:   startTime.Format(time.RFC3339Nano),
		EndTime:     endTime.Format(time.RFC3339Nano),
		Start:       startTime,
		End:         endTime,
		ServiceName: serviceName,
		Period:      fmt.Sprintf("PT%dH", stepHour),
		StepHour:    stepHour,
	}

	return &getUsageParams, nil

}

func parseTimeMinusBuffer(timeStr string) (*time.Time, error) {

	if len(timeStr) == 0 {
		return nil, fmt.Errorf("%s param missing in query", timeStr)
	}

	timeUnix, err := strconv.ParseInt(timeStr, 10, 64)
	if err != nil || len(timeStr) == 0 {
		return nil, fmt.Errorf("%s param is not in correct timestamp format", timeStr)
	}

	timeUnixNow := time.Now().UnixNano()
	if timeUnix > timeUnixNow-30000000000 {
		timeUnix = timeUnix - 30000000000
	}

	timeFmt := time.Unix(0, timeUnix)

	return &timeFmt, nil

}

func parseTimestamp(timeStr string) (*string, error) {
	if len(timeStr) == 0 {
		return nil, fmt.Errorf("%s param missing in query", timeStr)
	}
	return &timeStr, nil

}

func parseTags(tagsStr string) (*[]model.TagQuery, error) {

	tags := new([]model.TagQuery)

	if len(tagsStr) == 0 {
		return tags, nil
	}
	err := json.Unmarshal([]byte(tagsStr), tags)
	if err != nil {
		zap.S().Error("Error in parsig tags", zap.Error(err))
		return nil, fmt.Errorf("error in parsing %s ", tagsStr)
	}
	// zap.S().Info("Tags: ", *tags)

	return tags, nil
}
