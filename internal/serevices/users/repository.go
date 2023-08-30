package users

import (
	"fmt"
	"strings"

	"github.com/Dmitriy770/user-segmentation-service/internal/entities"
	"github.com/Dmitriy770/user-segmentation-service/internal/models"
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

func (r *repository) AddUser(user models.User) error {
	_, err := r.db.Exec(
		`
		 INSERT INTO users (id)
		 VALUES ($1);
		`,
		user.ID,
	)
	if err != nil {
		return errors.Wrap(err, "insert user")
	}

	return nil
}

func (r *repository) GetUser(user_id int) (*models.User, error) {

	rawUserWithSegment := make([]entities.UserSegment, 0)
	err := r.db.Select(
		&rawUserWithSegment,
		`
		SELECT user_id, slug_segment
		FROM segments_users
		WHERE user_id=$1;
		`,
		user_id,
	)
	if err != nil {
		return nil, errors.Wrap(err, "get user")
	}

	user := &models.User{ID: user_id}

	for _, segmentUser := range rawUserWithSegment {
		user.Segments = append(user.Segments, segmentUser.SegmentSlug)
	}

	return user, nil
}

func (r *repository) AddSegmentsToUser(user_id int, segments []string) error {
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
		return errors.Wrap(err, "add segments to user")
	}

	return nil
}

func (r *repository) DeleteSegmentsFromUser(user_id int, segments []string) error {
	values := make([]string, 0)
	for _, segment := range segments {
		value := fmt.Sprintf("'%s'", segment)
		values = append(values, value)
	}

	query := fmt.Sprintf(
		`
		DELETE FROM segments_users
		WHERE user_id=%d AND segment_slug IN  (%s);
		`,
		user_id,
		strings.Join(values, ", "),
	)

	_, err := r.db.Exec(query)
	if err != nil {
		return errors.Wrap(err, "add segments to user")
	}

	return nil
}
