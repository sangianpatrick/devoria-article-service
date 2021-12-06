package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"log"

	"github.com/mergermarket/go-pkcs7"
)

// Crypto is a collection behavior of encryption.
type Crypto interface {
	Encrypt(plaintext string, iv string) (encrypted string)
	Decrypt(encrypted string, iv string) (plaintext string)
}

// AES256CBC concrete struct of AES CBC algorithm.
type AES256CBC struct {
	secret string
}

// NewAES256CBC is constructor.
func NewAES256CBC(secret string) Crypto {
	return &AES256CBC{
		secret: secret,
	}
}

// Encrypt returns encrypted string.
func (a AES256CBC) Encrypt(plaintext string, iv string) (encrypted string) {
	bSecret := []byte(a.secret)
	bIV := []byte(iv)
	bPlaintext, _ := pkcs7.Pad([]byte(plaintext), aes.BlockSize)
	block, _ := aes.NewCipher(bSecret)
	ciphertext := make([]byte, len(bPlaintext))
	mode := cipher.NewCBCEncrypter(block, bIV)
	mode.CryptBlocks(ciphertext, bPlaintext)

	encrypted = hex.EncodeToString(ciphertext)

	return
}

// Decrypt returns plaintext of encrypted string.
func (a AES256CBC) Decrypt(encrypted string, iv string) (plaintext string) {
	bSecret := []byte(a.secret)
	bIV := []byte(iv)

	cipherText, _ := hex.DecodeString(encrypted)

	block, _ := aes.NewCipher(bSecret)

	if len(cipherText) < aes.BlockSize {
		log.Println("cipherText too short")
		return
	}

	mode := cipher.NewCBCDecrypter(block, bIV)
	mode.CryptBlocks(cipherText, cipherText)

	cipherText, _ = pkcs7.Unpad(cipherText, aes.BlockSize)
	plaintext = string(cipherText)

	return
}
