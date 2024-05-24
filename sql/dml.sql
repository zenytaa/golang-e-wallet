\c e_wallet_db;

INSERT INTO users(id, email, password, created_at, updated_at) 
VALUES 
(1, 'user1@gmail.com', '$2y$10$f/CnoRo3x.ynHo6OrATmmeGi09o0CgWfWDMpplR1ztuOPyTGBpTKe', NOW(), NOW()),
(2, 'user2@gmail.com', '$2y$10$fHogE3s9ElhO/KGqdBk19ejOaFnggW4H.snKhHjXQyjT7gYI0kIiK', NOW(), NOW()),
(3, 'user3@gmail.com', '$2y$10$Wnj.Jy5lKU0mW44u4SiPCec67xL.pZcJPbY4npc1i0z1kucCJimsi', NOW(), NOW()),
(4, 'user4@gmail.com', '$2y$10$ZwhGU5KAdVy5wbX2NZh6Y.M5ucp22tXfR/ruRxQt33MlTPD4fF9pC', NOW(), NOW()),
(5, 'user5@gmail.com', '$2y$10$KX98c8RRTp8OjcrcG9hzJO1q4DtUkqL7GNgHyMFHWmlj2JPF0TrjW', NOW(), NOW());

INSERT INTO source_of_funds(id, fund_name, created_at, updated_at)
VALUES 
(1, 'Bank Transfer', NOW(), NOW()), 
(2, 'Credit Card', NOW(), NOW()), 
(3, 'Cash', NOW(), NOW()), 
(4, 'Reward', NOW(), NOW());

INSERT INTO wallets(id, user_id, wallet_number, balance, created_at, updated_at)
VALUES 
(1, 1, '9990000000001', 100000000, NOW(), NOW()),
(2, 2, '9990000000002', 100000000, NOW(), NOW()),
(3, 3, '9990000000003', 100000000, NOW(), NOW()),
(4, 4, '9990000000004', 100000000, NOW(), NOW()),
(5, 5, '9990000000005', 100000000, NOW(), NOW());

INSERT INTO transactions(id, sender_wallet_id, recipient_wallet_id, amount, source_of_fund_id, description, created_at, updated_at)
VALUES 
(1, 1, 2, 5000000, 1, 'Receive money from 9990000000002 by Bank Transfer', NOW(), NOW()),
(2, 2, 1, 5000000, 1, 'Send money to 9990000000001 by Bank Transfer', NOW(), NOW()),
(3, 3, 4, 2500000, 1, 'Receive money from 9990000000004 by Bank Transfer', NOW(), NOW()),
(4, 4, 3, 2500000, 1, 'Send money to 9990000000003 by Bank Transfer', NOW(), NOW()),
(5, 5, 5, 60000, 2, 'Top Up from Credit Card', NOW(), NOW()),
(6, 1, 1, 300000, 1, 'Top Up from Bank Transfer', NOW(), NOW()),
(7, 2, 2, 400000, 3, 'Top Up from Cash', NOW(), NOW()),
(8, 5, 3, 250000, 1, 'Receive money from 9990000000003 by Bank Transfer', NOW(), NOW()),
(9, 3, 5, 250000, 1, 'Send money to 9990000000005 by Bank Transfer', NOW(), NOW()),
(10, 4, 1, 300000, 3, 'Receive money from 9990000000001 by Cash', NOW(), NOW()),
(11, 1, 4, 300000, 3, 'Send money to 9990000000004 by Cash', NOW(), NOW()),
(12, 4, 4, 250000, 1, 'Top Up from Bank Transfer', NOW(), NOW()),
(13, 5, 5, 200000, 2, 'Top Up from Credit Card', NOW(), NOW()),
(14, 2, 3, 450000, 1, 'Receive money from 9990000000003 by Bank Transfer', NOW(), NOW()),
(15, 3, 2, 450000, 1, 'Send money to 9990000000002 by Cash', NOW(), NOW()),
(16, 1, 1, 300000, 3, 'Top Up from Cash', NOW(), NOW()),
(17, 4, 5, 200000, 3, 'Receive money from 9990000000005 by Cash', NOW(), NOW()),
(18, 5, 4, 200000, 3, 'Send money to 9990000000004 by Cash', NOW(), NOW()),
(19, 1, 1, 150000, 3, 'Top Up from Cash', NOW(), NOW()),
(20, 1, 1, 150000, 3, 'Top Up from Cash', NOW(), NOW());

