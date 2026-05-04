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
	"errors"
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

func Decode(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, errors.New("empty input")
	}

	b0 := data[0]

	switch {
	case b0 < 0x80:
		return data[:1], nil
	case b0 <= 0xb7:
		length := int(b0 - 0x80)

		if len(data) < 1+length {
			return nil, errors.New("input too short")
		}

		out := data[1 : 1+length]
		if length == 1 && out[0] < 0x80 {
			return nil, errors.New("non-canonical single byte encoding")
		}

		return out, nil
	case b0 <= 0xbf:
		lenOfLen := int(b0 - 0xb7)

		if len(data) < 1+lenOfLen {
			return nil, errors.New("input too short for length")
		}

		lenBytes := data[1 : 1+lenOfLen]
		if lenBytes[0] == 0 {
			return nil, errors.New("non-canonical length encoding")
		}

		length := 0
		for _, b := range lenBytes {
			length = (length << 8) | int(b)
		}

		if length <= 55 {
			return nil, errors.New("long form used for short string")
		}

		if len(data) < 1+lenOfLen+length {
			return nil, errors.New("input too short for string")
		}

		return data[1+lenOfLen : 1+lenOfLen+length], nil
	default:
		return nil, errors.New("expected byte string, got list")
	}
}

func main() {
	// Example: "dog" -> 0x83 + "dog"
	enc := Encode([]byte("dog"))
	fmt.Printf("dog: %x\n", enc)

	dec, err := Decode(enc)
	if err != nil {
		panic(err)
	}

	fmt.Printf("result: %s\n", dec)

	// Example: "a" -> 0x61
	enc = Encode([]byte("a"))
	fmt.Printf("a:   %x\n", enc)

	dec, err = Decode(enc)
	if err != nil {
		panic(err)
	}

	fmt.Printf("result: %s\n", dec)

	// Example: "a" (repeated 55 times) ->
	enc = Encode([]byte("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"))
	fmt.Printf("a:   %x\n", enc)

	dec, err = Decode(enc)
	if err != nil {
		panic(err)
	}

	fmt.Printf("result: %s\n", dec)

	// Example: "a" (repeated 60 times) -> 0xb83c616161616161616161616161616161616161616161616161616161616161616161616161616161616161616161616161616161616161616161616161
	enc = Encode([]byte("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"))
	fmt.Printf("a:   %x\n", enc)

	dec, err = Decode(enc)
	if err != nil {
		panic(err)
	}

	fmt.Printf("result: %s\n", dec)
}
