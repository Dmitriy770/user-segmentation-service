package segments

import (
	"log/slog"

	"github.com/Dmitriy770/user-segmentation-service/internal/lib/logger/sl"
)

type SegmentsRep interface {
	AddSegment(slug string) error
	DeleteSegment(slug string) error
	GetSegment(slug string) (*string, error)
	GetSegmentUsers(slug string) (*[]int, error)
	DeleteSegmentUsers(slug string) error
}

type service struct {
	log         *slog.Logger
	segmentsRep SegmentsRep
}

func NewService(log *slog.Logger, segmentsRep SegmentsRep) *service {
	return &service{
		log:         log,
		segmentsRep: segmentsRep,
	}
}

func (s *service) AddSegment(slug string) error {
	const op = "segments.service.AddSegment"

	err := s.segmentsRep.AddSegment(slug)
	if err != nil {
		s.log.Error("failed to add segment", slog.String("op", op), sl.Err(err))
		return err
	}
	return nil
}

func (s *service) DeleteSegment(slug string) error {
	const op = "segments.service.DeleteSegment"

	dbSlug, err := s.segmentsRep.GetSegment(slug)
	if err != nil {
		s.log.Error("failed to delete segment", slog.String("op", op), sl.Err(err))
		return err
	}
	if dbSlug == nil {
		s.log.Error("failed to delete segment", slog.String("op", op), sl.Err(err))
		return ErrSlugNotFound
	}

	err = s.segmentsRep.DeleteSegmentUsers(slug)
	if err != nil {
		s.log.Error("failed to delete segment", slog.String("op", op), sl.Err(err))
		return err
	}

	err = s.segmentsRep.DeleteSegment(slug)
	if err != nil {
		s.log.Error("failed to delete segment", slog.String("op", op), sl.Err(err))
		return err
	}

	return nil
}
