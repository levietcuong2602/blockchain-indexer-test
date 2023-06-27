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
	// For example, check the functions balanceOf, transfer, transferFrom
	balanceOfSignature := []byte("70a08231b98ef4ca268c9cc3f6b4590e4bfec28280db06bb5d45e689f2a360be")
	transferSignature := []byte("a9059cbb2ab09eb219583f4a59a5d0623ade346d962bcd4e46b11da047c9049b")
	transferFromSignature := []byte("23b872dd7302113369cda2901243429419bec145408fa8b352b3dd92b66c680b")
	return bytes.Contains(bytecode, balanceOfSignature) &&
		bytes.Contains(bytecode, transferSignature) &&
		bytes.Contains(bytecode, transferFromSignature)
}

// Perform bytecode analysis to check if it corresponds to an ERC721 token
// Return true if the bytecode matches ERC721 patterns, otherwise return false
// Example implementation:
// Check if the bytecode contains the transfer function signature
func IsERC721Contract(bytecode []byte) bool {
	// Check ERC-721 feature elements in bytecode
	// Returns true if the featured element is found, else returns false
	// For example, check the functions balanceOf, transfer, transferFrom
	balanceOfSignature := []byte("70a08231b98ef4ca268c9cc3f6b4590e4bfec28280db06bb5d45e689f2a360be")
	ownerOfSignature := []byte("6352211e6566aa027e75ac9dbf2423197fbd9b82b9d981a3ab367d355866aa1c")
	transferFromSignature := []byte("23b872dd7302113369cda2901243429419bec145408fa8b352b3dd92b66c680b")
	approveSignature := []byte("d8b964e6357ee70c418910c9b0ff9da8e4722edeed138382a45eb9edd75da0c1")
	return bytes.Contains(bytecode, balanceOfSignature) &&
		bytes.Contains(bytecode, ownerOfSignature) &&
		bytes.Contains(bytecode, transferFromSignature) &&
		bytes.Contains(bytecode, approveSignature)
}
