package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/textileio/textile-go/cafe/models"
	"io"
	"io/ioutil"
	"net/http"
)

func unmarshalJSON(body io.ReadCloser, target interface{}) error {
	b, err := ioutil.ReadAll(body)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, target)
}

func ProfileChallenge(chal *models.ChallengeRequest, url string) (*models.ChallengeResponse, error) {
	payload, err := json.Marshal(chal)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	resp := &models.ChallengeResponse{}
	if err := unmarshalJSON(res.Body, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

func RegisterProfile(reg *models.ProfileRegistration, url string) (*models.SessionResponse, error) {
	payload, err := json.Marshal(reg)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	resp := &models.SessionResponse{}
	if err := unmarshalJSON(res.Body, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

func LoginProfile(cha *models.SignedChallenge, url string) (*models.SessionResponse, error) {
	payload, err := json.Marshal(cha)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	resp := &models.SessionResponse{}
	if err := unmarshalJSON(res.Body, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

func CreateReferral(rreq *models.ReferralRequest, url string) (*models.ReferralResponse, error) {
	params := fmt.Sprintf("count=%d&limit=%d&requested_by=%s", rreq.Count, rreq.Limit, rreq.RequestedBy)
	req, err := http.NewRequest("POST", fmt.Sprintf("%s?%s", url, params), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Referral-Key", rreq.Key)
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	resp := &models.ReferralResponse{}
	if err := unmarshalJSON(res.Body, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

func ListReferrals(key string, url string) (*models.ReferralResponse, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Referral-Key", key)
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	resp := &models.ReferralResponse{}
	if err := unmarshalJSON(res.Body, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

func RefreshSession(accessToken string, refreshToken string, url string) (*models.SessionResponse, error) {
	// build the request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(accessToken)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "text/plain")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", refreshToken))
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	resp := &models.SessionResponse{}
	if err := unmarshalJSON(res.Body, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

func Pin(accessToken string, reader io.Reader, url string, cType string) (*models.PinResponse, error) {
	// build the request
	req, err := http.NewRequest("POST", url, reader)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", cType)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	resp := &models.PinResponse{}
	if err := unmarshalJSON(res.Body, resp); err != nil {
		return nil, err
	}
	return resp, nil
}
