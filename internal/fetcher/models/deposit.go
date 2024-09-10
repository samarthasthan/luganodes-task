package models

type DepositDB struct {
	ID             int    `json:"id"`
	BlockNumber    int    `json:"blockNumber"`
	BlockTimestamp int    `json:"blockTimestamp"`
	Fee            int    `json:"fee"`
	Hash           string `json:"hash"`
	Pubkey         string `json:"pubkey"`
}

type Deposit struct {
	BlockNumber    int    `json:"blockNumber"`
	BlockTimestamp int    `json:"blockTimestamp"`
	Fee            int    `json:"fee"`
	Hash           string `json:"hash"`
	Pubkey         string `json:"pubkey"`
}
