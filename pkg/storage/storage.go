package storage

import (
	"context"

	"github.com/magdyamr542/go-web-service-template/pkg/domain"
)

type Storage interface {
	Resource() ResourceStorage

	Close(context.Context) error
}

type GetResourcesFilter struct {
	Tags  []string
	Type  string
	Level string
}
type ResourceStorage interface {
	GetByFilter(context.Context, GetResourcesFilter) ([]domain.Resource, error)
	Create(context.Context, *domain.Resource) error
}
