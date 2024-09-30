package internal

import (
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
	"strings"
)

type TransactionInput struct {
	BytecodeHex   string
	NetworkID     int64
	PrivateKeyHex string
	GasLimit      int64
	GasPrice      int64
	Nonce         uint64
}

func CreateSignedTransactionHex(input TransactionInput) (string, error) {
	bytecode, err := hex.DecodeString(strings.Replace(input.BytecodeHex, "0x", "", 1))
	if err != nil {
		return "", fmt.Errorf("failed to decode bytecode: %v", err)
	}

	privateKey, err := crypto.HexToECDSA(strings.Replace(input.PrivateKeyHex, "0x", "", 1))
	if err != nil {
		return "", fmt.Errorf("failed to decode private key: %v", err)
	}

	tx := types.NewTx(&types.LegacyTx{
		Nonce:    input.Nonce,
		GasPrice: big.NewInt(input.GasPrice),
		Gas:      uint64(input.GasLimit),
		To:       nil, // nil for contract creation
		Value:    big.NewInt(0),
		Data:     bytecode,
	})

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(big.NewInt(input.NetworkID)), privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign transaction: %v", err)
	}

	signedTxBytes, err := signedTx.MarshalBinary()
	if err != nil {
		return "", fmt.Errorf("failed to marshal transaction: %v", err)
	}

	return hex.EncodeToString(signedTxBytes), nil
}
