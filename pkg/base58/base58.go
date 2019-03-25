// Package base58 provides encoding and decoding functions for the base58 encoding used in Mazewire.
package base58

// Implementation: Leading zero bytes (the high-order bytes) are trimmed out, and then the alphabet
// defined here is used with the remaining non-zero bytes.

// The alphabet: no zero, no l, no capital I, no capital O.
const alphabet = "123456789abcdefghijkmnopqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ"

// Encode takes a number and returns its base58 encoding as a string.
func Encode(num uint64) string {
	str := make([]byte, 7)
	bytes := 0
	for num >= 58 {
		quot := num / 58
		str[bytes] = alphabet[num-58*quot]
		num = quot
		bytes++
	}
	return string(str[:bytes])
}