INSERT INTO transactions(id, sender_wallet_id, recipient_wallet_id, amount, source_of_fund_id, description, created_at, updated_at)
VALUES 
(1, 1, 2, 5000000, 1, 'Receive money from 9990000000002 by Bank Transfer', '2024-02-23T15:16:11.520661Z', '2024-02-23T15:16:11.520661Z'),
(2, 2, 1, 5000000, 1, 'Send money to 9990000000001 by Bank Transfer', '2024-02-23T15:16:11.520661Z', '2024-02-23T15:16:11.520661Z'),
(3, 3, 4, 2500000, 1, 'Receive money from 9990000000004 by Bank Transfer', '2024-02-24T15:18:11.520661Z', '2024-02-24T15:18:11.520661Z'),
(4, 4, 3, 2500000, 1, 'Send money to 9990000000003 by Bank Transfer', '2024-02-24T15:18:11.520661Z', '2024-02-24T15:18:11.520661Z'),
(5, 5, 5, 60000, 2, 'Top Up from Credit Card', '2024-02-25T15:16:11.520661Z', '2024-02-25T15:16:11.520661Z'),
(6, 1, 1, 300000, 1, 'Top Up from Bank Transfer', '2024-02-25T15:20:11.520661Z', '2024-02-25T15:20:11.520661Z'),
(7, 2, 2, 400000, 3, 'Top Up from Cash', '2024-02-25T17:16:11.520661Z', '2024-02-25T17:16:11.520661Z'),
(8, 5, 3, 250000, 1, 'Receive money from 9990000000003 by Bank Transfer', '2024-02-25T18:16:11.520661Z', '2024-02-25T18:16:11.520661Z'),
(9, 3, 5, 250000, 1, 'Send money to 9990000000005 by Bank Transfer', '2024-02-25T18:16:11.520661Z', '2024-02-25T18:16:11.520661Z'),
(10, 4, 1, 300000, 3, 'Receive money from 9990000000001 by Cash', '2024-02-25T21:55:11.520661Z', '2024-02-25T21:55:11.520661Z'),
(11, 1, 4, 300000, 3, 'Send money to 9990000000004 by Cash', '2024-02-25T21:55:11.520661Z', '2024-02-25T21:55:11.520661Z'),
(12, 4, 4, 250000, 1, 'Top Up from Bank Transfer', '2024-02-26T10:16:11.520661Z', '2024-02-26T10:16:11.520661Z'),
(13, 5, 5, 200000, 2, 'Top Up from Credit Card', '2024-02-26T11:16:11.520661Z', '2024-02-26T11:16:11.520661Z'),
(14, 2, 3, 450000, 1, 'Receive money from 9990000000003 by Bank Transfer', '2024-02-26T12:16:11.520661Z', '2024-02-26T12:16:11.520661Z'),
(15, 3, 2, 450000, 1, 'Send money to 9990000000002 by Cash', '2024-02-26T12:16:11.520661Z', '2024-02-26T12:16:11.520661Z'),
(16, 1, 1, 300000, 3, 'Top Up from Cash', '2024-02-26T14:16:11.520661Z', '2024-02-26T14:16:11.520661Z'),
(17, 4, 5, 200000, 3, 'Receive money from 9990000000005 by Cash', '2024-02-26T15:16:11.520661Z', '2024-02-26T15:16:11.520661Z'),
(18, 5, 4, 200000, 3, 'Send money to 9990000000004 by Cash', '2024-02-26T15:16:11.520661Z', '2024-02-26T15:16:11.520661Z'),
(19, 1, 1, 150000, 3, 'Top Up from Cash', '2024-02-26T17:16:11.520661Z', '2024-02-26T17:16:11.520661Z'),
(20, 1, 1, 150000, 3, 'Top Up from Cash', '2024-02-26T18:16:11.520661Z', '2024-02-26T18:16:11.520661Z');