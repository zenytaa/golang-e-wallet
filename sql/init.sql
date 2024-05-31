\c postgres

DROP DATABASE IF EXISTS e_wallet_db;
CREATE DATABASE e_wallet_db;

\c e_wallet_db;

CREATE TABlE users(
    id BIGSERIAL PRIMARY KEY,
    email VARCHAR UNIQUE NOT NULL,
    username VARCHAR NOT NULL,
    password VARCHAR NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE TABlE wallets(
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id),
    wallet_number CHAR(13) UNIQUE NOT NULL,
    balance DECIMAL NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE TABLE password_resets(
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id),
    token VARCHAR NOT NULL,
    expired_at TIMESTAMP NOT NULL DEFAULT NOW()+INTERVAL'10 minutes',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE TABLE source_of_funds(
    id BIGSERIAL PRIMARY KEY,
    fund_name TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE TABLE transaction_types (
    id BIGSERIAL PRIMARY KEY, 
    name VARCHAR,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE TABLE transactions(
    id BIGSERIAL PRIMARY KEY,
    amount DECIMAL NOT NULL CHECK(amount > 0),
    description TEXT NOT NULL,
    sender_wallet_id BIGINT NOT NULL REFERENCES wallets(id),
    recipient_wallet_id BIGINT NOT NULL REFERENCES wallets(id),
    source_of_fund_id BIGINT NOT NULL REFERENCES source_of_funds(id),
    transaction_type_id BIGINT NOT NULL REFERENCES transaction_types(id),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
);

INSERT INTO transaction_types(name) VALUES
('Transfer'),
('Top-Up');

INSERT INTO source_of_funds(fund_name) VALUES 
('Bank Transfer'), 
('Credit Card'), 
('Cash'), 
('Reward');
