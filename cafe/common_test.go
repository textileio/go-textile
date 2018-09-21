package cafe

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

var client = &http.Client{}
var (
	cafeAddr        = os.Getenv("CAFE_ADDR")
	cafeReferralKey = os.Getenv("CAFE_REFERRAL_KEY")
	cafeTokenSecret = os.Getenv("CAFE_TOKEN_SECRET")
)

// CAFE V0

func signUpUser(reg interface{}) (*http.Response, error) {
	url := fmt.Sprintf("%s/api/v0/users", cafeAddr)
	payload, err := json.Marshal(reg)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return client.Do(req)
}

func signInUser(creds interface{}) (*http.Response, error) {
	url := fmt.Sprintf("%s/api/v0/users", cafeAddr)
	payload, err := json.Marshal(creds)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return client.Do(req)
}

// CAFE V1

func profileChallenge(creq interface{}) (*http.Response, error) {
	url := fmt.Sprintf("%s/api/v1/profiles/challenge", cafeAddr)
	payload, err := json.Marshal(creq)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return client.Do(req)
}

func registerProfile(reg interface{}) (*http.Response, error) {
	url := fmt.Sprintf("%s/api/v1/profiles", cafeAddr)
	payload, err := json.Marshal(reg)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return client.Do(req)
}

func loginProfile(cha interface{}) (*http.Response, error) {
	url := fmt.Sprintf("%s/api/v1/profiles", cafeAddr)
	payload, err := json.Marshal(cha)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return client.Do(req)
}

func refreshSession(accessToken string, refreshToken string) (*http.Response, error) {
	url := fmt.Sprintf("%s/api/v1/tokens", cafeAddr)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(accessToken)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", refreshToken))
	return client.Do(req)
}

func createReferral(key string, count int, limit int, requestedBy string) (*http.Response, error) {
	url := fmt.Sprintf("%s/api/v1/referrals", cafeAddr)
	params := fmt.Sprintf("count=%d&limit=%d&requested_by=%s", count, limit, requestedBy)
	req, err := http.NewRequest("POST", fmt.Sprintf("%s?%s", url, params), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Referral-Key", key)
	return client.Do(req)
}

func listReferrals(key string) (*http.Response, error) {
	url := fmt.Sprintf("%s/api/v1/referrals", cafeAddr)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Referral-Key", key)
	return client.Do(req)
}

func pin(reader io.Reader, token string, cType string) (*http.Response, error) {
	url := fmt.Sprintf("%s/api/v1/pin", cafeAddr)
	req, err := http.NewRequest("POST", url, reader)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", cType)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	return client.Do(req)
}

// UTILS

func unmarshalJSON(body io.ReadCloser, target interface{}) error {
	b, err := ioutil.ReadAll(body)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, target)
}

func verify(_ *jwt.Token) (interface{}, error) {
	return []byte(cafeTokenSecret), nil
}
