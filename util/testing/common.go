package testing

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/textileio/textile-go/cafe/models"
)

var client = &http.Client{}
var (
	CafeAddr        = os.Getenv("CAFE_ADDR")
	CafeReferralKey = os.Getenv("CAFE_REFERRAL_KEY")
)

type ReferralResponse struct {
	Status   int      `json:"status,omitempty"`
	RefCodes []string `json:"ref_codes,omitempty"`
	Error    string   `json:"error,omitempty"`
}

func (r *ReferralResponse) Read(body io.ReadCloser) error {
	b, err := ioutil.ReadAll(body)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, r)
}

func CreateReferral(key string, num int, limit int, requestedBy string) (int, *ReferralResponse, error) {
	url := fmt.Sprintf("%s/api/v0/referrals?count=%d&limit=%d&requested_by=%s", CafeAddr, num, limit, requestedBy)
	req, err := http.NewRequest("POST", url, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Referral-Key", key)
	res, err := client.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 201 {
		return res.StatusCode, nil, nil
	}

	resp := &ReferralResponse{}
	if err := resp.Read(res.Body); err != nil {
		return res.StatusCode, nil, err
	}
	return res.StatusCode, resp, nil
}

func ListReferrals(key string) (int, *ReferralResponse, error) {
	url := fmt.Sprintf("%s/api/v0/referrals", CafeAddr)
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Referral-Key", key)
	res, err := client.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return res.StatusCode, nil, nil
	}

	resp := &ReferralResponse{}
	if err := resp.Read(res.Body); err != nil {
		return res.StatusCode, nil, err
	}
	return res.StatusCode, resp, nil
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
