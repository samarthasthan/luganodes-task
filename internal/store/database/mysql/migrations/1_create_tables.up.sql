CREATE TABLE IF NOT EXISTS Deposit (
    id INT AUTO_INCREMENT PRIMARY KEY,
    blockNumber INT NOT NULL,
    blockTimestamp INT NOT NULL,
    fee BIGINT NOT NULL,
    hash VARCHAR(255) NOT NULL,
    pubkey VARCHAR(255) NOT NULL
);