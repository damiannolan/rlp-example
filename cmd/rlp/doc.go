// Visual map of the RLP encoding byte space.
//
// [0x00 - 0x7f]: Single bytes (the data itself).
// [0x80 - 0xb7]: Short Strings (0x80 + length).
// [0xb8 - 0xbf]: Long Strings (0xb7 + length of length).
// [0xc0 - 0xf7]: Short Lists (0xc0 + length).
// [0xf8 - 0xff]: Long Lists (0xf7 + length of length).
package main
