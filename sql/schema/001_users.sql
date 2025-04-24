-- +goose Up
CREATE TABLE users(
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP NOT NULL,
	name TEXT UNIQUE NOT NULL
);

CREATE TABLE feeds(
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP NOT NULL,
	last_fetched_at TIMESTAMP,
	name TEXT NOT NULL,
	url TEXT UNIQUE NOT NULL,
	user_id UUID NOT NULL ,
	CONSTRAINT fk_user_id
	FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	

);

CREATE TABLE feed_follows(
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP NOT NULL,
	user_id UUID NOT NULL ,
	CONSTRAINT fk_user_id
	FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
	feed_id UUID NOT NULL ,
	CONSTRAINT fk_feed_id
	FOREIGN KEY (feed_id) REFERENCES feeds(id) ON DELETE CASCADE,
	CONSTRAINT unique_ids UNIQUE (user_id,feed_id)
);

CREATE TABLE posts(
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP NOT NULL,
	title TEXT,
	url TEXT UNIQUE NOT NULL,
	description TEXT,
	published_at TIMESTAMP NOT NULL,
	feed_id UUID NOT NULL ,
	CONSTRAINT fk_feed_id
	FOREIGN KEY (feed_id) REFERENCES feeds(id)


);

-- +goose Down
DROP TABLE posts;
DROP TABLE feed_follows;
DROP TABLE feeds;
DROP TABLE users;
