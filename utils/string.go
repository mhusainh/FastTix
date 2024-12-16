package utils

import (
	"crypto/rand"
	"math/big"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func RandomString(n int) string {
	b := make([]byte, n)
	for i := range b {
		index, _ := rand.Int(rand.Reader, big.NewInt(int64(len(letterBytes))))
		b[i] = letterBytes[index.Int64()]
	}
	return string(b)
}
