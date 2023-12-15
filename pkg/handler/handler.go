package handler

import (
	"github.com/magdyamr542/go-web-service-template/pkg/usecase"
	"go.uber.org/zap"
)

type Handler struct {
	resourcesHandler
	versionHandler
}

func New(getResourcesUsecase usecase.GetResources, createResourceUsecase usecase.CreateResource, logger *zap.Logger) *Handler {
	return &Handler{
		resourcesHandler: *newResourcesHandler(getResourcesUsecase, createResourceUsecase, logger),
		versionHandler:   *newVersionHandler(),
	}
}
