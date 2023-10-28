package utils

import (
	"bytes"
	"strings"
)

// Converts a byte slice into a string slice
func ToStrSlice(byteFile []byte) []string {

	fileStr := bytes.NewBuffer(byteFile).String()
	stringSlice := strings.Split(fileStr, " ")
	return stringSlice

}
