BEGIN;

CREATE TABLE IF NOT EXISTS transactions (
    id BIGSERIAL PRIMARY KEY,
    sender_id BIGINT,
    sender_type VARCHAR(10) CHECK (sender_type IN ('User', 'Auction')),
    recipient_id BIGINT,
    recipient_type VARCHAR(10) NOT NULL CHECK (recipient_type IN ('User', 'Auction')),
    amount NUMERIC(15, 2) NOT NULL CHECK (amount > 0),
    transaction_type VARCHAR(20) NOT NULL CHECK (transaction_type IN ('Deposit', 'Refund', 'Payment')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    description TEXT
);

COMMIT;
