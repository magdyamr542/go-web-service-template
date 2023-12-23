package usecase

import (
	"context"

	"go.uber.org/zap"

	"github.com/magdyamr542/go-web-service-template/pkg/domain"
	"github.com/magdyamr542/go-web-service-template/pkg/helpers/pointers"
	"github.com/magdyamr542/go-web-service-template/pkg/storage"
)

type GetResources struct {
	store  storage.ResourceStorage
	logger *zap.Logger
}

func NewGetResources(store storage.ResourceStorage, logger *zap.Logger) *GetResources {
	return &GetResources{store: store, logger: logger}
}

type GetResourcesOptions struct {
	domain.LimitOffset
	Tags  []string
	Type  string
	Level string
}

func (g *GetResources) GetResources(ctx context.Context, options GetResourcesOptions) ([]domain.Resource, error) {
	if options.Limit == nil {
		options.Limit = pointers.Ptr(domain.DefaultLimit)
	}
	if options.Offset == nil {
		options.Offset = pointers.Ptr(domain.DefaultOffset)
	}
	return g.store.GetByFilter(ctx, storage.GetResourcesFilter{
		Level:       options.Level,
		Tags:        options.Tags,
		Type:        options.Type,
		LimitOffset: options.LimitOffset,
	})
}
