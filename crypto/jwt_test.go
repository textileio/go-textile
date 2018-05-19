package crypto

import (
	"encoding/base64"
	"strings"
	"testing"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/ed25519"
)

var publicKey = "hacXwMUbI9kfLcHEHnixSW5/VyNA0H529OihiQH78mA="
var privateKey = "f3GGB4mLcK7JUSeC0sd4piOBMuICzH99l5yVfh6VhT+FpxfAxRsj2R8twcQeeLFJbn9XI0DQfnb06KGJAfvyYA=="

var ed25519TestData = []struct {
	name        string
	tokenString string
	alg         string
	claims      map[string]interface{}
	valid       bool
}{
	{
		"Ed25519",
		"eyJhbGciOiJFZDI1NTE5IiwidHlwIjoiSldUIn0.eyJqdGkiOiJmb28iLCJzdWIiOiJiYXIifQ.A07HdQgX2_rNRjj7S4zHynLkEjjiu9BzKlSYgm0iConUN1qKTG8bfpoS7Z4StdfXWN741Iv5ZpmHaxt5Kk1LBw",
		"Ed25519",
		map[string]interface{}{"jti": "foo", "sub": "bar"},
		true,
	},
	{
		"invalid key",
		"eyJhbGciOiJFZDI1NTE5IiwidHlwIjoiSldUIn0.eyJqdGkiOiJmb28iLCJzdWIiOiJiYXIifQ.7FSQFedbbRl42nvUWJqBswvjmyMaBBLKk0opiARjxtZmQ86dVMYs5wcZ0gItVV8YLVu6F5065IFD699tVcacBA",
		"Ed25519",
		map[string]interface{}{"jti": "foo", "sub": "bar"},
		false,
	},
}

func TestEd25519Verify(t *testing.T) {
	var pk ed25519.PublicKey
	var err error
	pk, err = base64.StdEncoding.DecodeString(publicKey)
	if err != nil {
		t.Fatal(err)
	}

	for _, data := range ed25519TestData {
		parts := strings.Split(data.tokenString, ".")

		method := jwt.GetSigningMethod(data.alg)
		err := method.Verify(strings.Join(parts[0:2], "."), parts[2], pk)
		if data.valid && err != nil {
			t.Errorf("[%v] error while verifying key: %v", data.name, err)
		}
		if !data.valid && err == nil {
			t.Errorf("[%v] invalid key passed validation", data.name)
		}
	}
}

func TestEd25519Sign(t *testing.T) {
	var sk ed25519.PrivateKey
	var err error
	sk, err = base64.StdEncoding.DecodeString(privateKey)
	if err != nil {
		t.Fatal(err)
	}
	for _, data := range ed25519TestData {
		if data.valid {
			parts := strings.Split(data.tokenString, ".")
			method := jwt.GetSigningMethod(data.alg)
			sig, err := method.Sign(strings.Join(parts[0:2], "."), sk)
			if err != nil {
				t.Errorf("[%v] error signing token: %v", data.name, err)
			}
			if sig != parts[2] {
				t.Errorf("[%v] incorrect signature.\nwas:\n%v\nexpecting:\n%v", data.name, sig, parts[2])
			}
		}
	}
}

func TestGenerateEd25519Token(t *testing.T) {
	var sk ed25519.PrivateKey
	var err error
	sk, err = base64.StdEncoding.DecodeString(privateKey)
	if err != nil {
		t.Fatal(err)
	}
	claims := jwt.StandardClaims{
		Id:      "bar",
		Subject: "foo",
	}
	_, err = jwt.NewWithClaims(SigningMethodEd25519, claims).SignedString(sk)
	if err != nil {
		t.Fatal(err)
	}
}
