CREATE TABLE IF NOT EXISTS accounts (
    id SERIAL PRIMARY KEY,
    owner VARCHAR(255) NOT NULL,
    balance BIGINT NOT NULL CHECK (balance >= 0),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS transfers (
    id SERIAL PRIMARY KEY,
    from_account_id INT NOT NULL REFERENCES accounts(id),
    to_account_id INT NOT NULL REFERENCES accounts(id),
    amount BIGINT NOT NULL CHECK (amount > 0),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for performance
CREATE INDEX idx_transfers_from ON transfers(from_account_id);
CREATE INDEX idx_transfers_to ON transfers(to_account_id);

-- Seed some data for manual testing
INSERT INTO accounts (owner, balance) VALUES ('Alice', 10000); -- 100.00
INSERT INTO accounts (owner, balance) VALUES ('Bob', 0); -- 0.00
