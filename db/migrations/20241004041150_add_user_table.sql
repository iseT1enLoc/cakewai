-- +goose Up
CREATE TABLE USERS (
    id             SERIAL CONSTRAINT PK_ID PRIMARY KEY,
    google_id      VARCHAR(100),
    profile_picture VARCHAR(100),
    name           VARCHAR(100),
    password      VARCHAR(255),  -- Increased size for hashed passwords
    email          VARCHAR(100) UNIQUE,
    phone          VARCHAR(100) UNIQUE,
    createdAt      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updatedAt      TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

-- +goose Down
DROP TABLE USERS;
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
