CREATE TABLE IF NOT EXISTS pools (
    id SERIAL PRIMARY KEY,
    name VARCHAR NOT NULL,
    creation_date DATE DEFAULT CURRENT_DATE,
    input VARCHAR,
    output VARCHAR,
    project_id BIGINT REFERENCES projects(id)
);
