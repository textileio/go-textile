package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/segmentio/ksuid"
	"github.com/textileio/textile-go/central/models"
)

var client = &http.Client{}

var apiURL string
var registration = map[string]interface{}{
	"username": ksuid.New().String(),
	"password": ksuid.New().String(),
	"identity": map[string]string{
		"type":  "email_address",
		"value": fmt.Sprintf("%s@textile.io", ksuid.New().String()),
	},
}
var credentials = map[string]interface{}{
	"username": registration["username"],
	"password": registration["password"],
}

func TestUsers_Setup(t *testing.T) {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	apiURL = fmt.Sprintf("http://%s:%s", os.Getenv("HOSTNAME"), os.Getenv("PORT"))
}

func TestUsers_SignUp(t *testing.T) {
	stat, res, err := signUp(registration)
	if err != nil {
		t.Error(err)
		return
	}
	if stat != 201 {
		t.Errorf("got bad status: %d", stat)
		return
	}
	if res.Session == nil {
		t.Error("got bad session")
		return
	}
}

func TestUsers_SignIn(t *testing.T) {
	stat, res, err := signIn(credentials)
	if err != nil {
		t.Error(err)
		return
	}
	if stat != 200 {
		t.Errorf("got bad status: %d", stat)
		return
	}
	if res.Session == nil {
		t.Error("got bad session")
		return
	}
	credentials["password"] = "doh!"
	stat1, _, err := signIn(credentials)
	if err != nil {
		t.Error(err)
		return
	}
	if stat1 != 403 {
		t.Errorf("got bad status: %d", stat1)
		return
	}
	credentials["username"] = "bart"
	stat2, _, err := signIn(credentials)
	if err != nil {
		t.Error(err)
		return
	}
	if stat2 != 404 {
		t.Errorf("got bad status: %d", stat2)
		return
	}
}

func signUp(reg interface{}) (int, *models.Response, error) {
	url := fmt.Sprintf("%s/api/v1/users", apiURL)
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

func signIn(creds interface{}) (int, *models.Response, error) {
	url := fmt.Sprintf("%s/api/v1/users", apiURL)
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
