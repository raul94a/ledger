CREATE TABLE IF NOT EXISTS  clients (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    surname1 VARCHAR(255) NOT NULL,
    surname2 VARCHAR(255),
    email VARCHAR(255) NOT NULL UNIQUE,
    identification VARCHAR(255) NOT NULL UNIQUE,
    nationality CHAR(2) NOT NULL,
    date_of_birth DATE NOT NULL,
    sex CHAR(1) NOT NULL CHECK (sex IN ('M', 'F', 'O')),
    address VARCHAR(255) NOT NULL,
    city VARCHAR(100) NOT NULL,
    province VARCHAR(100) NOT NULL,
    state VARCHAR(100),
    zip_code VARCHAR(20) NOT NULL,
    telephone VARCHAR(20) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP

);

CREATE TABLE IF NOT EXISTS accounts (
    id SERIAL PRIMARY KEY,
    client_id INTEGER REFERENCES clients(id),
    account_number VARCHAR(32) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE IF NOT EXISTS transactions (
    id SERIAL PRIMARY KEY,
    account_id INTEGER REFERENCES accounts(id),
    type VARCHAR(50) NOT NULL, 
    amount DECIMAL(15,2) NOT NULL,
    to_account_id INTEGER REFERENCES accounts(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP

);

CREATE TABLE IF NOT EXISTS  ledger_entries (
    id SERIAL PRIMARY KEY,
    transaction_id INTEGER REFERENCES transactions(id),
    account_id INTEGER REFERENCES accounts(id),
    type VARCHAR(50) NOT NULL, -- credit, debit
    amount DECIMAL(15,2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP

);

CREATE TABLE IF NOT EXISTS account_balances (
    id SERIAL PRIMARY KEY,
    account_id INTEGER REFERENCES accounts(id),
    balance DECIMAL(15,2)
);

CREATE  MATERIALIZED VIEW IF NOT EXISTS account_balances_mv AS
SELECT 
    account_id,
    COALESCE(
        SUM(CASE WHEN type = 'credit' THEN amount ELSE -amount END), 
        0.0
    ) AS balance,
    MAX(created_at) AS created_at,
    MAX(updated_at) AS updated_at
FROM ledger_entries
GROUP BY account_id
WITH DATA;

-- Índice para consultas rápidas
CREATE  UNIQUE INDEX IF NOT EXISTS idx_account_balances_mv_account_id ON account_balances_mv (account_id);