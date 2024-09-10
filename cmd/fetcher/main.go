package main

import (
	"context"
	"log"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/samarthasthan/luganodes-task/internal/fetcher/kafka"
	"github.com/samarthasthan/luganodes-task/internal/fetcher/models"

	"github.com/samarthasthan/luganodes-task/pkg/env"
)

var (
	KAFKA_PORT  string
	KAFKA_HOST  string
	WS_ENDPOINT string
	ADDRESS     string
)

func init() {
	KAFKA_PORT = env.GetEnv("KAFKA_PORT", "9092")
	KAFKA_HOST = env.GetEnv("KAFKA_HOST", "localhost")
	WS_ENDPOINT = env.GetEnv("WEB3_ENDPOINT", "wss://eth-mainnet.g.alchemy.com/v2/Zu9sVXs0A9UIlF-phnzLNo48ch8QMGm6")
	ADDRESS = env.GetEnv("ADDRESS", "0x00000000219ab540356cBB839Cbe05303d7705Fa")
}

func main() {
	// New kafka Producer
	producer := kafka.NewKafkaProducer(KAFKA_HOST + ":" + KAFKA_PORT)
	defer producer.Producer.Close()

	// Connect to the WebSocket
	client, err := ethclient.Dial(WS_ENDPOINT)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	// Subscribe to the new block header events
	headers := make(chan *types.Header)
	sub, err := client.SubscribeNewHead(context.Background(), headers)
	if err != nil {
		log.Fatalf("Failed to subscribe to new block headers: %v", err)
	}

	// Monitor each new block
	for {
		select {
		case err := <-sub.Err():
			log.Fatalf("Subscription error: %v", err)
		case header := <-headers:
			log.Printf("New block mined: %v", header)

			// Call a function to check for deposits in this block
			checkForDeposits(client, header.Number, header.Time, producer)
		}
	}
}

// checkForDeposits fetches the block and checks for deposit transactions
func checkForDeposits(client *ethclient.Client, blockNumber *big.Int, blockTimestamp uint64, kafka *kafka.KafkaProducer) {
	// Beacon deposit contract address
	contractAddress := common.HexToAddress(ADDRESS)

	// Get the block by number
	block, err := client.BlockByNumber(context.Background(), blockNumber)
	if err != nil {
		log.Printf("Failed to fetch block: %v", err)
		return
	}

	// Iterate through all transactions in the block
	for _, tx := range block.Transactions() {
		// Check if the transaction interacts with the Beacon Deposit Contract
		if tx.To() != nil && *tx.To() == contractAddress {
			log.Printf("Deposit found in tx: %s", tx.Hash().Hex())

			// Process the transaction and extract the required deposit details
			processDeposit(client, tx, blockNumber, blockTimestamp, kafka)
		}
	}
}

// processDeposit processes each deposit transaction and logs the required details
func processDeposit(client *ethclient.Client, tx *types.Transaction, blockNumber *big.Int, blockTimestamp uint64, kafka *kafka.KafkaProducer) {
	// Use the chain ID signer to extract the sender address
	chainID := big.NewInt(1) // Mainnet chain ID
	signer := types.NewEIP155Signer(chainID)

	// Extract the sender's address from the transaction
	senderAddress, err := types.Sender(signer, tx)
	if err != nil {
		log.Printf("Failed to extract sender's address: %v", err)
		return
	}

	// Get the transaction fee (Gas used * Gas price)
	receipt, err := client.TransactionReceipt(context.Background(), tx.Hash())
	if err != nil {
		log.Printf("Failed to get transaction receipt: %v", err)
		return
	}
	fee := new(big.Int).Mul(tx.GasPrice(), big.NewInt(int64(receipt.GasUsed)))

	// Log the deposit details
	log.Printf("Deposit Details:")
	log.Printf("Block Number: %s", blockNumber.String())
	log.Printf("Block Timestamp: %s", time.Unix(int64(blockTimestamp), 0).UTC().String()) // Convert UNIX timestamp to human-readable
	log.Printf("Transaction Fee: %s Wei", fee.String())
	log.Printf("Transaction Hash: %s", tx.Hash().Hex())
	log.Printf("Sender Address (Pubkey): %s", senderAddress.Hex())

	// produce the deposit to kafka
	deposit := &models.Deposit{
		BlockNumber:    int(blockNumber.Int64()),
		BlockTimestamp: int(blockTimestamp),
		Fee:            int(fee.Int64()),
		Hash:           tx.Hash().Hex(),
		Pubkey:         senderAddress.Hex(),
	}

	kafka.Produce(deposit)
}
