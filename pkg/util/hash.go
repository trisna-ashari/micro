package util

import (
	"fmt"
	"hash"
	"strings"
)

// CalculateFingerprint is a function uses to calculate fingerprint.
// Example output: 34:4B:9B:26:B2:27:CB:F2:7C:23:EF:BB:4D:A6:55:06:50:34:48:23
func CalculateFingerprint(hash hash.Hash, bytes []byte) string {
	var fingerprint []string
	// `Write` expects bytes. If you have a string `s`,
	// use `[]byte(s)` to coerce it to bytes.
	hash.Write(bytes)

	// This gets the finalized hash result as a byte
	// slice. The argument to `Sum` can be used to append
	// to an existing byte slice: it usually isn't needed.
	byteSlice := hash.Sum(nil)

	// SHA1 values are often printed in hex, for example
	// in git commits. Use the `%x` format verb to convert
	// a hash results to a hex string.
	for _, b := range byteSlice {
		fingerprint = append(fingerprint, fmt.Sprintf("%02X", b))
	}

	return strings.Join(fingerprint, ":")
}
