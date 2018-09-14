package testing

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
	CafeAddr        = os.Getenv("CAFE_ADDR")
	CafeReferralKey = os.Getenv("CAFE_REFERRAL_KEY")
	CafeTokenSecret = os.Getenv("CAFE_TOKEN_SECRET")
)

// CAFE V0

func SignUpUser(reg interface{}) (*http.Response, error) {
	url := fmt.Sprintf("%s/api/v0/users", CafeAddr)
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

func SignInUser(creds interface{}) (*http.Response, error) {
	url := fmt.Sprintf("%s/api/v0/users", CafeAddr)
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

func ProfileChallenge(creq interface{}) (*http.Response, error) {
	url := fmt.Sprintf("%s/api/v1/profiles/challenge", CafeAddr)
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

func RegisterProfile(reg interface{}) (*http.Response, error) {
	url := fmt.Sprintf("%s/api/v1/profiles", CafeAddr)
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

func LoginProfile(cha interface{}) (*http.Response, error) {
	url := fmt.Sprintf("%s/api/v1/profiles", CafeAddr)
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

func RefreshSession(accessToken string, refreshToken string) (*http.Response, error) {
	url := fmt.Sprintf("%s/api/v1/tokens", CafeAddr)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(accessToken)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", refreshToken))
	return client.Do(req)
}

func CreateReferral(key string, count int, limit int, requestedBy string) (*http.Response, error) {
	url := fmt.Sprintf("%s/api/v1/referrals", CafeAddr)
	params := fmt.Sprintf("count=%d&limit=%d&requested_by=%s", count, limit, requestedBy)
	req, err := http.NewRequest("POST", fmt.Sprintf("%s?%s", url, params), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Referral-Key", key)
	return client.Do(req)
}

func ListReferrals(key string) (*http.Response, error) {
	url := fmt.Sprintf("%s/api/v1/referrals", CafeAddr)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Referral-Key", key)
	return client.Do(req)
}

func Pin(reader io.Reader, token string, cType string) (*http.Response, error) {
	url := fmt.Sprintf("%s/api/v1/pin", CafeAddr)
	req, err := http.NewRequest("POST", url, reader)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", cType)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	return client.Do(req)
}

// UTILS

func UnmarshalJSON(body io.ReadCloser, target interface{}) error {
	b, err := ioutil.ReadAll(body)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, target)
}

func Verify(_ *jwt.Token) (interface{}, error) {
	return []byte(CafeTokenSecret), nil
}
