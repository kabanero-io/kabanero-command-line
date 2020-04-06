package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
)

func Create32BKey(phrase string) string {
	h := sha256.New()
	h.Write([]byte(phrase))
	return hex.EncodeToString(h.Sum(nil))
}

func EncryptString(value string, key string) string {

	keySlice, _ := hex.DecodeString(key)
	// keySlice := []byte(key)
	acipher, _ := aes.NewCipher(keySlice)

	gcm, _ := cipher.NewGCM(acipher)

	nonce := make([]byte, gcm.NonceSize())
	io.ReadFull(rand.Reader, nonce)
	encryptedText := gcm.Seal(nonce, nonce, []byte(value), nil)
	fmt.Println("encrypted jwt", encryptedText)
	return string(encryptedText)
	// return hex.EncodeToString(encryptedText)
}

func DecryptString(encryptedVal string, key string) string {
	keySlice, _ := hex.DecodeString(key)
	// keySlice := []byte(key)
	// dataSlice, _ := hex.DecodeString(encryptedVal)
	dataSlice := []byte(encryptedVal)
	block, _ := aes.NewCipher(keySlice)
	gcm, _ := cipher.NewGCM(block)
	nonce := make([]byte, gcm.NonceSize())
	nonceSize := gcm.NonceSize()
	nonce, cipherSlice := dataSlice[:nonceSize], dataSlice[nonceSize:]
	decryptedText, _ := gcm.Open(nil, nonce, cipherSlice, nil)
	return string(decryptedText)
}
