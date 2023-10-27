CREATE TABLE companies (
    id VARCHAR PRIMARY KEY,
    company_name VARCHAR,
    address VARCHAR,
    email VARCHAR,
    phone_number VARCHAR,
    is_active BOOLEAN,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    created_by VARCHAR,
    updated_by VARCHAR
);