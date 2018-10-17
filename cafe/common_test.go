package cafe

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

var client = &http.Client{}
var cafeAddr = os.Getenv("CAFE_ADDR")

func pin(reader io.Reader, token string, cType string) (*http.Response, error) {
	url := fmt.Sprintf("%s/api/v1/pin", cafeAddr)
	req, err := http.NewRequest("POST", url, reader)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", cType)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	return client.Do(req)
}

// UTILS

func unmarshalJSON(body io.ReadCloser, target interface{}) error {
	b, err := ioutil.ReadAll(body)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, target)
}
