package generator

import (
	"crypto/rand"
	"math/big"
)

type CodeGenerator interface {
	GenerateCode(length int) (string, error)
}

type codeGenerator struct {
	// You can add fields here if needed for the generator's state or configuration
}

func NewCodeGenerator() CodeGenerator {
	return &codeGenerator{}
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func (cg *codeGenerator) GenerateCode(length int) (string, error) {
	passWord := make([]byte, length)
	for i := 0; i < length; i++ {

		// Generate a random index to select a character from the charset
		index, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		// Append the selected character to the code
		passWord[i] = charset[index.Int64()]
	}
	return string(passWord), nil
}
