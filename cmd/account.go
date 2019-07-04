package cmd

import (
	"encoding/base64"
	"fmt"
	"github.com/textileio/go-textile/crypto"
	"github.com/textileio/go-textile/keypair"
	"net/http"
	"os"
	"strconv"

	"github.com/textileio/go-textile/pb"
)

func getAccount() (string, *pb.Contact, error) {
	var contact pb.Contact
	res, err := executeJsonPbCmd(http.MethodGet, "account", params{}, &contact)
	if err != nil {
		return "", nil, err
	}
	return res, &contact, err
}

func AccountGet() error {
	res, _, err := getAccount()
	if err != nil {
		return err
	}
	output(res)
	return nil
}

func AccountSeed() error {
	res, err := executeStringCmd(http.MethodGet, "account/seed", params{})
	if err != nil {
		return err
	}
	output(res)
	return nil
}

func AccountAddress() error {
	res, err := executeStringCmd(http.MethodGet, "account/address", params{})
	if err != nil {
		return err
	}
	output(res)
	return nil
}

func getAccountKeyPair() (keypair.KeyPair, error) {
	var kp keypair.KeyPair
	res, err := executeStringCmd(http.MethodGet, "account/seed", params{})
	if err != nil {
		return kp, err
	}
	return keypair.Parse(res)
}

func AccountSign(message []byte, privateKeyString string) error {
	var signed []byte

	if privateKeyString == "" {
		kp, err := getAccountKeyPair()
		if err != nil {
			return err
		}

		signed, err = kp.Sign(message)
		if err != nil {
			return err
		}
	} else  {
		kp, err := keypair.Parse(privateKeyString)
		if err != nil {
			return err
		}

		signed, err = kp.Sign(message)
	}

	sigString := base64.StdEncoding.EncodeToString(signed)
	fmt.Println(sigString)

	return nil
}


func AccountVerify(message []byte, sigString string, publicKeyString string) error {
	signed, err := base64.StdEncoding.DecodeString(sigString)
	if err != nil {
		return err
	}

	if publicKeyString == "" {
		kp, err := getAccountKeyPair()
		if err != nil {
			return err
		}

		err = kp.Verify(message, signed)
	} else {
		kp, err := keypair.Parse(publicKeyString)
		if err != nil {
			return err
		}

		err = kp.Verify(message, signed)
	}

	if err != nil {
		fmt.Errorf("fail: %s\n", err)
	} else {
		fmt.Println("pass")
	}

	return nil
}

func AccountEncrypt(message []byte, publicKeyString string) error {
	var encrypted []byte

	if publicKeyString == "" {
		kp, err := getAccountKeyPair()
		if err != nil {
			return err
		}

		publicKey, err := kp.LibP2PPubKey()
		if err != nil {
			return err
		}

		encrypted, err = crypto.Encrypt(publicKey, message)
		if err != nil {
			return err
		}
	} else {
		kp, err := keypair.Parse(publicKeyString)
		if err != nil {
			return err
		}

		encrypted, err = kp.Encrypt(message)
		if err != nil {
			return err
		}
	}


	os.Stdout.Write(encrypted)

	return nil
}

func AccountDecrypt(message []byte, privateKeyString string) error {
	var decrypted []byte

	if privateKeyString == "" {
		kp, err := getAccountKeyPair()
		if err != nil {
			return err
		}

		privateKey, err := kp.LibP2PPrivKey()
		if err != nil {
			return err
		}

		decrypted, err = crypto.Decrypt(privateKey, message)
		if err != nil {
			return err
		}
	} else {
		kp, err := keypair.Parse(privateKeyString)
		if err != nil {
			return err
		}

		decrypted, err = kp.Decrypt(message)
		if err != nil {
			return err
		}
	}

	os.Stdout.Write(decrypted)

	return nil
}

func AccountSync(wait int) error {
	results := handleSearchStream("snapshots/search", params{
		opts: map[string]string{
			"wait": strconv.Itoa(wait),
		},
	})

	var remote []pb.QueryResult
	for _, res := range results {
		if !res.Local {
			remote = append(remote, res)
		}
	}
	if len(remote) == 0 {
		output("No snapshots were found")
		return nil
	}

	var postfix string
	if len(remote) > 1 {
		postfix = "s"
	}
	if !confirm(fmt.Sprintf("Apply %d snapshot%s?", len(remote), postfix)) {
		return nil
	}

	for _, result := range remote {
		if err := applyThreadSnapshot(&result); err != nil {
			return err
		}
	}

	return nil
}
