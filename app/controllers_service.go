package controllers

import (
	"goapm/logger"
	"goapm/web"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"github.com/valyala/fasthttp"
	model "goapm/domain"
	"strings"
	"time"
)

func InitControllers(appName string, version string, build string, commit string, buildTime string) {

	app := web.InitApp(appName, version, build, commit, buildTime)
	InitHealthCheck()
	InitServices()

	ticker := time.NewTicker(1 * time.Minute)
	go func() {
		for range ticker.C {
			go traceFilterJob.Run()
		}
	}()

	app.Use(func(c *fiber.Ctx) error {
		c.Path(strings.ReplaceAll(c.Path(), "//", "/"))
		return c.Next()
	})

	app.Get("/api/v1/services", func(ctx *fiber.Ctx) error {
		query, err := parseGetServicesRequest(ctx)
		if err != nil {
			return ctx.Status(fasthttp.StatusBadRequest).SendString(fasthttp.StatusMessage(fasthttp.StatusBadRequest))
		}

		result, err := apmService.GetServices(ctx.UserContext(), query)
		if err != nil {
			return ctx.Status(fasthttp.StatusBadRequest).SendString(fasthttp.StatusMessage(fasthttp.StatusBadRequest))
		}

		if len(*result) == 0 {
			return ctx.JSON([]model.ServiceItem{})
		}
		return ctx.JSON(result)
	})

	app.Get("/api/v1/services/list", func(ctx *fiber.Ctx) error {
		result, err := apmService.GetServicesList(ctx.UserContext())
		if err != nil {
			return ctx.Status(fasthttp.StatusBadRequest).SendString(fasthttp.StatusMessage(fasthttp.StatusBadRequest))
		}
		return ctx.JSON(result)
	})

	app.Get("/api/v1/service/overview", func(ctx *fiber.Ctx) error {
		query, err := parseGetServiceOverviewRequest(ctx)
		if err != nil {
			return ctx.Status(fasthttp.StatusBadRequest).SendString(fasthttp.StatusMessage(fasthttp.StatusBadRequest))
		}
		result, err := apmService.GetServiceOverview(ctx.UserContext(), query)
		if err != nil {
			return ctx.Status(fasthttp.StatusBadRequest).SendString(fasthttp.StatusMessage(fasthttp.StatusBadRequest))
		}

		if len(*result) == 0 {
			return ctx.JSON([]model.ServiceOverviewItem{})
		}
		return ctx.JSON(result)
	})

	app.Get("/api/v1/service/dbOverview", func(ctx *fiber.Ctx) error {
		query, err := parseGetServiceOverviewRequest(ctx)
		if err != nil {
			return ctx.Status(fasthttp.StatusBadRequest).SendString(fasthttp.StatusMessage(fasthttp.StatusBadRequest))
		}
		result, err := apmService.GetServiceDBOverview(ctx.UserContext(), query)
		if err != nil {
			return ctx.Status(fasthttp.StatusBadRequest).SendString(fasthttp.StatusMessage(fasthttp.StatusBadRequest))
		}

		if len(*result) == 0 {
			return ctx.JSON([]model.ServiceDBOverviewItem{})
		}
		return ctx.JSON(result)
	})

	app.Get("/api/v1/service/externalAvgDuration", func(ctx *fiber.Ctx) error {
		query, err := parseGetServiceOverviewRequest(ctx)
		if err != nil {
			return ctx.Status(fasthttp.StatusBadRequest).SendString(fasthttp.StatusMessage(fasthttp.StatusBadRequest))
		}
		result, err := apmService.GetServiceExternalAvgDuration(ctx.UserContext(), query)
		if err != nil {
			return ctx.Status(fasthttp.StatusBadRequest).SendString(fasthttp.StatusMessage(fasthttp.StatusBadRequest))
		}

		if len(*result) == 0 {
			return ctx.JSON([]model.ServiceExternalItem{})
		}
		return ctx.JSON(result)
	})

	app.Get("/api/v1/service/externalErrors", func(ctx *fiber.Ctx) error {
		query, err := parseGetServiceOverviewRequest(ctx)
		if err != nil {
			return ctx.Status(fasthttp.StatusBadRequest).SendString(fasthttp.StatusMessage(fasthttp.StatusBadRequest))
		}
		result, err := apmService.GetServiceExternalErrors(ctx.UserContext(), query)
		if err != nil {
			return ctx.Status(fasthttp.StatusBadRequest).SendString(fasthttp.StatusMessage(fasthttp.StatusBadRequest))
		}
		if len(*result) == 0 {
			return ctx.JSON([]model.ServiceExternalItem{})
		}
		return ctx.JSON(result)
	})

	app.Get("/api/v1/service/external", func(ctx *fiber.Ctx) error {
		query, err := parseGetServiceOverviewRequest(ctx)
		if err != nil {
			return ctx.Status(fasthttp.StatusBadRequest).SendString(fasthttp.StatusMessage(fasthttp.StatusBadRequest))
		}
		result, err := apmService.GetServiceExternal(ctx.UserContext(), query)
		if err != nil {
			return ctx.Status(fasthttp.StatusBadRequest).SendString(fasthttp.StatusMessage(fasthttp.StatusBadRequest))
		}
		if len(*result) == 0 {
			return ctx.JSON([]model.ServiceExternalItem{})
		}
		return ctx.JSON(result)
	})

	app.Get("/api/v1/service/:service/operations", func(ctx *fiber.Ctx) error {
		serviceName := ctx.Params("service")
		var err error
		if len(serviceName) == 0 {
			err = fmt.Errorf("service param not found")
			return ctx.Status(fasthttp.StatusBadRequest).SendString(fasthttp.StatusMessage(fasthttp.StatusBadRequest))
		}
		result, err := apmService.GetOperations(ctx.UserContext(), serviceName)
		if err != nil {
			return ctx.Status(fasthttp.StatusBadRequest).SendString(fasthttp.StatusMessage(fasthttp.StatusBadRequest))
		}

		if len(*result) == 0 {
			return ctx.JSON([]string{})
		}
		return ctx.JSON(result)
	})

	app.Get("/api/v1/service/top_endpoints", func(ctx *fiber.Ctx) error {
		query, err := parseGetTopEndpointsRequest(ctx)
		if err != nil {
			return ctx.Status(fasthttp.StatusBadRequest).SendString(fasthttp.StatusMessage(fasthttp.StatusBadRequest))
		}
		result, err := apmService.GetTopEndpoints(ctx.UserContext(), query)
		if err != nil {
			return ctx.Status(fasthttp.StatusBadRequest).SendString(fasthttp.StatusMessage(fasthttp.StatusBadRequest))
		}

		if len(*result) == 0 {
			return ctx.JSON([]model.TopEndpointsItem{})
		}
		return ctx.JSON(result)
	})

	app.Get("/api/v1/spans", func(ctx *fiber.Ctx) error {
		query, err := parseSpanSearchRequest(ctx)
		if err != nil {
			return ctx.Status(fasthttp.StatusBadRequest).SendString(fasthttp.StatusMessage(fasthttp.StatusBadRequest))
		}
		result, err := apmService.SearchSpans(ctx.UserContext(), query)
		if err != nil {
			return ctx.Status(fasthttp.StatusBadRequest).SendString(fasthttp.StatusMessage(fasthttp.StatusBadRequest))
		}

		if len(*result) == 0 {
			return ctx.JSON([]model.SearchSpansResult{})
		}
		return ctx.JSON(result)
	})

	app.Get("/api/v1/spans/aggregates", func(ctx *fiber.Ctx) error {
		query, err := parseSearchSpanAggregatesRequest(ctx)
		if err != nil {
			return ctx.Status(fasthttp.StatusBadRequest).SendString(fasthttp.StatusMessage(fasthttp.StatusBadRequest))
		}
		result, err := apmService.SearchSpansAggregate(ctx.UserContext(), query)
		if err != nil {
			return ctx.Status(fasthttp.StatusBadRequest).SendString(fasthttp.StatusMessage(fasthttp.StatusBadRequest))
		}

		if len(result) == 0 {
			return ctx.JSON([]model.SpanSearchAggregatesResponseItem{})
		}
		return ctx.JSON(result)
	})

	app.Get("/api/v1/tags", func(ctx *fiber.Ctx) error {
		serviceName := ctx.Query("service")
		result, err := apmService.GetTags(ctx.UserContext(), serviceName)
		if err != nil {
			return ctx.Status(fasthttp.StatusBadRequest).SendString(fasthttp.StatusMessage(fasthttp.StatusBadRequest))
		}

		if len(*result) == 0 {
			return ctx.JSON([]model.TagItem{})
		}
		return ctx.JSON(result)
	})

	app.Get("/api/v1/traces/:traceID", func(ctx *fiber.Ctx) error {
		traceId := ctx.Params("traceID")
		result, err := apmService.SearchTraces(ctx.UserContext(), traceId)
		if err != nil {
			return ctx.Status(fasthttp.StatusBadRequest).SendString(fasthttp.StatusMessage(fasthttp.StatusBadRequest))
		}

		if len(*result) == 0 {
			return ctx.JSON([]model.SearchSpansResult{})
		}
		return ctx.JSON(result)
	})

	app.Get("/api/v1/usage", func(ctx *fiber.Ctx) error {
		query, err := parseGetUsageRequest(ctx)
		if err != nil {
			return ctx.Status(fasthttp.StatusBadRequest).SendString(fasthttp.StatusMessage(fasthttp.StatusBadRequest))
		}
		result, err := apmService.GetUsage(ctx.UserContext(), query)
		if err != nil {
			return ctx.Status(fasthttp.StatusBadRequest).SendString(fasthttp.StatusMessage(fasthttp.StatusBadRequest))
		}

		if len(*result) == 0 {
			return ctx.JSON([]model.UsageItem{})
		}
		return ctx.JSON(result)
	})

	app.Get("/api/v1/serviceMapDependencies", func(ctx *fiber.Ctx) error {
		query, err := parseGetServicesRequest(ctx)
		if err != nil {
			return ctx.Status(fasthttp.StatusBadRequest).SendString(fasthttp.StatusMessage(fasthttp.StatusBadRequest))
		}
		result, err := apmService.GetServiceMapDependencies(ctx.UserContext(), query)
		if err != nil {
			return ctx.Status(fasthttp.StatusBadRequest).SendString(fasthttp.StatusMessage(fasthttp.StatusBadRequest))
		}
		if len(*result) == 0 {
			return ctx.JSON([]model.GetServicesParams{})
		}
		return ctx.JSON(result)
	})

	logger.LOGGER.Fatal(app.Listen(":" + viper.GetString("server.port")).Error())
	//logger.LOGGER.Fatal(app.Listen(":7778").Error())
}
