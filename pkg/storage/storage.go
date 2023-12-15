package storage

import "github.com/magdyamr542/go-web-service-template/pkg/domain"

type Storage interface {
	Resources() ResourcesStorage
	Tags() TagsStorage
}

type GetResourcesFilter struct {
	Tags  []string
	Type  string
	Level string
}
type ResourcesStorage interface {
	GetByFilter(filter GetResourcesFilter) ([]domain.Resource, error)
	Create(*domain.Resource) (*domain.Resource, error)
}

type TagsStorage interface {
	GetByNames(name []string) ([]domain.Tag, error)
	Insert(tag domain.Tag) error
}
