package protocol

import (
	"bytes"
	"errors"
)

var (
	ErrorPascalStringTooLong = errors.New("pascal string too long")
)

// Clen returns the position of the first null byte in a given byte slice
// for null-terminated c strings, this returns the length of the string
// function was yanked from a golang stdlib private method.
func Clen(n []byte) int {
	i := bytes.IndexByte(n, 0)
	if i == -1 {
		i = len(n)
	}
	return i
}

func WriteCString(input string) []byte {
	output := make([]byte, len(input)+1)
	copy(output, input)
	output[len(input)] = 0x00
	return output
}

func ReadCString(input []byte) string {
	length := Clen(input)
	output := make([]byte, length)
	copy(output, input[0:length])
	return string(output)
}

func ReadPascalStringStream(input []byte) (string, []byte) {
	if len(input) <= 0 {
		return "", input
	}

	length, input := input[0], input[1:]
	if length <= 0 || len(input) <= 0 {
		return "", input
	}
	return string(input[0:length]), input[length:]
}

func ReadPascalString(input []byte) string {
	if len(input) <= 0 {
		return ""
	}
	
	length, input := input[0], input[1:]
	if length <= 0 || len(input) <= 0 {
		return ""
	}
	return string(input[0:length])
}

func WritePascalString(input string) ([]byte, error) {
	if len(input) > 0xff {
		return []byte{}, ErrorPascalStringTooLong
	}

	output := make([]byte, len(input)+1)
	output[0] = byte(len(input))
	copy(output[1:], input)
	return output, nil
}
