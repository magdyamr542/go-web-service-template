package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/brpaz/echozap"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"

	"github.com/magdyamr542/go-web-service-template/pkg/api"
	"github.com/magdyamr542/go-web-service-template/pkg/handler"
	"github.com/magdyamr542/go-web-service-template/pkg/storage/squirrel"
	"github.com/magdyamr542/go-web-service-template/pkg/usecase"
)

func getLogger(env string) (*zap.Logger, error) {
	switch env {
	case "development":
		return zap.NewDevelopment()
	case "production":
		return zap.NewProduction()
	}
	return nil, nil
}

func main() {
	port := flag.String("port", "8000", "port to listen to")
	environment := flag.String("environment", "development", "current environment. values:(development,production)")
	flag.Parse()

	ctx := context.Background()

	// Setup echo
	e := echo.New()

	logger, err := getLogger(*environment)
	defer logger.Sync()
	if err != nil {
		log.Fatal(err)
	}

	e.Use(echozap.ZapLogger(logger))
	e.Use(middleware.Recover())

	// Setup the storage
	store, err := squirrel.NewDb(ctx, squirrel.Config{
		Host:     os.Getenv("DATABASE_HOST"),
		Port:     os.Getenv("DATABASE_PORT"),
		User:     os.Getenv("DATABASE_USER"),
		Password: os.Getenv("DATABASE_PASSWORD"),
		Database: os.Getenv("DATABASE_NAME"),
	}, logger)
	if err != nil {
		log.Fatal(err)
	}
	defer store.Close(ctx)

	// Setup the usecases
	getResourcesUsecase := usecase.NewGetResources(store.Resource(), logger)
	createResourceUsecase := usecase.NewCreateResource(store.Resource(), logger)

	// Setup the handler
	handler := handler.New(*getResourcesUsecase, *createResourceUsecase, logger)
	api.RegisterHandlers(e, handler)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", *port)))
}
