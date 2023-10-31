CREATE TABLE transaction_headers (
    id VARCHAR PRIMARY KEY,
    date DATE,
    inv_number VARCHAR,
    customer_id VARCHAR,
    name VARCHAR,
    address VARCHAR,
    company VARCHAR,
    phone_number VARCHAR,
    tx_type VARCHAR,
    total NUMERIC,
    is_active BOOLEAN,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    created_by VARCHAR,
    updated_by VARCHAR
);
