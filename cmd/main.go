package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fatih/color"
	"gopkg.in/abiosoft/ishell.v2"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

var (
	Grey   = color.New(color.FgHiBlack).SprintFunc()
	Green  = color.New(color.FgHiGreen).SprintFunc()
	Cyan   = color.New(color.FgHiCyan).SprintFunc()
	Yellow = color.New(color.FgHiYellow).SprintFunc()
)

var cmds []Cmd

type Cmd interface {
	Name() string
	Short() string
	Long() string
	Shell() *ishell.Cmd
}

func Cmds() []Cmd {
	return cmds
}

func register(cmd Cmd) {
	cmds = append(cmds, cmd)
}

type method string

const (
	GET  method = "GET"
	POST method = "POST"
	PUT  method = "PUT"
	DEL  method = "DEL"
)

type params struct {
	args    []string
	opts    map[string]string
	payload io.Reader
	ctype   string
}

func executeStringCmd(meth method, name string, pars params) (string, error) {
	req, err := request(meth, name, pars)
	if err != nil {
		return "", err
	}
	defer req.Body.Close()
	res, err := unmarshalString(req.Body)
	if err != nil {
		return "", err
	}
	return res, nil
}

func executeJsonCmd(meth method, name string, pars params, target interface{}) error {
	req, err := request(meth, name, pars)
	if err != nil {
		return err
	}
	defer req.Body.Close()
	if req.StatusCode >= 400 {
		res, err := unmarshalString(req.Body)
		if err != nil {
			return err
		}
		return errors.New(res)
	}
	return unmarshalJSON(req.Body, target)
}

func request(meth method, cmd string, pars params) (*http.Response, error) {
	url := fmt.Sprintf("http://127.0.0.1:8000/api/v0/%s", cmd)
	req, err := http.NewRequest(string(meth), url, pars.payload)
	if err != nil {
		return nil, err
	}
	if len(pars.args) > 0 {
		req.Header.Set("X-Textile-Args", strings.Join(pars.args, ","))
	}
	if len(pars.opts) > 0 {
		var items []string
		for k, v := range pars.opts {
			items = append(items, k+"="+v)
		}
		req.Header.Set("X-Textile-Opts", strings.Join(items, ","))
	}
	if pars.ctype != "" {
		req.Header.Set("Content-Type", pars.ctype)
	}
	client := &http.Client{}
	return client.Do(req)
}

func unmarshalString(body io.ReadCloser) (string, error) {
	data, err := ioutil.ReadAll(body)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func unmarshalJSON(body io.ReadCloser, target interface{}) error {
	data, err := ioutil.ReadAll(body)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, target)
}
