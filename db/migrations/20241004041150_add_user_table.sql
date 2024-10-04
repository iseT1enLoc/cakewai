-- +goose Up
CREATE TABLE USERS(
    id VARCHAR(40) PRIMARY KEY,
    user_name VARCHAR(50),
    user_email VARCHAR(60),
    pass_word VARCHAR(100)
);
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

-- +goose Down
DROP TABLE USERS;
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
