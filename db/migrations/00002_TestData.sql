-- +goose Up
-- SQL in this section is executed when the migration is applied.
INSERT INTO boards (created_at, updated_at) VALUES (NOW(), NOW());
INSERT INTO boards (created_at, updated_at) VALUES (NOW(), NOW());
INSERT INTO rooms (board_id, status, created_at, updated_at) VALUES (1, 0, NOW(), NOW());
INSERT INTO rooms (board_id, status, created_at, updated_at) VALUES (2, 0, NOW(), NOW());
INSERT INTO tiles (board_id, game_tile, created_at, updated_at) VALUES (1, 1, NOW(), NOW());
INSERT INTO tiles (board_id, game_tile, created_at, updated_at) VALUES (1, 2, NOW(), NOW());
INSERT INTO tiles (board_id, game_tile, created_at, updated_at) VALUES (1, 3, NOW(), NOW());
INSERT INTO tiles (board_id, game_tile, created_at, updated_at) VALUES (1, 4, NOW(), NOW());
INSERT INTO tiles (board_id, game_tile, created_at, updated_at) VALUES (1, 5, NOW(), NOW());
INSERT INTO tiles (board_id, game_tile, created_at, updated_at) VALUES (1, 6, NOW(), NOW());
INSERT INTO tiles (board_id, game_tile, created_at, updated_at) VALUES (1, 7, NOW(), NOW());
INSERT INTO tiles (board_id, game_tile, created_at, updated_at) VALUES (1, 8, NOW(), NOW());
INSERT INTO tiles (board_id, game_tile, created_at, updated_at) VALUES (1, 9, NOW(), NOW());
INSERT INTO tiles (board_id, game_tile, created_at, updated_at) VALUES (2, 1, NOW(), NOW());
INSERT INTO tiles (board_id, game_tile, created_at, updated_at) VALUES (2, 2, NOW(), NOW());
INSERT INTO tiles (board_id, game_tile, created_at, updated_at) VALUES (2, 3, NOW(), NOW());
INSERT INTO tiles (board_id, game_tile, created_at, updated_at) VALUES (2, 4, NOW(), NOW());
INSERT INTO tiles (board_id, game_tile, created_at, updated_at) VALUES (2, 5, NOW(), NOW());
INSERT INTO tiles (board_id, game_tile, created_at, updated_at) VALUES (2, 6, NOW(), NOW());
INSERT INTO tiles (board_id, game_tile, created_at, updated_at) VALUES (2, 7, NOW(), NOW());
INSERT INTO tiles (board_id, game_tile, created_at, updated_at) VALUES (2, 8, NOW(), NOW());
INSERT INTO tiles (board_id, game_tile, created_at, updated_at) VALUES (2, 9, NOW(), NOW());

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DELETE FROM tiles;
DELETE FROM players;
DELETE FROM rooms;
DELETE FROM boards;
