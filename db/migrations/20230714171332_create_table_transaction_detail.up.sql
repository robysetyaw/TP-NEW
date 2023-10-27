CREATE TABLE transaction_details (
    id VARCHAR PRIMARY KEY,
    transaction_id VARCHAR,
    meat_id VARCHAR,
    meat_name TEXT,
    qty NUMERIC,
    price NUMERIC,
    total NUMERIC,
    is_active BOOLEAN,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    created_by TEXT,
    updated_by TEXT
);
