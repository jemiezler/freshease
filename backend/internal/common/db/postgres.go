package db

import (
	"context"
	"time"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"

	"freshease/backend/ent"
)

// dsn example: postgres://user:pass@localhost:5432/app?sslmode=disable
func NewEntClientPGX(ctx context.Context, dsn string, debug bool) (*ent.Client, func(context.Context) error, error) {
	cfg, err := pgx.ParseConfig(dsn)
	if err != nil {
		return nil, nil, err
	}
	sqlDB := stdlib.OpenDB(*cfg)

	// Pool tuning
	sqlDB.SetMaxOpenConns(50)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxIdleTime(5 * time.Minute)
	sqlDB.SetConnMaxLifetime(60 * time.Minute)

	// Verify connectivity
	if err := sqlDB.PingContext(ctx); err != nil {
		_ = sqlDB.Close()
		return nil, nil, err
	}

	drv := entsql.OpenDB(dialect.Postgres, sqlDB)
	client := ent.NewClient(ent.Driver(drv))
	if debug {
		client = client.Debug()
	}

	closeFn := func(c context.Context) error { return sqlDB.Close() }
	return client, closeFn, nil
}
