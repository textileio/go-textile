package cmd

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/textileio/go-textile/crypto"
	"github.com/textileio/go-textile/keypair"
	"github.com/textileio/go-textile/util"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

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

func getAccountPrivateKey() (string, error) {
	var privateKey string
	privateKey, err := executeStringCmd(http.MethodGet, "account/seed", params{})
	if err != nil {
		return privateKey, err
	}
	return privateKey, nil
}

func AccountSeed() error {
	privateKey, err := getAccountPrivateKey()
	if err != nil {
		return err
	}
	output(privateKey)
	return nil
}

func getAccountPublicKey() (string, error) {
	var publicKey string
	publicKey, err := executeStringCmd(http.MethodGet, "account/address", params{})
	if err != nil {
		return publicKey, err
	}
	return publicKey, nil
}

func AccountAddress() error {
	publicKey, err := getAccountPublicKey()
	if err != nil {
		return err
	}
	output(publicKey)
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

func accountSign(message []byte, privateKeyString string) (string, error) {
	var signed []byte
	var sigString string

	if privateKeyString == "" {
		kp, err := getAccountKeyPair()
		if err != nil {
			return sigString, err
		}

		signed, err = kp.Sign(message)
		if err != nil {
			return sigString, err
		}
	} else  {
		kp, err := keypair.Parse(privateKeyString)
		if err != nil {
			return sigString, err
		}

		signed, err = kp.Sign(message)
	}

	sigString = base64.StdEncoding.EncodeToString(signed)

	return sigString, nil
}

func AccountSign(message []byte, privateKeyString string) error {
	sigString, err := accountSign(message, privateKeyString)
	if err != nil {
		return err
	}
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

func AccountAuthGithub(username string) error {
	privateKeyString, err := getAccountPrivateKey()
	if err != nil {
		return err
	}

	publicKeyString, err := getAccountPublicKey()
	if err != nil {
		return err
	}

	message, err := json.Marshal([]string{publicKeyString, username, time.Now().UTC().String()})
	if err != nil {
		return err
	}

	sigString, err := accountSign(message, privateKeyString)
	if err != nil {
		return err
	}

	fmt.Printf("Post a gist with the following:\n\n%s\n\n", sigString)

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Once posted, post the gist URL here, then press <enter>\n\n")
	input, err := reader.ReadString('\n')
	if err != nil {
		return err
	}

	url := strings.TrimSpace(input)
	fmt.Print("\n\n")

	r := regexp.MustCompile(`^https://gist\.github\.com/([^/]+)/([^/]+)/?$`)

	matches := r.FindAllStringSubmatch(url, -1)

	if len(matches) != 1 && len(matches[0]) != 3 {
		return fmt.Errorf("Gist URL was not constructed as expected")
	}

	if matches[0][1] != username {
		return fmt.Errorf("Gist Username was not as expected")
	}

	gistID := matches[0][2]

	rawURL := fmt.Sprintf("https://gist.githubusercontent.com/%s/%s/raw/", username, gistID)

	resp, err := http.Get(rawURL)
	if err != nil {
		return err
	}

	rawResult, err := util.UnmarshalString(resp.Body)
	if err != nil {
		return err
	}
	result := strings.TrimSpace(rawResult)

	if result != sigString {
		return fmt.Errorf("Signature did not match what we expected\nActual:    %s\nExpected:  %s", result, sigString)
	}

	// signature is ok, write it all the the auth thread

	fmt.Println("Verified", result) // false

	return nil
}
