CREATE TABLE IF NOT EXISTS pools (
    id UUID PRIMARY KEY DEFAULT (gen_random_uuid()),
    name VARCHAR NOT NULL,
    creation_date DATE DEFAULT CURRENT_DATE,
    input VARCHAR,
    output VARCHAR,
    project_id BIGINT REFERENCES projects(id)
);
