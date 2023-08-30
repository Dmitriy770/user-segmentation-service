package users

import "github.com/pkg/errors"

var (
	ErrUserDoesntHaveSegment = errors.New("user doesn`t have segment")
	ErrUserHaveSegment       = errors.New("user have segment")
	ErrSegmentDoesNotExist   = errors.New("segment does not exist")
)
