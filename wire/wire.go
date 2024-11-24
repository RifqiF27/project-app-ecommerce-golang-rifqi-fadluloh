//go:build wireinject
// +build wireinject

package wire

import (
	"ecommerce/database"
	"ecommerce/handler"
	"ecommerce/repository"
	"ecommerce/router"
	"ecommerce/service"
	"ecommerce/util"
	"os"
	"strconv"
	"log"

	"github.com/go-chi/chi/v5"
	"github.com/google/wire"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func ProvideConfiguration() util.Configuration {
	return util.ReadConfiguration()
}

var ConfigSet = wire.NewSet(ProvideConfiguration)

func ProvideLogger() (*zap.Logger, error) {

	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found. Using system environment variables.")
		return nil, err
	}

	debugStr := os.Getenv("DEBUG")
	debug, err := strconv.ParseBool(debugStr)
	if err != nil {
		debug = false
	}

	logLevel := zap.InfoLevel
	if debug {
		logLevel = zap.DebugLevel
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "T",
		LevelKey:       "L",
		MessageKey:     "M",
		CallerKey:      "C",
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	file, err := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Panicf("Failed to open log file: %v", err)
		return nil, err
	}

	core := zapcore.NewTee(
		zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			zapcore.AddSync(file),
			logLevel,
		),
		zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoderConfig),
			zapcore.AddSync(os.Stdout),
			logLevel,
		),
	)

	logger := zap.New(core)
	logger.Info("Logger initialized successfully")
	return logger, nil
}

func InitializeRouterHandler() (*chi.Mux, error) {
	wire.Build(
		ConfigSet,
		database.InitDB,
		ProvideLogger,

		repository.NewAuthRepository,
		service.NewAuthService,
		handler.NewAuthHandler,

		repository.NewHomePageRepository,
        service.NewHomePageService,
        handler.NewHomePageHandler,

		repository.NewCheckoutRepository,
		service.NewCheckoutService,
		handler.NewCheckoutHandler,
		
		router.NewRouter,
	)
	return nil, nil
}
