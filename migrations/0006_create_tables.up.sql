CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    balance NUMERIC(15, 2) DEFAULT 0 CHECK (balance >= 0)
);

CREATE TABLE IF NOT EXISTS lots (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    starting_bid NUMERIC(15, 2) NOT NULL CHECK (starting_bid >= 0),
    seller_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status VARCHAR(20) NOT NULL CHECK (status IN ('Active', 'Closed'))
);

CREATE TABLE IF NOT EXISTS auctions (
    id BIGSERIAL PRIMARY KEY,
    lot_id BIGINT NOT NULL UNIQUE REFERENCES lots(id) ON DELETE CASCADE,
    status VARCHAR(20) NOT NULL CHECK (status IN ('Active', 'Ended')),
    winner_id BIGINT REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS bids (
    id BIGSERIAL PRIMARY KEY,
    auction_id BIGINT NOT NULL REFERENCES auctions(id) ON DELETE CASCADE,
    bidder_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    amount NUMERIC(15, 2) NOT NULL CHECK (amount > 0),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (auction_id, bidder_id)
);

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
