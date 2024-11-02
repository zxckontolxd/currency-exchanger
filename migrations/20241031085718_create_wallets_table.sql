-- +goose Up
-- +goose StatementBegin
CREATE TABLE wallets (
    id SERIAL PRIMARY KEY,
    balance_usd NUMERIC(20, 2) DEFAULT 0.0,
    balance_rub NUMERIC(20, 2) DEFAULT 0.0,
    balance_eur NUMERIC(20, 2) DEFAULT 0.0
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS wallets;
-- +goose StatementEnd
