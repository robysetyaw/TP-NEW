CREATE TABLE customers (
    id VARCHAR PRIMARY KEY,
    fullname VARCHAR,
    address VARCHAR,
    company_id VARCHAR,
    phone_number VARCHAR,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    created_by VARCHAR,
    updated_by VARCHAR
);