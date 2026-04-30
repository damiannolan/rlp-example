// Write a simplified function to RLP-encode a string or byte array.
//
// Rules to follow:
// Single Byte: If the string is 1 byte long and it is < 0x80 (128 decimal), the RLP encoding is the byte itself.
//
// Short String: If the string is 0–55 bytes long, the RLP encoding is the string length +0x80 (128),
// followed by the string.
//
// Long String: If the string is greather than 55 bytes long, the RLP encoding is
// 0xb7 (183) + the length of the length (in bytes), followed by the actual length, followed by the string
package main

import (
	"encoding/binary"
	"fmt"
)

func Encode(data []byte) []byte {
	length := len(data)

	if length == 1 && data[0] < 0x80 {
		return data
	}

	if length <= 55 {
		prefix := byte(length + 0x80)
		return append([]byte{prefix}, data...)
	}

	lenBz := make([]byte, 8)
	binary.BigEndian.PutUint64(lenBz, uint64(length))

	i := 0
	for i < len(lenBz) && lenBz[i] == 0 {
		i++
	}

	// Trim leading zeros from the length bytes
	trimmedLen := lenBz[i:]

	prefix := byte(0xb7 + len(trimmedLen))
	res := append([]byte{prefix}, trimmedLen...)
	return append(res, data...)
}

func main() {
	enc := Encode([]byte("dog"))
	// Example: "dog" -> 0x83 + "dog"
	fmt.Printf("dog: %x\n", enc)

	enc = Encode([]byte("a"))
	// Example: "a" -> 0x61
	fmt.Printf("a:   %x\n", enc)

	enc = Encode([]byte("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"))
	// Example: "a" (repeated 55 times) ->
	fmt.Printf("a:   %x\n", enc)

	enc = Encode([]byte("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"))
	// Example: "a" (repeated 60 times) -> 0xb83c616161616161616161616161616161616161616161616161616161616161616161616161616161616161616161616161616161616161616161616161
	fmt.Printf("a:   %x\n", enc)
}
