-- +goose Up
-- +goose StatementBegin
CREATE TABLE JWTTokens (
    id SERIAL PRIMARY KEY,
    token TEXT UNIQUE NOT NULL,
    expiration TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    user_id INTEGER REFERENCES users(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
