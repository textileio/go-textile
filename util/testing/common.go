package testing

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/textileio/textile-go/cafe/models"
	cafe "github.com/textileio/textile-go/core/cafe"
	"io"
	"net/http"
	"os"
)

var client = &http.Client{}
var (
	CafeAddr        = os.Getenv("CAFE_ADDR")
	CafeReferralKey = os.Getenv("CAFE_REFERRAL_KEY")
)

func CreateReferral(key string, count int, limit int, requestedBy string) (*models.ReferralResponse, error) {
	req := &models.ReferralRequest{Key: key, Count: count, Limit: limit, RequestedBy: requestedBy}
	return cafe.CreateReferral(req, fmt.Sprintf("%s/api/v0/referrals", CafeAddr))
}

func ListReferrals(key string) (*models.ReferralResponse, error) {
	return cafe.ListReferrals(key, fmt.Sprintf("%s/api/v0/referrals", CafeAddr))
}

func SignUp(reg interface{}) (int, *models.Response, error) {
	url := fmt.Sprintf("%s/api/v0/users", CafeAddr)
	payload, err := json.Marshal(reg)
	if err != nil {
		return 0, nil, err
	}
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 201 {
		return res.StatusCode, nil, nil
	}

	resp := &models.Response{}
	if err := resp.Read(res.Body); err != nil {
		return res.StatusCode, nil, err
	}
	return res.StatusCode, resp, nil
}

func SignIn(creds interface{}) (int, *models.Response, error) {
	url := fmt.Sprintf("%s/api/v0/users", CafeAddr)
	payload, err := json.Marshal(creds)
	if err != nil {
		return 0, nil, err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return res.StatusCode, nil, nil
	}

	resp := &models.Response{}
	if err := resp.Read(res.Body); err != nil {
		return res.StatusCode, nil, err
	}
	return res.StatusCode, resp, nil
}

func Pin(reader io.Reader, token string, cType string) (int, *models.Response, error) {
	url := fmt.Sprintf("%s/api/v0/pin", CafeAddr)
	req, err := http.NewRequest("POST", url, reader)
	req.Header.Set("Content-Type", cType)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	res, err := client.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 201 {
		return res.StatusCode, nil, nil
	}

	resp := &models.Response{}
	if err := resp.Read(res.Body); err != nil {
		return res.StatusCode, nil, err
	}
	return res.StatusCode, resp, nil
}
