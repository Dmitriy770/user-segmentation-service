package postgres

import (
	"fmt"

	"github.com/Dmitriy770/user-segmentation-service/internal/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

func New(cfg config.PostgreSQL) (*sqlx.DB, error) {
	dataSource := fmt.Sprintf(
		"user=%s password=%s host=%s port=%d dbname=%s sslmode=disable",
		cfg.User, cfg.Password, cfg.Address, cfg.Port, cfg.DB,
	)

	conn, err := sqlx.Connect("postgres", dataSource)
	if err != nil {
		return nil, errors.Wrap(err, "sqlx connect")
	}

	err = conn.Ping()
	if err != nil {
		return nil, errors.Wrap(err, "ping failed")
	}

	return conn, nil
}
