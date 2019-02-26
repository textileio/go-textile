package crypto_test

import (
	"testing"

	. "github.com/textileio/go-textile/crypto"
)

var symmetricTestData = struct {
	plaintext  []byte
	key        []byte
	ciphertext []byte
}{
	plaintext: []byte("yoyoyoyo!"),
}

func TestGenerateAESKey(t *testing.T) {
	key, err := GenerateAESKey()
	if err != nil {
		t.Fatal(err)
	}
	symmetricTestData.key = key
}

func TestEncryptAES(t *testing.T) {
	ciphertext, err := EncryptAES(symmetricTestData.plaintext, symmetricTestData.key)
	if err != nil {
		t.Fatal(err)
	}
	symmetricTestData.ciphertext = ciphertext
}

func TestDecryptAES(t *testing.T) {
	plaintext, err := DecryptAES(symmetricTestData.ciphertext, symmetricTestData.key)
	if err != nil {
		t.Fatal(err)
	}
	if string(symmetricTestData.plaintext) != string(plaintext) {
		t.Error("decrypt AES failed")
	}
	key, err := GenerateAESKey()
	if err != nil {
		t.Fatal(err)
	}
	plaintext, err = DecryptAES(symmetricTestData.ciphertext, key)
	if err == nil {
		t.Error("decrypt AES with bad key succeeded")
	}
}
