package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/brpaz/echozap"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/magdyamr542/go-web-service-template/pkg/api"
	"github.com/magdyamr542/go-web-service-template/pkg/handler"
	"github.com/magdyamr542/go-web-service-template/pkg/usecase"
	"go.uber.org/zap"
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

	// Setup the usecases
	getResourcesUsecase := usecase.NewGetResources(nil, logger)
	createResourceUsecase := usecase.NewCreateResource(nil, logger)

	// Setup the handler
	handler := handler.New(*getResourcesUsecase, *createResourceUsecase, logger)
	api.RegisterHandlers(e, handler)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", *port)))
}
