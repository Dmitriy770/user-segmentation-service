package sqlite

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/Dmitriy770/user-segmentation-service/internal/storage"
	"github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlute.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS segments(
		slug VARCHAR(50) PRIMARY KEY
	);
	`)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	stmt, err = db.Prepare(`
	CREATE TABLE IF NOT EXISTS segments_users(
		user_id INTEGER NOT NULL,
		segment_slug VARCHAR(50) NOT NULL,

		PRIMARY KEY(user_id, segment_slug),

		FOREIGN KEY (segment_slug) REFERENCES segments(slug) ON DELETE CASCADE
	);
	`)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) CreateSegment(slug string) error {
	const op = "storage.sqlite.DeleteURL"

	stmt, err := s.db.Prepare("INSERT INTO segments(slug) VALUES (?)")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec(slug)
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintPrimaryKey {
			return fmt.Errorf("%s: %w", op, storage.ErrSegmentExists)
		}

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) DeleteSegment(slug string) error {
	const op = "storage.sqlite.DeleteSegment"

	stmt, err := s.db.Prepare("DELETE FROM segments WHERE slug=?")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec(slug)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) AddUserToSegment(userId int, segmentSlug string) error {
	const op = "storage.sqlite.AddUserToSegment"

	stmt, err := s.db.Prepare("SELECT slug FROM segments WHERE slug=?")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	var dbSlug string
	err = stmt.QueryRow(segmentSlug).Scan(&dbSlug)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return storage.ErrSegmetNotFound
		}

		return fmt.Errorf("%s: %w", op, err)
	}

	stmt, err = s.db.Prepare("INSERT INTO segments_users (user_id, segment_slug) VALUES (?, ?)")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec(userId, segmentSlug)
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintPrimaryKey {
			return fmt.Errorf("%s: %w", op, storage.ErrUserHaveThisSegment)
		}

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) DeleteUserFromSegment(userId int, segmentSlug string) error {
	const op = "storage.sqlite.DeleteUserFromSegment"

	stmt, err := s.db.Prepare("DELETE FROM segments_users WHERE user_id=? AND segment_slug=?")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec(userId, segmentSlug)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) GetUserSegments(userId int) ([]string, error) {
	const op = "storage.sqlite.GetUserSegments"

	stmt, err := s.db.Prepare("SELECT segment_slug FROM segments_users WHERE user_id=?")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec(userId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	rows, err := stmt.Query(userId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var segments []string
	for rows.Next() {
		segment := ""
		err = rows.Scan(&segment)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		segments = append(segments, segment)
	}

	return segments, nil
}
