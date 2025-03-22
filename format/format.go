package format

import (
	"math/big"
	"strings"
)

// formatAmount formats token amount according to its decimals
func FormatAmountBigInt(amount *big.Int, decimals uint8) string {
	if amount == nil {
		return "0"
	}

	// Clone the amount to avoid modifying the original
	result := new(big.Int).Set(amount)

	// If decimals is 0, just return the amount as string
	if decimals == 0 {
		return result.String()
	}

	// Calculate the divisor (10^decimals)
	divisor := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil)

	// Calculate the integer part
	intPart := new(big.Int).Div(result, divisor)

	// Calculate the fractional part
	fracPart := new(big.Int).Mod(result, divisor)

	// Convert to string with proper padding
	fracStr := fracPart.String()

	// Pad with leading zeros if necessary
	for uint8(len(fracStr)) < decimals {
		fracStr = "0" + fracStr
	}

	// Trim trailing zeros
	fracStr = strings.TrimRight(fracStr, "0")

	// Format the final string
	if fracStr == "" {
		return intPart.String()
	}
	return intPart.String() + "." + fracStr
}
