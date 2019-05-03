package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
)

// GenerateAESKey returns 44 random bytes, 32 for the key and 12 for a nonce.
func GenerateAESKey() ([]byte, error) {
	key := make([]byte, 44)
	if _, err := rand.Read(key); err != nil {
		return nil, err
	}
	return key, nil
}

// EncryptAES performs AES-256 GCM encryption on the provided bytes with key
func EncryptAES(bytes []byte, key []byte) ([]byte, error) {
	if len(key) != 44 {
		return nil, fmt.Errorf("invalid key")
	}
	block, err := aes.NewCipher(key[:32])
	if err != nil {
		return nil, err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	ciph := aesgcm.Seal(nil, key[32:], bytes, nil)
	return ciph, nil
}

// DecryptAES uses key (:32 key, 32:12 nonce) to perform AES-256 GCM decryption on bytes.
func DecryptAES(bytes []byte, key []byte) ([]byte, error) {
	if len(key) != 44 {
		return nil, fmt.Errorf("invalid key")
	}
	block, err := aes.NewCipher(key[:32])
	if err != nil {
		return nil, err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	plain, err := aesgcm.Open(nil, key[32:], bytes, nil)
	if err != nil {
		return nil, err
	}
	return plain, nil
}
