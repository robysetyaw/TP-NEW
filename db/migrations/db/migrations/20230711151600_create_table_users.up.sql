CREATE TABLE users (
    id VARCHAR PRIMARY KEY,
    username VARCHAR,
    password VARCHAR,
    is_active BOOLEAN,
    role VARCHAR,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    created_by VARCHAR,
    updated_by VARCHAR
);
