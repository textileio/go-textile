package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/textileio/textile-go/cafe/models"
	"github.com/textileio/textile-go/repo"
	"io"
	"net/http"
)

func CreateReferral(rreq *models.ReferralRequest, url string) (*models.ReferralResponse, error) {
	params := fmt.Sprintf("count=%d&limit=%d&requested_by=%s", rreq.Count, rreq.Limit, rreq.RequestedBy)
	req, err := http.NewRequest("POST", fmt.Sprintf("%s?%s", url, params), nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Referral-Key", rreq.Key)
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	resp := &models.ReferralResponse{}
	if err := resp.Read(res.Body); err != nil {
		return nil, err
	}
	return resp, nil
}

func ListReferrals(key string, url string) (*models.ReferralResponse, error) {
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Referral-Key", key)
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	resp := &models.ReferralResponse{}
	if err := resp.Read(res.Body); err != nil {
		return nil, err
	}
	return resp, nil
}

func SignUp(reg *models.Registration, url string) (*models.Response, error) {
	payload, err := json.Marshal(reg)
	if err != nil {
		return nil, err
	}

	// build the request
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// read response
	resp := &models.Response{}
	if err := resp.Read(res.Body); err != nil {
		return nil, err
	}
	return resp, nil
}

func SignIn(creds *models.Credentials, url string) (*models.Response, error) {
	payload, err := json.Marshal(creds)
	if err != nil {
		return nil, err
	}

	// build the request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// read response
	resp := &models.Response{}
	if err := resp.Read(res.Body); err != nil {
		return nil, err
	}
	return resp, nil
}

func Pin(tokens *repo.CafeTokens, reader io.Reader, url string, cType string) (*models.Response, error) {
	// build the request
	req, err := http.NewRequest("POST", url, reader)
	req.Header.Set("Content-Type", cType)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tokens.Access))
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// read response
	resp := &models.Response{}
	if err := resp.Read(res.Body); err != nil {
		return nil, err
	}
	return resp, nil
}
