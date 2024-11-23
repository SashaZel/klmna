-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS pools (
    id UUID PRIMARY KEY DEFAULT (gen_random_uuid()),
    name VARCHAR NOT NULL,
    description VARCHAR,
    created_at TIMESTAMP,
    project_id UUID REFERENCES projects(id)
);
-- +goose StatementEnd
