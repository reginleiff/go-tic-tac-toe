-- +goose Up
CREATE TABLE boards (
	PRIMARY KEY (id),
	id		SERIAL		NON NULL,
	created_at	TIMESTAMP	NON NULL,
	updated_at	TIMESTAMP	NON NULL
)

CREATE TABLE rooms (
	PRIMARY KEY (id),
	id		SERIAL		NON NULL,
	board_id	INT		NON NULL,
	status		INT		NON NULL,
	created_at	TIMESTAMP	NON NULL,
	updated_at	TIMESTAMP	NON NULL,

	FOREIGN KEY (board_id) 
		REFERENCES boards(id)
)

CREATE TABLE players (
	PRIMARY KEY (id),
	id		SERIAL		NON NULL,
	room_id		INT			,	
	created_at	TIMESTAMP	NON NULL,
	updated_at	TIMESTAMP	NON NULL,
	
	FOREIGN KEY (room_id)
		REFERENCES players(id)
)

CREATE TABLE tiles (
	PRIMARY KEY (id),
	id		SERIAL		NON NULL,
	board_id	INT		NON NULL,
	created_at	TIMESTAMP	NON NULL,
	updated_at	TIMESTAMP	NON NULL,
	captured_by	INT			,

	FOREIGN KEY (board_id)
		REFERENCES boards(id),

	FOREIGN KEY (captured_by)
		REFERENCES players(id) 
)

-- +goose Down
DROP TABLE tiles;
DROP TABLE players;
DROP TABLE rooms;
DROP TABLE boards;
