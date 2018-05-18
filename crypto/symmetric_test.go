package crypto

import (
	"testing"
)

var plaintext = "yoyoyoyo!"
var key string
var ciph []byte

func TestEncryptAES(t *testing.T) {
	var err error
	key, ciph, err = EncryptAES([]byte(plaintext))
	if err != nil {
		t.Error(err)
		return
	}
}

func TestDecryptAES(t *testing.T) {
	plain, err := DecryptAES(ciph, key)
	if err != nil {
		t.Error(err)
		return
	}
	if string(plain) != plaintext {
		t.Error("decrypt aes failed")
	}
}
