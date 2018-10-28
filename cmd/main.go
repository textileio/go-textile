package cmd

import (
	"encoding/json"
	"fmt"
	"gopkg.in/abiosoft/ishell.v2"
	"io"
	"io/ioutil"
	"net/http"
)

type Cmd interface {
	Name() string
	Short() string
	Long() string
	Shell() *ishell.Cmd
}

func Cmds() []Cmd {
	return []Cmd{
		&peerCmd,
		&addressCmd,
	}
}

func executeStringCmd(name string, args []string) error {
	req, err := request("GET", name, args)
	if err != nil {
		return err
	}
	defer req.Body.Close()
	res, err := unmarshalString(req.Body)
	if err != nil {
		return err
	}
	fmt.Println(res)
	return nil
}

func request(meth string, cmd string, args []string) (*http.Response, error) {
	url := fmt.Sprintf("http://127.0.0.1:8000/api/v0/%s", cmd)
	req, err := http.NewRequest(meth, url, nil)
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	return client.Do(req)
}

func unmarshalJSON(body io.ReadCloser, target interface{}) error {
	data, err := ioutil.ReadAll(body)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, target)
}

func unmarshalString(body io.ReadCloser) (string, error) {
	data, err := ioutil.ReadAll(body)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
