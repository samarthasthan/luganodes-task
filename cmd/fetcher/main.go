package main

import (
	"context"

	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/samarthasthan/luganodes-task/internal/fetcher/kafka"
	"github.com/samarthasthan/luganodes-task/internal/fetcher/models"

	"github.com/samarthasthan/luganodes-task/pkg/env"
	"github.com/samarthasthan/luganodes-task/pkg/logger"
)

var (
	KAFKA_PORT  string
	KAFKA_HOST  string
	WS_ENDPOINT string
	ADDRESS     string
	log         *logger.Logger
)

func init() {
	KAFKA_PORT = env.GetEnv("KAFKA_PORT", "9092")
	KAFKA_HOST = env.GetEnv("KAFKA_HOST", "localhost")
	WS_ENDPOINT = env.GetEnv("WEB3_ENDPOINT", "wss://eth-mainnet.g.alchemy.com/v2/Zu9sVXs0A9UIlF-phnzLNo48ch8QMGm6")
	ADDRESS = env.GetEnv("ADDRESS", "0x00000000219ab540356cBB839Cbe05303d7705Fa")
}

func main() {
	// Initialize the logger
	log = logger.NewLogger("fetcher")

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
			log.Printf("New block mined: %v", header.Number)

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

func processDeposit(client *ethclient.Client, tx *types.Transaction, blockNumber *big.Int, blockTimestamp uint64, kafka *kafka.KafkaProducer) {
	// Chain ID for Mainnet
	chainID := big.NewInt(1)

	var senderAddress common.Address
	var err error

	// Check transaction type and use the appropriate signer
	switch tx.Type() {
	case types.LegacyTxType:
		// For legacy transactions (pre-EIP-155)
		signer := types.NewEIP155Signer(chainID)
		senderAddress, err = types.Sender(signer, tx)

	case types.AccessListTxType:
		// For EIP-2930 access list transactions
		signer := types.NewEIP2930Signer(chainID)
		senderAddress, err = types.Sender(signer, tx)

	case types.DynamicFeeTxType:
		// For EIP-1559 dynamic fee transactions
		signer := types.NewLondonSigner(chainID)
		senderAddress, err = types.Sender(signer, tx)

	default:
		log.Printf("Unsupported transaction type: %d", tx.Type())
		return
	}

	if err != nil {
		log.Printf("Failed to extract sender's address: %v", err)
		return
	}

	// Get the transaction receipt to calculate the transaction fee
	receipt, err := client.TransactionReceipt(context.Background(), tx.Hash())
	if err != nil {
		log.Printf("Failed to get transaction receipt: %v", err)
		return
	}

	// Calculate the fee in Wei
	fee := new(big.Int).Mul(tx.GasPrice(), big.NewInt(int64(receipt.GasUsed)))

	// Convert block timestamp to a human-readable format
	blockTime := time.Unix(int64(blockTimestamp), 0).UTC()

	// Log the deposit details
	log.Printf("Deposit Details:")
	log.Printf("Block Number: %s", blockNumber.String())
	log.Printf("Block Timestamp: %s", blockTime.String())
	log.Printf("Transaction Fee: %s Wei", fee.String())
	log.Printf("Transaction Hash: %s", tx.Hash().Hex())
	log.Printf("Sender Address (Pubkey): %s", senderAddress.Hex())

	// Prepare the deposit data for Kafka
	deposit := &models.Deposit{
		BlockNumber:    int(blockNumber.Int64()),
		BlockTimestamp: int(blockTimestamp),
		Fee:            fee.Int64(),
		Hash:           tx.Hash().Hex(),
		Pubkey:         senderAddress.Hex(),
	}

	// Send the deposit details to Kafka
	kafka.Produce(deposit)
}
