package repotools

import "encoding/hex"

// UUIDVersion4 takes an array of 16 (random) bytes and converts it to a UUIDv4 value.
func UUIDVersion4(u [16]byte) string {
	// https://en.wikipedia.org/wiki/Universally_unique_identifier#Version_4_.28random.29
	const dash = '-'

	// 13th character is "4"
	u[6] = (u[6] & 0x0f) | 0x40 // Version 4
	// 17th character is "8", "9", "a", or "b"
	u[8] = (u[8] & 0x3f) | 0x80 // Variant is 10

	var scratch [36]byte

	hex.Encode(scratch[:8], u[0:4])
	scratch[8] = dash
	hex.Encode(scratch[9:13], u[4:6])
	scratch[13] = dash
	hex.Encode(scratch[14:18], u[6:8])
	scratch[18] = dash
	hex.Encode(scratch[19:23], u[8:10])
	scratch[23] = dash
	hex.Encode(scratch[24:], u[10:])

	return string(scratch[:])
}
