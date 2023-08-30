CREATE TABLE IF NOT EXISTS segments(
	slug VARCHAR(50) PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS users(
    id INT PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS segments_users(
    user_id INTEGER NOT NULL,
    segment_slug VARCHAR(50) NOT NULL,

    PRIMARY KEY(user_id, segment_slug),

    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE RESTRICT,
    FOREIGN KEY (segment_slug) REFERENCES segments(slug) ON DELETE RESTRICT
);