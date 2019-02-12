package util

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"strings"
)

func UnmarshalString(body io.ReadCloser) (string, error) {
	data, err := ioutil.ReadAll(body)
	if err != nil {
		return "", err
	}
	return trimQuotes(string(data)), nil
}

func UnmarshalJSON(body io.ReadCloser, target interface{}) error {
	data, err := ioutil.ReadAll(body)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, target)
}

func SplitString(in string, sep string) []string {
	list := make([]string, 0)
	for _, s := range strings.Split(in, sep) {
		t := strings.TrimSpace(s)
		if t != "" {
			list = append(list, t)
		}
	}
	return list
}

func trimQuotes(s string) string {
	if len(s) > 0 && s[0] == '"' {
		s = s[1:]
	}
	if len(s) > 0 && s[len(s)-1] == '"' {
		s = s[:len(s)-1]
	}
	return s
}
