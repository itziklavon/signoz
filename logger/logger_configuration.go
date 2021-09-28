package logger

import (
	"goapm/ds_utils"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"strings"
)

var LOGGER = InitLogger()
var logLevels = ds_utils.NewSyncedMap()

// NewLoggerWithName init logger with specific name and additional fileds, see example in test
func NewLoggerWithName(name string, fileds ...zap.Field) *zap.SugaredLogger {
	pe := zap.NewProductionEncoderConfig()
	pe.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
	pe.FunctionKey = "F"
	consoleEncoder := zapcore.NewJSONEncoder(pe)
	atomicLevel := zap.NewAtomicLevel()
	logLevels.Put(name, atomicLevel)
	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), atomicLevel),
	)
	logger := zap.New(core, zap.AddCaller())
	logger.Named(name)
	if fileds != nil && len(fileds) > 0 {
		logger.With(fileds...)
	}
	return logger.Sugar()
}

// InitLogger init base logger, returns new instance
func InitLogger() *zap.SugaredLogger {
	pe := zap.NewProductionEncoderConfig()
	pe.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
	pe.FunctionKey = "F"
	consoleEncoder := zapcore.NewJSONEncoder(pe)
	atomicLevel := zap.NewAtomicLevel()
	logLevels.Put("ROOT", atomicLevel)
	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), atomicLevel),
	)
	logger := zap.New(core, zap.AddCaller())
	return logger.Sugar()
}

// SetLevel set different log level can send specific logger name,
//if none is sent,
//then sets root log level
func SetLevel(logLevel string, loggerName ...string) {
	atomicLevelInterface, _ := logLevels.Get("ROOT")
	atomicLevel := atomicLevelInterface.(zap.AtomicLevel)
	if len(loggerName) > 0 && len(loggerName[0]) > 0 {
		atomicLevelInterface, _ = logLevels.Get(loggerName[0])
		atomicLevel = atomicLevelInterface.(zap.AtomicLevel)
	}
	if strings.EqualFold("debug", logLevel) {
		atomicLevel.SetLevel(zap.DebugLevel)
	} else if strings.EqualFold("Info", logLevel) {
		atomicLevel.SetLevel(zap.InfoLevel)
	} else if strings.EqualFold("warn", logLevel) {
		atomicLevel.SetLevel(zap.WarnLevel)
	} else if strings.EqualFold("error", logLevel) {
		atomicLevel.SetLevel(zap.ErrorLevel)
	}
}
