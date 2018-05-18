package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
)

// EncryptAES performs AES-256 GCM encryption on the provided bytes and returns
// the 32 byte key + the 12 byte nonce concatenated and base64 encoded.
func EncryptAES(bytes []byte) (string, []byte, error) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return "", nil, err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", nil, err
	}
	nonce := make([]byte, 12)
	_, err = rand.Read(nonce)
	if err != nil {
		return "", nil, err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", nil, err
	}
	ciph := aesgcm.Seal(nil, nonce, bytes, nil)
	key = append(key, nonce...)
	key64 := base64.StdEncoding.EncodeToString(key)
	return key64, ciph, nil
}

// DecryptAES used the provided 44 byte key (:32 key, 32:12 nonce) to perform AES-256 GCM decryption.
func DecryptAES(bytes []byte, key string) ([]byte, error) {
	keyb, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return nil, err
	}
	if len(keyb) != 44 {
		return nil, errors.New("invalid key")
	}
	block, err := aes.NewCipher(keyb[:32])
	if err != nil {
		return nil, err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	plain, err := aesgcm.Open(nil, keyb[32:], bytes, nil)
	if err != nil {
		return nil, err
	}
	return plain, nil
}
