package common

import (
	"crypto/rand"
	"math/big"
)

func RandomNString(n int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
	charsetLen := big.NewInt(int64(len(charset)))

	result := make([]byte, n)
	for i := range result {
		randomIndex, err := rand.Int(rand.Reader, charsetLen)
		if err != nil {
			panic(err) // Xử lý lỗi nếu có
		}
		result[i] = charset[randomIndex.Int64()]
	}
	return string(result)
}
