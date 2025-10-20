package bootstrap

import (
	"github.com/prankevich/Auth_service/internal/config"
	"github.com/prankevich/Auth_service/internal/usecase"
	"net/http"

	http2 "github.com/prankevich/Auth_service/internal/adapter/driving/http"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

func initDB(cfg config.Postgres) (*sqlx.DB, error) {
	connConfig, err := pgx.ParseConfig(cfg.ConnectionURL())
	if err != nil {
		return nil, err
	}

	connStr := stdlib.RegisterConnConfig(connConfig)

	db, err := sqlx.Connect("pgx", connStr)
	if err != nil {
		return db, err
	}

	db.SetMaxOpenConns(cfg.MaxOpenConnections)
	db.SetMaxIdleConns(cfg.MaxIdleConnections)
	db.SetConnMaxLifetime(cfg.ConnectionMaxLifetime)

	return db, nil
}

func initHTTPService(
	cfg *config.Config,
	uc *usecase.UseCases,
) *http.Server {
	return http2.New(
		cfg,
		uc,
	)
}
