CREATE TABLE IF NOT EXISTS pools (
    id UUID PRIMARY KEY DEFAULT (gen_random_uuid()),
    name VARCHAR NOT NULL,
    created_at TIMESTAMP,
    input VARCHAR,
    output VARCHAR,
    project_id UUID REFERENCES projects(id)
);
