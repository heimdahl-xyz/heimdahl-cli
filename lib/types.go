package lib

import (
	"math/big"
)

type FungibleTokenTransfer struct {
	Timestamp    int64    `json:"timestamp"`
	FromAddress  string   `json:"from_address"`
	FromOwner    string   `json:"from_owner,omitempty"`
	ToAddress    string   `json:"to_addresss"`
	ToOwner      string   `json:"to_owner,omitempty"`
	Amount       *big.Int `json:"amount"`
	TokenAddress string   `json:"token_address"`
	Symbol       string   `json:"symbol"`
	Chain        string   `json:"chain"`
	Network      string   `json:"network"`
	TxHash       string   `json:"tx_hash"`
	Decimals     uint8    `json:"decimals"`
	Position     uint64   `json:"position"`
}
