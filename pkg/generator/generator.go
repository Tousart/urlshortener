package generator

import (
	"crypto/rand"
	"math/big"
)

type Generator struct {
}

var (
	alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
)

func NewGenerator() *Generator {
	return &Generator{}
}

func (g *Generator) GenerateRandomStringWithLength(length int) (string, error) {
	res := make([]byte, length)
	symbolsCount := int64(len(alphabet))
	for i := range length {
		num, err := rand.Int(rand.Reader, big.NewInt(symbolsCount))
		if err != nil {
			return "", err
		}
		res[i] = alphabet[num.Int64()]
	}
	return string(res), nil
}
