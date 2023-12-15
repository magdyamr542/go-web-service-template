package usecase

import (
	"context"
	"fmt"

	"github.com/magdyamr542/go-web-service-template/pkg/domain"
	"github.com/magdyamr542/go-web-service-template/pkg/storage"
	"go.uber.org/zap"
)

type GetResources struct {
	store  storage.ResourcesStorage
	logger *zap.Logger
}

func NewGetResources(store storage.ResourcesStorage, logger *zap.Logger) *GetResources {
	return &GetResources{store: store, logger: logger}
}

type GetResourcesOptions struct {
	Tags  []string
	Type  string
	Level string
}

func (g *GetResources) GetResources(ctx context.Context, options GetResourcesOptions) ([]domain.Resource, error) {
	return nil, fmt.Errorf("not implemented")
}
