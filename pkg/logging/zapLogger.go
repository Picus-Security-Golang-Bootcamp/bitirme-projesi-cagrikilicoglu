package logging

import (
	"fmt"

	"github.com/cagrikilicoglu/shopping-basket/pkg/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewLogger creates a new logger using zap
func NewZapLogger(config *config.Config) {
	logLevel, err := zapcore.ParseLevel(config.Logger.Level)
	if err != nil {
		panic(fmt.Sprintf("Unknown log level: %v", logLevel))
	}
	var cfg zap.Config
	if config.Logger.Development {
		cfg = zap.NewDevelopmentConfig()
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	} else {
		cfg = zap.NewProductionConfig()
	}
	logger, err := cfg.Build()
	if err != nil {
		logger = zap.NewNop()
	}
	zap.ReplaceGlobals(logger)
}

func Close() {
	defer zap.L().Sync()
}
