package usecase

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"

	"github.com/magdyamr542/go-web-service-template/pkg/domain"
	"github.com/magdyamr542/go-web-service-template/pkg/storage"
)

type GetResources struct {
	store      storage.ResourceStorage
	logger     *zap.Logger
	tagsMetric *prometheus.CounterVec
}

func NewGetResources(store storage.ResourceStorage, logger *zap.Logger, tagsMetric *prometheus.CounterVec) *GetResources {
	return &GetResources{store: store, logger: logger, tagsMetric: tagsMetric}
}

type GetResourcesOptions struct {
	Tags  []string
	Type  string
	Level string
}

func (g *GetResources) GetResources(ctx context.Context, options GetResourcesOptions) ([]domain.Resource, error) {
	if g.tagsMetric != nil {
		go func() {
			for _, t := range options.Tags {
				g.tagsMetric.With(prometheus.Labels{"tags": t}).Inc()
			}
		}()
	}

	return g.store.GetByFilter(ctx, storage.GetResourcesFilter{
		Level: options.Level,
		Tags:  options.Tags,
		Type:  options.Type,
	})
}
