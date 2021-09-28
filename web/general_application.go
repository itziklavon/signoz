package web

import (
	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	recover2 "github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/spf13/viper"
	"goapm/config"
	"goapm/http"
	"goapm/logger"
	"goapm/utils"
	web_filters "goapm/web/filters"
	"os"
	"runtime"
	"strconv"
	"time"
)

var HealthChecksToRun = make(map[string]utils.HealthCheckService)

// InitApp create new web application using - (https://github.com/gofiber/fiber)
// Initialize logger, config server properties and the following routes -
// /actuator/prometheus - metrics collector - using (https://github.com/ansrivas/fiberprometheus)
// /actuator/info - need to pass values of version, build, commit, build time in order to display the correct build info
// /actuator/health - check health of selected services, some defaults healthcheck are automatically registered -
// MS - sql/redis/clickhouse/mongo (first connector in map)
// Multi Tenant - sql, redis (brand - HEALTH_CHECK_BRAND property)
// /actuator/loggers - changes the log level to chosen one, if we choose to create specific logger for specific service, we can only change it
// /actuator/refresh - reload config server properties
func InitApp(appName, version, build, commit, buildTime string) *fiber.App {
	logger.InitLogger()

	profile := "prod"
	configBranch := "master"

	if len(os.Getenv("profile")) > 0 {
		profile = os.Getenv("profile")
	}

	if len(os.Getenv("configBranch")) > 0 {
		configBranch = os.Getenv("configBranch")
	}

	viper.Set("profile", profile)
	viper.Set("configServerUrl", os.Getenv("confsrvDomain"))
	viper.Set("configBranch", configBranch)

	restClient := http.NewRestClient()
	confService := config.NewConfSrvService(restClient)
	confService.LoadConfigurationFromBranch(
		viper.GetString("configServerUrl"),
		appName,
		viper.GetString("profile"),
		viper.GetString("configBranch"))
	logger.SetLevel(viper.GetString("LOG_LEVEL"))

	location, err := time.LoadLocation("UTC")
	if err == nil {
		time.Local = location
	}
	app := fiber.New()
	app.Use(recover2.New())
	app.Use(web_filters.NewXssFilter())
	app.Use(pprof.New())
	prometheus := fiberprometheus.New(appName)
	prometheus.RegisterAt(app, "/actuator/prometheus")
	app.Use(prometheus.Middleware)
	app.Use(cors.New(cors.Config{
		AllowMethods:  "GET,PUT,POST,DELETE,OPTIONS",
		AllowHeaders:  "*",
		AllowOrigins:  "*",
		ExposeHeaders: "x-auth-token,Date,jwt-auth",
	}))

	app.Get("/actuator/info", func(ctx *fiber.Ctx) error {
		jsonMap := make(map[string]interface{})
		buildMap := make(map[string]interface{})
		exectuable, _ := os.Executable()
		buildMap["executable"] = exectuable
		buildMap["build"] = build
		buildMap["version"] = version
		buildMap["buid_time"] = buildTime
		buildMap["commit"] = commit
		buildMap["go_version"] = runtime.Version()
		jsonMap["build_info"] = buildMap
		jsonMap["timestamp"] = time.Now().UTC()
		jsonMap["go_routines_num"] = runtime.NumGoroutine()
		jsonMap["vcpu"] = runtime.NumCPU()
		jsonMap["memory_stats"] = getMemUsage()
		return ctx.JSON(jsonMap)
	})

	app.Get("/actuator/health", func(ctx *fiber.Ctx) error {
		var servicesHealthResponses []utils.ServiceHealth
		for _, healthRequest := range HealthChecksToRun {
			servicesHealthResponses = append(servicesHealthResponses, healthRequest.CheckService(ctx.Context()))
		}
		healthCheckResponse := utils.ConstructHealthCheckResponse(servicesHealthResponses...)
		return ctx.Status(healthCheckResponse.Status).JSON(healthCheckResponse)
	})

	app.Post("/actuator/loggers/:level", func(ctx *fiber.Ctx) error {
		level := ctx.Params("level", "error")
		loggerName := ctx.Query("logger_name")
		logger.SetLevel(level, loggerName)
		x := make(map[string]interface{})
		x["logger"] = level
		return ctx.JSON(x)
	})

	app.Post("/actuator/refresh", func(ctx *fiber.Ctx) error {
		confService.LoadConfigurationFromBranch(
			viper.GetString("configServerUrl"),
			appName,
			viper.GetString("profile"),
			viper.GetString("configBranch"))
		logger.SetLevel(viper.GetString("LOG_LEVEL"))
		return ctx.SendString("Done....")
	})

	return app
}

func getMemUsage() map[string]interface{} {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	mapToReturn := make(map[string]interface{})
	mapToReturn["alloc"] = strconv.Itoa(int(bToMb(m.Alloc))) + " MB"
	mapToReturn["total_alloc"] = strconv.Itoa(int(bToMb(m.TotalAlloc))) + " MB"
	mapToReturn["sys"] = strconv.Itoa(int(bToMb(m.Sys))) + " MB"
	mapToReturn["num_gc"] = m.NumGC
	return mapToReturn
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
