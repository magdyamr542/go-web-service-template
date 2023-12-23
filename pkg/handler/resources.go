package handler

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
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
	deleteResourceUsecase usecase.DeleteResource
	logger                *zap.SugaredLogger
}

func newResourcesHandler(
	getResourcesUsecase usecase.GetResources,
	createResourceUsecase usecase.CreateResource,
	deleteResourceUsecase usecase.DeleteResource,
	logger *zap.Logger) *resourcesHandler {
	return &resourcesHandler{
		getResourcesUsecase:   getResourcesUsecase,
		createResourceUsecase: createResourceUsecase,
		deleteResourceUsecase: deleteResourceUsecase,
		logger:                logger.Sugar().Named("resouces_handler"),
	}
}

func (h *resourcesHandler) GetResources(ctx echo.Context, params api.GetResourcesParams) error {
	h.logger.Debugw("handling get resources with", "params", params)

	if err := h.validateGetResources(params); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	resources, err := h.getResourcesUsecase.GetResources(ctx.Request().Context(), usecase.GetResourcesOptions{
		Tags:  strings.Split(params.Tags, ","),
		Type:  pointers.DefaultIfNil((*string)(params.Type)),
		Level: pointers.DefaultIfNil((*string)(params.Level)),
		LimitOffset: domain.LimitOffset{
			Limit:  params.Limit,
			Offset: params.Offset,
		},
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
		func() error { return validation.Fielded("description", validation.RequiredStr(body.Description)) },
		func() error { return validation.Fielded("reference", validation.RequiredStr(body.Reference)) },
		func() error { return validation.Fielded("reference", validation.RequiredStr(body.Reference)) },
		func() error { return validation.Fielded("tags", validation.MinLen(body.Tags, 1)) },
		func() error { return validation.Fielded("level", validation.OneOf(body.Level, resourceLevels)) },
		func() error { return validation.Fielded("type", validation.OneOf(body.Type, resourceTypes)) },
		func() error {
			return validation.Fielded("tags", validation.ValidItems(body.Tags, validation.RequiredStr))
		},
	}
	return validation.Validate(funcs...)
}

func (h *resourcesHandler) validateGetResources(params api.GetResourcesParams) error {
	funcs := []validation.ValidationFunc{func() error { return validation.Fielded("tags", validation.RequiredStr(params.Tags)) }}
	if params.Level != nil {
		funcs = append(funcs, func() error { return validation.Fielded("level", validation.OneOf(*params.Level, resourceLevels)) })
	}
	if params.Type != nil {
		funcs = append(funcs, func() error { return validation.Fielded("type", validation.OneOf(*params.Type, resourceTypes)) })
	}
	if params.Offset != nil {
		funcs = append(funcs, func() error { return validation.Fielded("offset", validation.Min(*params.Offset, 0, true)) })
	}
	if params.Limit != nil {
		funcs = append(funcs, func() error { return validation.Fielded("limit", validation.Max(*params.Limit, domain.MaxLimit, true)) })
	}
	return validation.Validate(funcs...)
}

func (h *resourcesHandler) DeleteResource(ctx echo.Context, id uuid.UUID) error {
	idStr := id.String()
	h.logger.Debugw("handling delete resource", "resource_id", idStr)

	err := h.deleteResourceUsecase.DeleteResource(ctx.Request().Context(), idStr)
	if err != nil {
		h.logger.With(zap.Error(err)).Errorw("error deleting resource", "resource_id", idStr)
		return fmt.Errorf("error deleting the resource")
	}

	return ctx.NoContent(http.StatusNoContent)
}
