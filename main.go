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

	"github.com/brpaz/echozap"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"

	"github.com/magdyamr542/go-web-service-template/pkg/api"
	"github.com/magdyamr542/go-web-service-template/pkg/handler"
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
	flag.Parse()

	ctx := context.Background()

	// Setup echo.
	e := echo.New()

	logger, err := getLogger(*environment)
	defer logger.Sync()
	if err != nil {
		log.Print(err)
		return 1
	}

	e.Use(echozap.ZapLogger(logger))
	e.Use(middleware.Recover())

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

	// Setup the handlers.
	handler := handler.New(*getResourcesUsecase, *createResourceUsecase, logger)
	api.RegisterHandlers(e, handler)

	// Start server.
	serverErrCh := make(chan error, 1)
	go func() {
		if err := e.Start(fmt.Sprintf(":%s", *port)); err != nil && err != http.ErrServerClosed {
			logger.With(zap.Error(err)).Error("error shutting down the server")
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
			logger.With(zap.Error(err)).Error("can't shut down the server")
		}
	}
	logger.Info("Shutting down db...")
	if err := store.Close(ctx); err != nil {
		logger.With(zap.Error(err)).Error("can't close db connection")
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
