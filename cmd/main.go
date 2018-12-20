package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/fatih/color"
)

type ClientOptions struct {
	ApiAddr    string `long:"api" description:"API address to use" default:"http://127.0.0.1:40600"`
	ApiVersion string `long:"api-version" description:"API version to use" default:"v0"`
}

var (
	Grey   = color.New(color.FgHiBlack).SprintFunc()
	Green  = color.New(color.FgHiGreen).SprintFunc()
	Cyan   = color.New(color.FgHiCyan).SprintFunc()
	Yellow = color.New(color.FgHiYellow).SprintFunc()
)

var apiAddr, apiVersion string

func setApi(opts ClientOptions) {
	apiAddr = opts.ApiAddr
	apiVersion = opts.ApiVersion
}

var cmds []Cmd

type Cmd interface {
	Name() string
	Short() string
	Long() string
}

func Cmds() []Cmd {
	return cmds
}

func register(cmd Cmd) {
	cmds = append(cmds, cmd)
}

type method string

const (
	GET   method = "GET"
	POST  method = "POST"
	PUT   method = "PUT"
	DEL   method = "DELETE"
	PATCH method = "PATCH"
)

type params struct {
	args    []string
	opts    map[string]string
	payload io.Reader
	ctype   string
}

func executeStringCmd(meth method, pth string, pars params) (string, error) {
	req, err := request(meth, pth, pars)
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

func executeJsonCmd(meth method, pth string, pars params, target interface{}) (string, error) {
	req, err := request(meth, pth, pars)
	if err != nil {
		return "", err
	}
	defer req.Body.Close()
	if req.StatusCode >= 400 {
		res, err := unmarshalString(req.Body)
		if err != nil {
			return "", err
		}
		return "", errors.New(res)
	}
	if err := unmarshalJSON(req.Body, target); err != nil {
		return "", err
	}
	data, err := json.MarshalIndent(target, "", "    ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func request(meth method, pth string, pars params) (*http.Response, error) {
	apiUrl := fmt.Sprintf("%s/api/%s/%s", apiAddr, apiVersion, pth)
	req, err := http.NewRequest(string(meth), apiUrl, pars.payload)
	if err != nil {
		return nil, err
	}
	if len(pars.args) > 0 {
		var args []string
		for _, arg := range pars.args {
			args = append(args, url.PathEscape(arg))
		}
		req.Header.Set("X-Textile-Args", strings.Join(args, ","))
	}
	if len(pars.opts) > 0 {
		var items []string
		for k, v := range pars.opts {
			items = append(items, k+"="+url.PathEscape(v))
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
	return trimQuotes(string(data)), nil
}

func unmarshalJSON(body io.ReadCloser, target interface{}) error {
	data, err := ioutil.ReadAll(body)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, target)
}

func output(value interface{}) {
	fmt.Println(value)
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

func printSplash(version string) {
	url := fmt.Sprintf("%s/api/%s", apiAddr, apiVersion)
	fmt.Println(Grey("Textile shell version v" + version))
	fmt.Println(Grey("Textile API: " + url))
	fmt.Println(Grey("type 'help' for available commands"))
}
