package entities

type UserSegment struct {
	UserId      int    `db:"user_id"`
	SegmentSlug string `db:"segment_slug"`
}
