package squirrel

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/zap"

	"github.com/magdyamr542/go-web-service-template/pkg/storage"
)

type db struct {
	conn     *sqlx.DB
	resource *resource
}

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

func NewDb(ctx context.Context, cfg Config, logger *zap.Logger) (storage.Storage, error) {
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable", cfg.User, cfg.Password, cfg.Database, cfg.Host, cfg.Port)
	conn, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		return nil, err
	}

	stmtCache := sq.NewStmtCache(conn)
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).RunWith(stmtCache)

	return &db{
		conn:     conn,
		resource: newResource(conn, builder, logger),
	}, nil
}

func (d *db) Resource() storage.ResourceStorage {
	return d.resource
}

func (d *db) Close(ctx context.Context) error {
	if d.conn != nil {
		return d.conn.Close()
	}
	return nil
}
