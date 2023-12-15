package dto

import (
	"github.com/magdyamr542/go-web-service-template/pkg/api"
	"github.com/magdyamr542/go-web-service-template/pkg/domain"
	"github.com/magdyamr542/go-web-service-template/pkg/helpers/slices"
	"github.com/oapi-codegen/runtime/types"

	"github.com/google/uuid"
)

func ResourceToDTO(r domain.Resource) api.Resource {
	return api.Resource{
		Id:          uuid.MustParse(r.Id),
		CreatedAt:   types.Date{Time: r.CreatedAt},
		UpdatedAt:   types.Date{Time: r.UpdatedAt},
		Description: r.Description,
		Level:       api.ResourceLevel(r.Level),
		Reference:   r.Reference,
		Tags:        slices.MapSlice(r.Tags, func(t domain.Tag) string { return t.Name }),
		Type:        api.ResourceType(r.Type),
	}
}
