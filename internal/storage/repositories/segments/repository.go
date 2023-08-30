package segments

import (
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *repository {
	return &repository{
		db: db,
	}
}

func (r *repository) AddSegment(slug string) error {
	_, err := r.db.Exec(
		`
		 INSERT INTO segments (slug)
		 VALUES ($1)
		`,
		slug,
	)
	if err != nil {
		return errors.Wrap(err, "insert slug")
	}

	return nil
}

func (r *repository) DeleteSegment(slug string) error {
	_, err := r.db.Exec(
		`
		 DELETE FROM segments
		 WHERE slug=$1
		`,
		slug,
	)
	if err != nil {
		return errors.Wrap(err, "delete slug")
	}

	return nil
}
