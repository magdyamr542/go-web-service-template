package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"

	"github.com/magdyamr542/go-web-service-template/pkg/api"
	"github.com/magdyamr542/go-web-service-template/pkg/handler"
	"github.com/magdyamr542/go-web-service-template/pkg/logging"
	"github.com/magdyamr542/go-web-service-template/pkg/metrics"
	"github.com/magdyamr542/go-web-service-template/pkg/storage/squirrel"
	"github.com/magdyamr542/go-web-service-template/pkg/usecase"
)

func main() {
	if code := realMain(); code != 0 {
		os.Exit(code)
	}
}

func realMain() int {
	port := flag.String("port", "8000", "port to listen to")
	environment := flag.String("environment", "development", "current environment. values:(development,production)")
	enableMetrics := flag.Bool("with-metrics", false, "whether to enable exporting prometheus metrics")
	flag.Parse()

	ctx := context.Background()

	// Setup echo.
	e := echo.New()

	// Setup logger.
	logger, err := getLogger(*environment)
	defer logger.Sync()
	if err != nil {
		log.Print(err)
		return 1
	}

	// Setup metrics.
	var mtrcs metrics.Metrics
	if *enableMetrics {
		logger.Info("Setting up the metrics...")
		mtrcs = metrics.New()
	}

	// Setup middlewares.
	e.Use(logging.ZapLogger(logging.LoggerMigglewareConfig{
		Skip: []string{
			"/metrics", // Don't log prometheus scraping requests.
		},
	}, logger))
	e.Use(middleware.Recover())

	if *enableMetrics {
		logger.Info("Will server prometheus metrics on /metrics")
		e.Use(echoprometheus.NewMiddlewareWithConfig(echoprometheus.MiddlewareConfig{
			Registerer: mtrcs.Registry,
			Subsystem:  "app",
		}))
		e.GET("/metrics", echoprometheus.NewHandlerWithConfig(echoprometheus.HandlerConfig{
			Gatherer: mtrcs.Registry.(prometheus.Gatherer),
		}))
	}

	// Setup the storage.
	store, err := squirrel.NewDb(ctx, squirrel.Config{
		Host:     os.Getenv("DATABASE_HOST"),
		Port:     os.Getenv("DATABASE_PORT"),
		User:     os.Getenv("DATABASE_USER"),
		Password: os.Getenv("DATABASE_PASSWORD"),
		Database: os.Getenv("DATABASE_NAME"),
	}, logger)
	if err != nil {
		logger.With(zap.Error(err)).Error("can't create db")
		return 1
	}
	defer store.Close(ctx)

	// Setup the usecases.
	getResourcesUsecase := usecase.NewGetResources(store.Resource(), logger)
	createResourceUsecase := usecase.NewCreateResource(store.Resource(), logger)
	deleteResourceUsecase := usecase.NewDeleteResource(store.Resource(), logger)

	// Setup the handlers.
	handler := handler.New(*getResourcesUsecase, *createResourceUsecase, *deleteResourceUsecase, logger)
	api.RegisterHandlers(e, handler)

	// Start server.
	serverErrCh := make(chan error, 1)
	go func() {
		if err := e.Start(fmt.Sprintf(":%s", *port)); err != nil && err != http.ErrServerClosed {
			serverErrCh <- err
		}
	}()

	// Graceful shutdown.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	var serverErr error
	select {
	case <-quit:
	case serverErr = <-serverErrCh:
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if serverErr == nil {
		logger.Info("Shutting down server...")
		if err := e.Shutdown(ctx); err != nil {
			logger.With(zap.Error(err)).Error("error shutting down the server")
		}
	} else {
		logger.With(zap.Error(err)).Error("error shutting down the server")
	}

	logger.Info("Shutting down db...")
	if err := store.Close(ctx); err != nil {
		logger.With(zap.Error(err)).Error("error closing the db")
	}

	return 0
}

func getLogger(env string) (*zap.Logger, error) {
	switch env {
	case "development":
		return zap.NewDevelopment()
	case "production":
		return zap.NewProduction()
	}
	return nil, nil
}
