package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"io"
)

func Create32BKey(phrase string) string {
	h := sha256.New()
	_, err := h.Write([]byte(phrase))
	if err != nil {
		panic(err.Error())
	}
	return hex.EncodeToString(h.Sum(nil))
}

func EncryptString(value string, key string) string {
	keySlice, _ := hex.DecodeString(key)
	acipher, _ := aes.NewCipher(keySlice)
	gcm, _ := cipher.NewGCM(acipher)
	nonce := make([]byte, gcm.NonceSize())
	_, err := io.ReadFull(rand.Reader, nonce)
	if err != nil {
		panic(err.Error())
	}
	encryptedText := gcm.Seal(nonce, nonce, []byte(value), nil)
	return base64.StdEncoding.EncodeToString(encryptedText)
}

func DecryptString(encryptedVal string, key string) string {
	keySlice, _ := hex.DecodeString(key)
	dataSlice, _ := base64.StdEncoding.DecodeString(encryptedVal)
	block, _ := aes.NewCipher(keySlice)
	gcm, _ := cipher.NewGCM(block)
	nonce := make([]byte, gcm.NonceSize())
	nonceSize := gcm.NonceSize()
	nonce, cipherSlice := dataSlice[:nonceSize], dataSlice[nonceSize:]
	decryptedText, _ := gcm.Open(nil, nonce, cipherSlice, nil)
	return string(decryptedText)
}
