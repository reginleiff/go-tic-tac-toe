-- +goose Up
CREATE TABLE boards (
	PRIMARY KEY (id),
	id		SERIAL		NOT NULL,
	created_at	TIMESTAMP	NOT NULL,
	updated_at	TIMESTAMP	NOT NULL
);

CREATE TABLE rooms (
	PRIMARY KEY (id),
	id		SERIAL		NOT NULL,
	board_id	INT		NOT NULL,
	status		INT		NOT NULL,
	created_at	TIMESTAMP	NOT NULL,
	updated_at	TIMESTAMP	NOT NULL,

	FOREIGN KEY (board_id) 
		REFERENCES boards(id)
);

CREATE TABLE players (
	PRIMARY KEY (id),
	id		SERIAL		NOT NULL,
	room_id		INT			,	
	created_at	TIMESTAMP	NOT NULL,
	updated_at	TIMESTAMP	NOT NULL,
	
	FOREIGN KEY (room_id)
		REFERENCES rooms(id)
);

CREATE TABLE tiles (
	PRIMARY KEY (id),
	id		SERIAL		NOT NULL,
	board_id	INT		NOT NULL,
	game_tile	INT		NOT NULL,
	created_at	TIMESTAMP	NOT NULL,
	updated_at	TIMESTAMP	NOT NULL,
	player_id	INT			,

	FOREIGN KEY (board_id)
		REFERENCES boards(id),

	FOREIGN KEY (player_id)
		REFERENCES players(id) 
);

-- +goose Down
DROP TABLE tiles;
DROP TABLE players;
DROP TABLE rooms;
DROP TABLE boards;
