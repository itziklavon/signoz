package logger

import (
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"testing"
)

func TestGetBrandsConnectors(t *testing.T) {
	SetLevel("debug")
	level, _ := logLevels.Get("ROOT")
	assert.Equal(t, level.(zap.AtomicLevel).Level(), zap.DebugLevel)
	SetLevel("info")
	level, _ = logLevels.Get("ROOT")
	assert.Equal(t, level.(zap.AtomicLevel).Level(), zap.InfoLevel)
	SetLevel("warn")
	level, _ = logLevels.Get("ROOT")
	assert.Equal(t, level.(zap.AtomicLevel).Level(), zap.WarnLevel)
	SetLevel("error")
	level, _ = logLevels.Get("ROOT")
	assert.Equal(t, level.(zap.AtomicLevel).Level(), zap.ErrorLevel)
}

func TestNewLoggerWithName(t *testing.T) {
	NewLoggerWithName("test", zap.String("test", "test"))
	SetLevel("debug", "test")
	level, _ := logLevels.Get("test")
	assert.Equal(t, level.(zap.AtomicLevel).Level(), zap.DebugLevel)
	SetLevel("info", "test")
	level, _ = logLevels.Get("test")
	assert.Equal(t, level.(zap.AtomicLevel).Level(), zap.InfoLevel)
	SetLevel("warn", "test")
	level, _ = logLevels.Get("test")
	assert.Equal(t, level.(zap.AtomicLevel).Level(), zap.WarnLevel)
	SetLevel("error", "test")
	level, _ = logLevels.Get("test")
	assert.Equal(t, level.(zap.AtomicLevel).Level(), zap.ErrorLevel)
}
