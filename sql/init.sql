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
    user_id INT NOT NULL REFERENCES users(id),
    wallet_number CHAR(13) UNIQUE NOT NULL,
    balance DECIMAL NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE TABLE password_resets(
    id BIGSERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id),
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

CREATE TABLE transactions(
    id BIGSERIAL PRIMARY KEY,
    sender_wallet_id INT NOT NULL REFERENCES wallets(id),
    recipient_wallet_id INT NOT NULL REFERENCES wallets(id),
    recipient_username VARCHAR NOT NULL,
    amount DECIMAL NOT NULL CHECK(amount > 0),
    source_of_fund_id INT NOT NULL REFERENCES source_of_funds(id),
    description TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
);

INSERT INTO source_of_funds(id, fund_name, created_at, updated_at)
VALUES 
(1, 'Bank Transfer', NOW(), NOW()), 
(2, 'Credit Card', NOW(), NOW()), 
(3, 'Cash', NOW(), NOW()), 
(4, 'Reward', NOW(), NOW());
