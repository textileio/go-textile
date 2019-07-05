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
	var signatureBytes []byte
	var signatureString string

	if privateKeyString == "" {
		kp, err := getAccountKeyPair()
		if err != nil {
			return signatureString, err
		}

		signatureBytes, err = kp.Sign(message)
		if err != nil {
			return signatureString, err
		}
	} else  {
		kp, err := keypair.Parse(privateKeyString)
		if err != nil {
			return signatureString, err
		}

		signatureBytes, err = kp.Sign(message)
	}

	signatureString = base64.StdEncoding.EncodeToString(signatureBytes)

	return signatureString, nil
}

func AccountSign(message []byte, privateKeyString string) error {
	signatureString, err := accountSign(message, privateKeyString)
	if err != nil {
		return err
	}
	fmt.Println(signatureString)
	return nil
}

func AccountVerify(message []byte, signatureString string, publicKeyString string) error {
	signatureBytes, err := base64.StdEncoding.DecodeString(signatureString)
	if err != nil {
		return err
	}

	if publicKeyString == "" {
		kp, err := getAccountKeyPair()
		if err != nil {
			return err
		}

		err = kp.Verify(message, signatureBytes)
	} else {
		kp, err := keypair.Parse(publicKeyString)
		if err != nil {
			return err
		}

		err = kp.Verify(message, signatureBytes)
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
	// The private key is used to generate the signature
	privateKeyString, err := getAccountPrivateKey()
	if err != nil {
		return err
	}

	// The public key is used to verify the signature
	publicKeyString, err := getAccountPublicKey()
	if err != nil {
		return err
	}

	// Create the message to be signed, which is a JSON array of the
	// textile account id
	// external account id
	// timestamp
	message, err := json.Marshal([]string{publicKeyString, username, time.Now().UTC().String()})
	if err != nil {
		return err
	}

	// Create the signature
	signatureString, err := accountSign(message, privateKeyString)
	if err != nil {
		return err
	}

	// Output the instructions for the user
	fmt.Printf("Create a GitHub Gist — https://gist.github.com/new — with the following:\n\n%s\n\n", signatureString)

	// Read the input verification
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Once the GitHub Gist has been created, copy its URL and paste it here, then press <enter>\n\n")
	input, err := reader.ReadString('\n')
	fmt.Print("\n\n")
	if err != nil {
		return err
	}

	// Prepare the input for verification
	inputURL := strings.TrimSpace(input)

	// Extract the verifiable components of the input
	r := regexp.MustCompile(`^https://gist\.github\.com/([^/]+)/([^/]+)/?$`)
	matches := r.FindAllStringSubmatch(inputURL, -1)

	// Verify the input is indeed verifiable
	if len(matches) != 1 && len(matches[0]) != 3 {
		return fmt.Errorf("The URL of the GitHub Gist was not constructed as expected")
	}

	// Verify the username is as expected
	// This is to prevent someone posting the verification on a different external account
	if matches[0][1] != username {
		return fmt.Errorf("The username of the GitHub Gist was not as expected")
	}

	// The verification id, in this case it is the gistID
	verificationID := matches[0][2]

	// The verification URL, in this case it is the GitHub Gist URL for the raw content
	verificationURL := fmt.Sprintf("https://gist.githubusercontent.com/%s/%s/raw/", username, verificationID)

	// Fetch the contents of the verification URL
	verificationResponse, err := http.Get(verificationURL)
	if err != nil {
		return err
	}
	verificationResponseBody, err := util.UnmarshalString(verificationResponse.Body)
	if err != nil {
		return err
	}
	verificationResponseContent := strings.TrimSpace(verificationResponseBody)

	// Verify that the contents of the verification URL match the signature
	if verificationResponseContent != signatureString {
		return fmt.Errorf("Signature did not match what we expected\nActual:    %s\nExpected:  %s", verificationResponseContent, signatureString)
	}
	fmt.Println("Verified")

	// Complete
	return nil
}
