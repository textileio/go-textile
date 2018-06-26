package testing

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/textileio/textile-go/central/models"
)

var client = &http.Client{}
var (
	CentralApiURL = fmt.Sprintf("http://%s", os.Getenv("HOST"))
	RefKey        = os.Getenv("REF_KEY")
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

func CreateReferral(key string, num int, limit int, requested_by string) (int, *ReferralResponse, error) {
	url := fmt.Sprintf("%s/api/v1/referrals?count=%d&limit=%d&requested_by=%s", CentralApiURL, num, limit, requested_by)
	req, err := http.NewRequest("POST", url, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Referral-Key", key)
	res, err := client.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer res.Body.Close()
	resp := &ReferralResponse{}
	if err := resp.Read(res.Body); err != nil {
		return res.StatusCode, nil, err
	}
	return res.StatusCode, resp, nil
}

func ListReferrals(key string) (int, *ReferralResponse, error) {
	url := fmt.Sprintf("%s/api/v1/referrals", CentralApiURL)
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Referral-Key", key)
	res, err := client.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer res.Body.Close()
	resp := &ReferralResponse{}
	if err := resp.Read(res.Body); err != nil {
		return res.StatusCode, nil, err
	}
	return res.StatusCode, resp, nil
}

func SignUp(reg interface{}) (int, *models.Response, error) {
	url := fmt.Sprintf("%s/api/v1/users", CentralApiURL)
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
	resp := &models.Response{}
	if err := resp.Read(res.Body); err != nil {
		return res.StatusCode, nil, err
	}
	return res.StatusCode, resp, nil
}

func SignIn(creds interface{}) (int, *models.Response, error) {
	url := fmt.Sprintf("%s/api/v1/users", CentralApiURL)
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
	resp := &models.Response{}
	if err := resp.Read(res.Body); err != nil {
		return res.StatusCode, nil, err
	}
	return res.StatusCode, resp, nil
}
