-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS tasks (
    id UUID PRIMARY KEY DEFAULT (gen_random_uuid()),
    created_at TIMESTAMP,
    assigned_at TIMESTAMP,
    input VARCHAR,
    output VARCHAR,
    project_id UUID REFERENCES projects(id),
    pool_id UUID REFERENCES pools(id)
);
-- +goose StatementEnd