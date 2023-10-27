CREATE TABLE meats(
    id VARCHAR PRIMARY KEY,
    name VARCHAR,
    stock NUMERIC,
    price NUMERIC,
    is_active BOOLEAN,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    created_by VARCHAR,
    updated_by VARCHAR
)