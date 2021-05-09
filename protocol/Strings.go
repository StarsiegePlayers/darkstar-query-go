package protocol

import "bytes"

// yanked from the golang stdlib???
// why they don't provide a standard version of this, we will never know
func Clen(n []byte) int {
	i := bytes.IndexByte(n, 0)
	if i == -1 {
		i = len(n)
	}
	return i
}
