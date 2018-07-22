package client

import (
	"bytes"
	"encoding/json"
	"github.com/textileio/textile-go/cafe/models"
	"net/http"
)

func SignIn(creds *models.Credentials, api string) (*models.Response, error) {
	payload, err := json.Marshal(creds)
	if err != nil {
		return nil, err
	}

	// build the request
	req, err := http.NewRequest("POST", api, bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// convert to structured response
	resp := &models.Response{}
	if err := resp.Read(res.Body); err != nil {
		return nil, err
	}
	return resp, nil
}

func SignUp(reg *models.Registration, api string) (*models.Response, error) {
	payload, err := json.Marshal(reg)
	if err != nil {
		return nil, err
	}

	// build the request
	req, err := http.NewRequest("PUT", api, bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// convert to structured response
	resp := &models.Response{}
	if err := resp.Read(res.Body); err != nil {
		return nil, err
	}
	return resp, nil
}
