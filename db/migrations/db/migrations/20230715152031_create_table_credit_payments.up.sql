CREATE TABLE credit_payments (
    id VARCHAR PRIMARY KEY,
    inv_number VARCHAR,
    payment_date DATE,
    amount NUMERIC,
    created_at TIMESTAMP,
    created_by VARCHAR,
    updated_at TIMESTAMP,
    updated_by VARCHAR
);
