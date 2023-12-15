package usecase

import (
	"context"
	"fmt"

	"github.com/magdyamr542/go-web-service-template/pkg/domain"
	"github.com/magdyamr542/go-web-service-template/pkg/helpers/slices"
	"github.com/magdyamr542/go-web-service-template/pkg/storage"
	"go.uber.org/zap"
)

type CreateResource struct {
	store  storage.ResourcesStorage
	logger *zap.SugaredLogger
}

func NewCreateResource(store storage.ResourcesStorage, logger *zap.Logger) *CreateResource {
	return &CreateResource{store: store, logger: logger.Sugar().Named("create_resource_usecase")}
}

type CreateResourceRequest struct {
	Description string
	Level       domain.ResourceLevel
	Reference   string
	Tags        []string
	Type        domain.ResourceType
}

func (g *CreateResource) CreateResource(ctx context.Context, request CreateResourceRequest) (*domain.Resource, error) {
	tags := slices.MapSlice[string, domain.Tag](request.Tags, func(s string) domain.Tag {
		return domain.Tag{
			DefaultFields: domain.NewDefaultFields(),
			Name:          s,
		}
	})
	resource, err := g.store.Create(&domain.Resource{
		DefaultFields: domain.NewDefaultFields(),
		Description:   request.Description,
		Reference:     request.Reference,
		Level:         request.Level,
		Type:          request.Type,
		Tags:          tags,
	})
	if err != nil {
		g.logger.Errorw("error creating new resource", "err", err, "request", request)
		return nil, fmt.Errorf("error creating resource")
	}
	return resource, nil
}
