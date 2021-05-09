package protocol

import "bytes"

// Clen counts the number of bytes a given byte slice has, up to the first 0x00
// this code was yanked from the golang stdlib.
// why they don't provide a standard version of this, we will never know
func Clen(n []byte) int {
	i := bytes.IndexByte(n, 0)
	if i == -1 {
		i = len(n)
	}
	return i
}
