package users

import (
	"log/slog"

	"github.com/Dmitriy770/user-segmentation-service/internal/lib/logger/sl"
	"github.com/Dmitriy770/user-segmentation-service/internal/models"
	"github.com/pkg/errors"
)

type UsersRep interface {
	AddUser(user models.User) error
	GetUser(user_id int) (*models.User, error)
	AddSegmentsToUser(user_id int, segments []string) error
	DeleteSegmentsFromUser(user_id int, segments []string) error
}

type service struct {
	log      *slog.Logger
	usersRep UsersRep
}

func NewService(log *slog.Logger, usersRep UsersRep) *service {
	return &service{
		log:      log,
		usersRep: usersRep,
	}
}

func (s *service) GetUser(id int) (*models.User, error) {
	const op = "users.service.GetUser"

	user, err := s.usersRep.GetUser(id)
	if err != nil {
		s.log.Error("failed to get user", slog.String("op", op), sl.Err(err))
		return nil, err
	}

	return user, err
}

func (s *service) UpdateUser(userId int, slugsForAdd []string, slugsForDelete []string) error {
	const op = "users.service.UpdateUser"

	s.usersRep.AddUser(models.User{ID: userId})

	if len(slugsForDelete) > 0 {
		user, err := s.usersRep.GetUser(userId)
		if err != nil {
			s.log.Error("failed to update user", slog.String("op", op), sl.Err(err))
			return errors.Wrap(err, "update users")
		}

		mapOldSlugs := make(map[string]struct{}, 0)
		for _, slug := range user.Segments {
			mapOldSlugs[slug] = struct{}{}
		}

		for _, slug := range slugsForDelete {
			if _, ok := mapOldSlugs[slug]; !ok {
				s.log.Error("failed to update user", slog.String("op", op), sl.Err(ErrUserDoesntHaveSegment))
				return ErrUserDoesntHaveSegment
			}
		}
	}

	if len(slugsForAdd) > 0 {
		err := s.usersRep.AddSegmentsToUser(userId, slugsForAdd)
		if err != nil {
			s.log.Error("failed to update user", slog.String("op", op), sl.Err(err))
			return errors.Wrap(err, "update users")
		}
	}

	if len(slugsForDelete) > 0 {
		err := s.usersRep.DeleteSegmentsFromUser(userId, slugsForDelete)
		if err != nil {
			s.log.Error("failed to update user", slog.String("op", op), sl.Err(err))
			return errors.Wrap(err, "update users")
		}
	}

	return nil
}
