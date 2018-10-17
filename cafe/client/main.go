package client

import (
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
