package handler

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/magdyamr542/go-web-service-template/pkg/api"
	"github.com/magdyamr542/go-web-service-template/pkg/domain"
	"github.com/magdyamr542/go-web-service-template/pkg/dto"
	"github.com/magdyamr542/go-web-service-template/pkg/helpers/pointers"
	"github.com/magdyamr542/go-web-service-template/pkg/usecase"
	"github.com/magdyamr542/go-web-service-template/pkg/validation"
	"go.uber.org/zap"
)

var (
	resourceTypes = []api.ResourceType{
		api.ARTICLE,
		api.VIDEO,
	}

	resourceLevels = []api.ResourceLevel{
		api.BEGINNER,
		api.INTERMEDIATE,
		api.ADVANCED,
	}
)

type resourcesHandler struct {
	getResourcesUsecase   usecase.GetResources
	createResourceUsecase usecase.CreateResource
	logger                *zap.SugaredLogger
}

func newResourcesHandler(getResourcesUsecase usecase.GetResources, createResourceUsecase usecase.CreateResource, logger *zap.Logger) *resourcesHandler {
	return &resourcesHandler{
		getResourcesUsecase:   getResourcesUsecase,
		createResourceUsecase: createResourceUsecase,
		logger:                logger.Sugar().Named("resouces_handler"),
	}
}

func (h *resourcesHandler) GetResources(ctx echo.Context, params api.GetResourcesParams) error {
	h.logger.Debugw("handling get resources with", "tags", params.Tags, "level", params.Level, "type", params.Type)

	resources, err := h.getResourcesUsecase.GetResources(ctx.Request().Context(), usecase.GetResourcesOptions{
		Tags:  params.Tags,
		Type:  pointers.DefaultIfNil((*string)(params.Type)),
		Level: pointers.DefaultIfNil((*string)(params.Level)),
	})
	if err != nil {
		h.logger.Errorw("error getting resources", "err", err)
		return fmt.Errorf("error getting the resources: %v", err)
	}

	return ctx.JSON(http.StatusOK, resources)
}

func (h *resourcesHandler) CreateResource(ctx echo.Context) error {
	var body api.NewResource
	if err := ctx.Bind(&body); err != nil {
		return err
	}

	h.logger.Debugw("handling create resource", "body", body)

	// Validate the request.
	if err := h.validateNewResource(body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Create the resource.
	resource, err := h.createResourceUsecase.CreateResource(ctx.Request().Context(), usecase.CreateResourceRequest{
		Description: body.Description,
		Reference:   body.Reference,
		Tags:        body.Tags,
		Type:        domain.ResourceType(body.Type),
		Level:       domain.ResourceLevel(body.Level),
	})
	if err != nil {
		h.logger.Errorw("error creating resource", "err", err)
		return echo.ErrInternalServerError
	}

	return ctx.JSON(http.StatusCreated, dto.ResourceToDTO(*resource))
}

func (h *resourcesHandler) validateNewResource(body api.NewResource) error {

	funcs := []validation.ValidationFunc{
		func() error {
			return validation.RequiredStrField("description", body.Description)
		},
		func() error {
			return validation.RequiredStrField("reference", body.Reference)
		},
		func() error {
			return validation.RequiredStrField("reference", body.Reference)
		},
		func() error {
			return validation.MinLenField("tags", body.Tags, 1)
		},
		func() error {
			return validation.ValidItemsField("tags", body.Tags, validation.RequiredStr)
		},
		func() error {
			return validation.OneOfField("level", body.Level, resourceLevels)
		},
		func() error {
			return validation.OneOfField("type", body.Type, resourceTypes)
		},
	}

	return validation.Validate(funcs...)
}
