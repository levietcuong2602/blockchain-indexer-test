package address

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"strings"

	"golang.org/x/crypto/sha3"
)

// Decode decodes a hex string with 0x prefix.
func Remove0x(input string) string {
	if strings.HasPrefix(input, "0x") {
		return input[2:]
	}

	return input
}

// Hex returns an EIP55-compliant hex string representation of the address.
func EIP55Checksum(unchecksummed string) (string, error) {
	v := []byte(Remove0x(strings.ToLower(unchecksummed)))

	if _, err := hex.DecodeString(string(v)); err != nil {
		return "", fmt.Errorf("failed to decode a string: %w", err)
	}

	sha := sha3.NewLegacyKeccak256()
	if _, err := sha.Write(v); err != nil {
		return "", fmt.Errorf("failed to write sha: %w", err)
	}
	hash := sha.Sum(nil)

	result := v
	for i := 0; i < len(result); i++ {
		hashByte := hash[i/2]
		if i%2 == 0 {
			hashByte >>= 4
		} else {
			hashByte &= 0xf
		}
		if result[i] > '9' && hashByte > 7 {
			result[i] -= 32
		}
	}
	val := string(result)

	return "0x" + val, nil
}

// Perform bytecode analysis to check if it corresponds to an ERC20 token
// Return true if the bytecode matches ERC20 patterns, otherwise return false
// Example implementation:
// Check if the bytecode contains the transfer function signature
func IsERC20Contract(bytecode []byte) bool {
	// Check ERC-20 feature elements in bytecode
	// Returns true if the featured element is found, else returns false
	// For example, check the functions transfer and event Transfer
	transferSignature := []byte("a9059cbb")
	transferEventSignature := []byte("dd62ed3e")
	return bytes.Contains(bytecode, transferSignature) &&
		bytes.Contains(bytecode, transferEventSignature)
}

// Perform bytecode analysis to check if it corresponds to an ERC721 token
// Return true if the bytecode matches ERC721 patterns, otherwise return false
// Example implementation:
// Check if the bytecode contains the transfer function signature
func IsERC721Contract(bytecode []byte) bool {
	// Check ERC-721 feature elements in bytecode
	// Returns true if the featured element is found, else returns false
	// For example, check the functions transferFrom and event Transfer
	transferFromSignature := []byte("23b872dd") // Hàm transferFrom có mã hash là keccak256("transferFrom(address,address,uint256)")
	transferEventSignature := []byte("ddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef")

	return bytes.Contains(bytecode, transferFromSignature) &&
		bytes.Contains(bytecode, transferEventSignature)
}
