-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS projects (
    id UUID PRIMARY KEY DEFAULT (gen_random_uuid()),
    name VARCHAR NOT NULL,
    created_at TIMESTAMP,
    template VARCHAR
);
-- +goose StatementEnd
