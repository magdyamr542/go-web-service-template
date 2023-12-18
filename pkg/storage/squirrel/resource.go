package squirrel

import (
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"go.uber.org/zap"

	"github.com/magdyamr542/go-web-service-template/pkg/domain"
	"github.com/magdyamr542/go-web-service-template/pkg/helpers/slices"
	"github.com/magdyamr542/go-web-service-template/pkg/storage"
)

type resourceTableIdentifiers struct {
	tableName   string
	id          string
	createdAt   string
	updatedAt   string
	description string
	reference   string
	level       string
	typ         string
}

type resourcesTagsTableIdentifiers struct {
	tableName  string
	id         string
	createdAt  string
	updatedAt  string
	resourceId string
	tagId      string
}

var (
	rti = resourceTableIdentifiers{
		tableName:   "app.resource",
		id:          "id",
		createdAt:   "created_at",
		updatedAt:   "updated_at",
		description: "description",
		reference:   "reference",
		level:       "level",
		typ:         "type",
	}

	rtti = resourcesTagsTableIdentifiers{
		tableName:  "app.resources_tags",
		id:         "id",
		createdAt:  "created_at",
		updatedAt:  "updated_at",
		resourceId: "resource_id",
		tagId:      "tag_id",
	}
)

type resourceDbRow struct {
	domain.DefaultFields
	Description string         `db:"description"`
	Reference   string         `db:"reference"`
	Level       string         `db:"level"`
	Type        string         `db:"type"`
	Tags        pq.StringArray `db:"tags"`
}

func (r *resourceDbRow) ToDomain() domain.Resource {
	return domain.Resource{
		DefaultFields: r.DefaultFields,
		Description:   r.Description,
		Reference:     r.Reference,
		Level:         domain.ResourceLevel(r.Level),
		Type:          domain.ResourceType(r.Type),
		Tags:          r.Tags,
	}
}

type resource struct {
	conn    *sqlx.DB
	builder sq.StatementBuilderType
	logger  *zap.SugaredLogger
}

func newResource(conn *sqlx.DB, builder sq.StatementBuilderType, logger *zap.Logger) *resource {
	return &resource{
		conn:    conn,
		builder: builder,
		logger:  logger.Sugar().Named("resource_storage"),
	}
}

func (r *resource) GetByFilter(ctx context.Context, filter storage.GetResourcesFilter) ([]domain.Resource, error) {
	query := r.builder.Select(fmt.Sprintf("r.*, array_agg(t.%s) as tags", tti.name)).
		From(rti.tableName + " as r").
		Join(fmt.Sprintf("%s as rt on r.%s = rt.%s", rtti.tableName, rti.id, rtti.resourceId)).
		Join(fmt.Sprintf("%s as t on t.%s = rt.%s", tti.tableName, tti.id, rtti.tagId)).
		GroupBy("r." + rti.id).
		RunWith(r.conn)

	if filter.Level != "" {
		query = query.Where(sq.Eq{rti.level: filter.Level})
	}

	if filter.Type != "" {
		query = query.Where(sq.Eq{rti.typ: filter.Type})
	}

	if len(filter.Tags) > 0 {
		query = query.Where(sq.Eq{fmt.Sprintf("t.%s", tti.name): filter.Tags})
	}

	rows, err := query.QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []resourceDbRow
	if err := sqlx.StructScan(rows, &result); err != nil {
		return nil, err
	}

	return slices.MapSlice(result, func(rdr resourceDbRow) domain.Resource { return rdr.ToDomain() }), nil

}

func (r *resource) Create(ctx context.Context, resource *domain.Resource) error {
	return WithTx(ctx, r.conn, r.logger, func(tx *sqlx.Tx) error {

		// Insert the tags
		tagInsertBuilder := r.builder.
			Insert(tti.tableName).
			Columns(tti.id, tti.createdAt, tti.updatedAt, tti.name)

		for _, tag := range resource.Tags {
			tagInsertBuilder = tagInsertBuilder.Values(uuid.NewString(), time.Now(), time.Now(), tag)
		}

		r.logSql(tagInsertBuilder.ToSql())

		_, err := tagInsertBuilder.
			Suffix("ON CONFLICT DO NOTHING").
			RunWith(tx).
			ExecContext(ctx)
		if err != nil {
			return err
		}

		// Get those tags
		tagsQuery := r.builder.Select("*").From(tti.tableName).Where(sq.Eq{
			tti.name: resource.Tags,
		}).RunWith(tx)

		r.logSql(tagsQuery.ToSql())

		rows, err := tagsQuery.QueryContext(ctx)
		if err != nil {
			return err
		}
		defer rows.Close()

		var tags []domain.Tag
		if err := sqlx.StructScan(rows, &tags); err != nil {
			return err
		}
		r.logger.Debugw("tags from db", "tags", tags)

		// Insert the resource
		_, err = r.builder.
			Insert(rti.tableName).
			Columns(rti.id, rti.createdAt, rti.updatedAt, rti.description, rti.reference, rti.level, rti.typ).
			Values(resource.Id, resource.CreatedAt, resource.UpdatedAt, resource.Description, resource.Reference, resource.Level, resource.Type).
			RunWith(tx).
			ExecContext(ctx)

		if err != nil {
			return err
		}

		// Insert the (resource,tag) mapping
		mappingBuilder := r.builder.
			Insert(rtti.tableName).
			Columns(rtti.id, rtti.createdAt, rtti.updatedAt, rtti.resourceId, rtti.tagId)
		for _, tag := range tags {
			mappingBuilder = mappingBuilder.Values(uuid.NewString(), time.Now(), time.Now(), resource.Id, tag.Id)
		}
		r.logSql(mappingBuilder.ToSql())

		_, err = mappingBuilder.
			Suffix("ON CONFLICT DO NOTHING").
			RunWith(tx).
			ExecContext(ctx)
		if err != nil {
			return err
		}

		return nil
	})
}

func (r *resource) logSql(sql string, args []interface{}, err error) {
	if err != nil {
		r.logger.With("err", err).Errorf("error generating sql")
		return
	}
	r.logger.Debugw("", "sql", sql, "args", args)

}
