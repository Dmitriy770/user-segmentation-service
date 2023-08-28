package storage

import "errors"

var (
	ErrSegmetNotFound      = errors.New("segment not found")
	ErrSegmentExists       = errors.New("segment exists")
	ErrUserHaveThisSegment = errors.New("user have this segment")
)
