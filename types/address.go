package types

import (
	"bytes"
	"crypto/sha512"
	"encoding/base32"
)

const (
	checksumLenBytes = 4
	hashLenBytes = sha512.Size256
)

type Address [hashLenBytes]byte

func (a Address) String() string {
	// Compute the checksum
	checksumHash := sha512.Sum512_256(a[:])
	checksumLenBytes := checksumHash[hashLenBytes-checksumLenBytes:]

	// Append the checksum and encode as base32
	checksumAddress := append(a[:], checksumLenBytes...)
	return base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(checksumAddress)
}

func DecodeAddress(addr string) (a Address, err error) {
	// Interpret the address as base32
	decoded, err := base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(addr)
	if err != nil {
		return
	}

	// Ensure the decoded address is the correct length
	if len(decoded) != len(a) + checksumLenBytes {
		err = errWrongAddressLen
		return
	}

	// Split into address + checksum
	addressBytes := decoded[:len(a)]
	checksumBytes := decoded[len(a):]

	// Compute the expected checksum
	checksumHash := sha512.Sum512_256(addressBytes)
	expectedChecksumBytes := checksumHash[hashLenBytes-checksumLenBytes:]

	// Check the checksum
	if !bytes.Equal(expectedChecksumBytes, checksumBytes) {
		err = errWrongChecksum
		return
	}

	// Checksum is good, copy address bytes into output
	copy(a[:], addressBytes)
	return a, nil
}
