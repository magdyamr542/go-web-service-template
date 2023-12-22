package usecase

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"github.com/magdyamr542/go-web-service-template/pkg/storage"
)

type DeleteResource struct {
	store  storage.ResourceStorage
	logger *zap.Logger
}

func NewDeleteResource(store storage.ResourceStorage, logger *zap.Logger) *DeleteResource {
	return &DeleteResource{store: store, logger: logger}
}

func (g *DeleteResource) DeleteResource(ctx context.Context, id string) error {
	if err := g.store.Delete(ctx, id); err != nil {
		g.logger.With(zap.Error(err)).Error("failed to delete resource", zap.String("resource_id", id))
		return fmt.Errorf("error deleting resource")
	}
	return nil
}
