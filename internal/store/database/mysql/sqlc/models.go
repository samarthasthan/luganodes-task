// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package sqlc

type Deposit struct {
	ID             int32
	Blocknumber    int32
	Blocktimestamp int32
	Fee            int32
	Hash           string
	Pubkey         string
}
