package users

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/Dmitriy770/user-segmentation-service/internal/entities"
	"github.com/Dmitriy770/user-segmentation-service/internal/lib/logger/sl"
	"github.com/Dmitriy770/user-segmentation-service/internal/models"
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

func (r *repository) AddUser(user models.User) error {
	const op = "users.repository.AddUser"

	_, err := r.db.Exec(
		`
		 INSERT INTO users (id)
		 VALUES ($1);
		`,
		user.ID,
	)
	if err != nil {
		r.log.Error("failed to add user", slog.String("op", op), sl.Err(err))
		return errors.Wrap(err, "insert user")
	}

	return nil
}

func (r *repository) GetUser(user_id int) (*models.User, error) {
	const op = "users.repository.GetUser"

	rawUserWithSegment := make([]entities.UserSegment, 0)
	err := r.db.Select(
		&rawUserWithSegment,
		`
		SELECT user_id, segment_slug
		FROM segments_users
		WHERE user_id=$1;
		`,
		user_id,
	)
	if err != nil {
		r.log.Error("failed to get user", slog.String("op", op), sl.Err(err))
		return nil, errors.Wrap(err, "get user")
	}

	user := &models.User{
		ID:       user_id,
		Segments: make([]string, 0),
	}

	for _, segmentUser := range rawUserWithSegment {
		user.Segments = append(user.Segments, segmentUser.SegmentSlug)
	}

	return user, nil
}

func (r *repository) AddSegmentsToUser(user_id int, segments []string) error {
	const op = "users.repository.AddSegmentsToUser"

	values := make([]string, 0)
	for _, segment := range segments {
		value := fmt.Sprintf("(%d,'%s')", user_id, segment)
		values = append(values, value)
	}

	query := fmt.Sprintf(
		`
		INSERT INTO segments_users (user_id, segment_slug)
		VALUES %s;
		`,
		strings.Join(values, ", "),
	)

	_, err := r.db.Exec(query)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				return ErrUserHaveSegment
			case "foreign_key_violation":
				return ErrSegmentDoesNotExist
			default:
				return errors.Wrap(err, "insert slug")
			}
		}
		r.log.Error("failed to add segmets to user", slog.String("op", op), sl.Err(err))
		return errors.Wrap(err, "add segments to user")
	}

	return nil
}

func (r *repository) DeleteSegmentsFromUser(user_id int, segments []string) error {
	const op = "users.repository.DeleteSegmentsFromUser"

	values := make([]string, 0)
	for _, segment := range segments {
		value := fmt.Sprintf("'%s'", segment)
		values = append(values, value)
	}

	query := fmt.Sprintf(
		`
		DELETE FROM segments_users
		WHERE user_id=%d AND segment_slug IN (%s);
		`,
		user_id,
		strings.Join(values, ", "),
	)

	_, err := r.db.Exec(query)
	if err != nil {
		r.log.Error("failed to delete segmets from user", slog.String("op", op), sl.Err(err))
		return errors.Wrap(err, "add segments to user")
	}

	return nil
}
