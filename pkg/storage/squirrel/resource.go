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
	resourceT = resourceTableIdentifiers{
		tableName:   "app.resource",
		id:          "id",
		createdAt:   "created_at",
		updatedAt:   "updated_at",
		description: "description",
		reference:   "reference",
		level:       "level",
		typ:         "type",
	}

	resourceTagsT = resourcesTagsTableIdentifiers{
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
	query := r.builder.Select(fmt.Sprintf("r.*, array_agg(t.%s) as tags", tagT.name)).
		From(resourceT.tableName + " as r").
		Join(fmt.Sprintf("%s as rt on r.%s = rt.%s", resourceTagsT.tableName, resourceT.id, resourceTagsT.resourceId)).
		Join(fmt.Sprintf("%s as t on t.%s = rt.%s", tagT.tableName, tagT.id, resourceTagsT.tagId)).
		GroupBy("r." + resourceT.id).
		RunWith(r.conn)

	if filter.Level != "" {
		query = query.Where(sq.Eq{resourceT.level: filter.Level})
	}

	if filter.Type != "" {
		query = query.Where(sq.Eq{resourceT.typ: filter.Type})
	}

	if len(filter.Tags) > 0 {
		query = query.Where(sq.Eq{fmt.Sprintf("t.%s", tagT.name): filter.Tags})
	}

	query = applyLimitOffset(query, filter.LimitOffset)

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
			Insert(tagT.tableName).
			Columns(tagT.id, tagT.createdAt, tagT.updatedAt, tagT.name)

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
		tagsQuery := r.builder.Select("*").From(tagT.tableName).Where(sq.Eq{
			tagT.name: resource.Tags,
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
			Insert(resourceT.tableName).
			Columns(resourceT.id, resourceT.createdAt, resourceT.updatedAt, resourceT.description, resourceT.reference, resourceT.level, resourceT.typ).
			Values(resource.Id, resource.CreatedAt, resource.UpdatedAt, resource.Description, resource.Reference, resource.Level, resource.Type).
			RunWith(tx).
			ExecContext(ctx)

		if err != nil {
			return err
		}

		// Insert the (resource,tag) mapping
		mappingBuilder := r.builder.
			Insert(resourceTagsT.tableName).
			Columns(resourceTagsT.id, resourceTagsT.createdAt, resourceTagsT.updatedAt, resourceTagsT.resourceId, resourceTagsT.tagId)
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

func (r *resource) Delete(ctx context.Context, id string) error {
	return WithTx(ctx, r.conn, r.logger, func(tx *sqlx.Tx) error {

		query := r.builder.Delete(resourceT.tableName).
			Where(sq.Eq{resourceT.id: id}).
			RunWith(tx)

		result, err := query.ExecContext(ctx)
		if err != nil {
			return err
		}

		deletedRows, _ := result.RowsAffected()
		r.logger.Debugf("deleted %d rows when deleting resource %s", deletedRows, id)

		return nil
	})
}

func applyLimitOffset(builder sq.SelectBuilder, limitOffset domain.LimitOffset) sq.SelectBuilder {
	if limitOffset.Offset != nil {
		builder = builder.Offset(uint64(*limitOffset.Offset))
	}

	if limitOffset.Limit != nil {
		builder = builder.Limit(uint64(*limitOffset.Limit))
	}

	return builder
}
