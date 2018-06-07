-- +goose Up
-- SQL in this section is executed when the migration is applied.
INSERT INTO boards (created_at, updated_at) VALUES (NOW(), NOW());
INSERT INTO boards (created_at, updated_at) VALUES (NOW(), NOW());
INSERT INTO rooms (board_id, status, created_at, updated_at) VALUES (1, 0, NOW(), NOW());
INSERT INTO rooms (board_id, status, created_at, updated_at) VALUES (2, 0, NOW(), NOW());
INSERT INTO tiles (board_id, created_at, updated_at) VALUES (1, NOW(), NOW());
INSERT INTO tiles (board_id, created_at, updated_at) VALUES (1, NOW(), NOW());
INSERT INTO tiles (board_id, created_at, updated_at) VALUES (1, NOW(), NOW());
INSERT INTO tiles (board_id, created_at, updated_at) VALUES (1, NOW(), NOW());
INSERT INTO tiles (board_id, created_at, updated_at) VALUES (1, NOW(), NOW());
INSERT INTO tiles (board_id, created_at, updated_at) VALUES (1, NOW(), NOW());
INSERT INTO tiles (board_id, created_at, updated_at) VALUES (1, NOW(), NOW());
INSERT INTO tiles (board_id, created_at, updated_at) VALUES (1, NOW(), NOW());
INSERT INTO tiles (board_id, created_at, updated_at) VALUES (1, NOW(), NOW());
INSERT INTO tiles (board_id, created_at, updated_at) VALUES (2, NOW(), NOW());
INSERT INTO tiles (board_id, created_at, updated_at) VALUES (2, NOW(), NOW());
INSERT INTO tiles (board_id, created_at, updated_at) VALUES (2, NOW(), NOW());
INSERT INTO tiles (board_id, created_at, updated_at) VALUES (2, NOW(), NOW());
INSERT INTO tiles (board_id, created_at, updated_at) VALUES (2, NOW(), NOW());
INSERT INTO tiles (board_id, created_at, updated_at) VALUES (2, NOW(), NOW());
INSERT INTO tiles (board_id, created_at, updated_at) VALUES (2, NOW(), NOW());
INSERT INTO tiles (board_id, created_at, updated_at) VALUES (2, NOW(), NOW());
INSERT INTO tiles (board_id, created_at, updated_at) VALUES (2, NOW(), NOW());

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DELETE * FROM tiles;
DELETE * FROM rooms;
DELETE * FROM boards;
