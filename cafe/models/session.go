package models

import (
	"encoding/json"
	"io"
	"io/ioutil"
)

type Session struct {
	AccessToken      string `json:"access_token"`
	ExpiresAt        int64  `json:"expires_at"`
	RefreshToken     string `json:"refresh_token"`
	RefreshExpiresAt int64  `json:"refresh_expires_at"`
	SubjectId        string `json:"subject_id"`
	TokenType        string `json:"token_type"`
}

type Response struct {
	Status  int      `json:"status,omitempty"`
	Session *Session `json:"session,omitempty"`
	Error   *string  `json:"error,omitempty"`
	Id      *string  `json:"id,omitempty"`
}

type ReferralResponse struct {
	Response
	RefCodes []string `json:"ref_codes,omitempty"`
}

func (r *Response) Read(body io.ReadCloser) error {
	return unmarshalJson(body, r)
}

func (r *ReferralResponse) Read(body io.ReadCloser) error {
	return unmarshalJson(body, r)
}

func unmarshalJson(body io.ReadCloser, target interface{}) error {
	b, err := ioutil.ReadAll(body)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, target)
}

type ReferralRequest struct {
	Key         string
	Count       int
	Limit       int
	RequestedBy string
}
