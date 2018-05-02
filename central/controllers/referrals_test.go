package controllers_test

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
)

type referralResponse struct {
	Status   int      `json:"status,omitempty"`
	RefCodes []string `json:"ref_codes,omitempty"`
	Error    string   `json:"error,omitempty"`
}

func (r *referralResponse) Read(body io.ReadCloser) error {
	b, err := ioutil.ReadAll(body)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, r)
}

func TestReferrals_Setup(t *testing.T) {
	apiURL = fmt.Sprintf("http://%s", os.Getenv("HOST"))
}

func TestReferrals_CreateReferral(t *testing.T) {
	num := 10
	stat, res, err := createReferral(refKey, num)
	if err != nil {
		t.Error(err)
		return
	}
	if stat != 201 {
		t.Errorf("got bad status: %d", stat)
		return
	}
	if len(res.RefCodes) != num {
		t.Error("got bad ref codes")
		return
	}
}

func TestReferrals_ListReferrals(t *testing.T) {
	stat, res, err := listReferrals(refKey)
	if err != nil {
		t.Error(err)
		return
	}
	if stat != 200 {
		t.Errorf("got bad status: %d", stat)
		return
	}
	if len(res.RefCodes) == 0 {
		t.Error("got bad ref codes")
		return
	}
}

func createReferral(key string, num int) (int, *referralResponse, error) {
	url := fmt.Sprintf("%s/api/v1/referrals?count=%d", apiURL, num)
	req, err := http.NewRequest("POST", url, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Referral-Key", key)
	res, err := client.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer res.Body.Close()
	resp := &referralResponse{}
	if err := resp.Read(res.Body); err != nil {
		return res.StatusCode, nil, err
	}
	return res.StatusCode, resp, nil
}

func listReferrals(key string) (int, *referralResponse, error) {
	url := fmt.Sprintf("%s/api/v1/referrals", apiURL)
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Referral-Key", key)
	res, err := client.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer res.Body.Close()
	resp := &referralResponse{}
	if err := resp.Read(res.Body); err != nil {
		return res.StatusCode, nil, err
	}
	return res.StatusCode, resp, nil
}
