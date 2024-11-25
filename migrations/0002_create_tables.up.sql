BEGIN;

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

CREATE INDEX IF NOT EXISTS idx_lots_seller_id ON lots(seller_id);
CREATE INDEX IF NOT EXISTS idx_bids_auction_id ON bids(auction_id);
CREATE INDEX IF NOT EXISTS idx_auctions_status ON auctions(status);

COMMIT;
