package segments

import "errors"

var (
	ErrSlugBusy     = errors.New("segment with the same slug already exists")
	ErrSlugNotFound = errors.New("segment with same slug not found")
)
