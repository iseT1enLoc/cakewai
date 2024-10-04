-- +goose Up
CREATE TABLE USERS(
	id             VARCHAR(100),
	google_id       VARCHAR(100),
	profile_picture VARCHAR(100),
	name           VARCHAR(100),
	pass_word     VARCHAR(100),
	email         VARCHAR(100),
	phone         VARCHAR(100),
	createdAt      VARCHAR(100),
	updatedAt     VARCHAR(100),
);
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

-- +goose Down
DROP TABLE USERS;
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
