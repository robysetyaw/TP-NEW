CREATE TABLE daily_expenditures (
    id VARCHAR PRIMARY KEY,
    user_id VARCHAR,
    amount NUMERIC,
    description VARCHAR,
    is_active BOOLEAN,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    created_by VARCHAR,
    updated_by VARCHAR
);