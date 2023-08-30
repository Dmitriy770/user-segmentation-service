package segments

import (
	"log/slog"

	"github.com/Dmitriy770/user-segmentation-service/internal/entities"
	"github.com/Dmitriy770/user-segmentation-service/internal/lib/logger/sl"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/pkg/errors"
)

type repository struct {
	log *slog.Logger
	db  *sqlx.DB
}

func NewRepository(log *slog.Logger, db *sqlx.DB) *repository {
	return &repository{
		log: log,
		db:  db,
	}
}

func (r *repository) AddSegment(slug string) error {
	const op = "segments.repository.AddSegment"

	_, err := r.db.Exec(
		`
		 INSERT INTO segments (slug)
		 VALUES ($1)
		`,
		slug,
	)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				return ErrSlugBusy
			default:
				return errors.Wrap(err, "insert slug")
			}

		}
		r.log.Error(err.Error(), slog.String("op", op))
		return errors.Wrap(err, "insert slug")
	}

	return nil
}

func (r *repository) DeleteSegment(slug string) error {
	const op = "segments.repository.DeleteSegment"

	_, err := r.db.Exec(
		`
		 DELETE FROM segments
		 WHERE slug=$1
		`,
		slug,
	)
	if err != nil {
		r.log.Error("failed to delte segment", slog.String("op", op), sl.Err(err))
		return errors.Wrap(err, "delete slug")
	}

	return nil
}

func (r *repository) GetSegment(slug string) (*string, error) {
	const op = "segments.repository.GetSegment"

	rawSegment := make([]entities.Segment, 0)
	err := r.db.Select(
		&rawSegment,
		`
		SELECT slug
		FROM segments
		WHERE slug=$1;
		`,
		slug,
	)
	if err != nil {
		r.log.Error("failed to get segment", slog.String("op", op), sl.Err(err))
		return nil, err
	}

	if len(rawSegment) > 0 {
		return &rawSegment[0].Slug, nil
	} else {
		return nil, nil
	}
}

func (r *repository) GetSegmentUsers(slug string) (*[]int, error) {
	const op = "segments.repository.GetSegmentUsers"

	rawUserWithSegment := make([]entities.UserSegment, 0)
	err := r.db.Select(
		&rawUserWithSegment,
		`
		SELECT user_id, segment_slug
		FROM segments_users
		WHERE segment_slug=$1;
		`,
		slug,
	)
	if err != nil {
		r.log.Error("failed to get users with segment", slog.String("op", op), sl.Err(err))
		return nil, errors.Wrap(err, "get segment users")
	}

	users := make([]int, 0)
	for _, userSegment := range rawUserWithSegment {
		users = append(users, userSegment.UserId)
	}

	return &users, nil
}

func (r *repository) DeleteSegmentUsers(slug string) error {
	const op = "segments.repository.DeleteSegmentUsers"

	_, err := r.db.Exec(
		`
				DELETE FROM segments_users
				WHERE segment_slug=$1	
			`,
		slug,
	)
	if err != nil {
		r.log.Error("failed to dekete segment from users", slog.String("op", op), sl.Err(err))
		return errors.Wrap(err, "delete segment users")
	}

	return nil
}
