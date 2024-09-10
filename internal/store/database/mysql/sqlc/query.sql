-- name: InsertDeposit :exec
INSERT INTO Deposit (blockNumber, blockTimestamp, fee, hash, pubkey) VALUES (?, ?, ?, ?, ?);


-- name: GetDeposites :many
SELECT * FROM Deposit ORDER BY blockTimestamp DESC LIMIT ? OFFSET ?;