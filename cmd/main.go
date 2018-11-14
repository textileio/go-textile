package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/fatih/color"
	"github.com/textileio/textile-go/core"
	"gopkg.in/abiosoft/ishell.v2"
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

func RunShell(shell *ishell.Shell, opts ClientOptions) {
	setApi(opts)
	printSplash(core.Version)
	shell.Run()
}

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
	DEL  method = "DELETE"
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
	url := fmt.Sprintf("%s/api/%s/%s", apiAddr, apiVersion, pth)
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

func output(value interface{}, ctx *ishell.Context) {
	if ctx != nil {
		ctx.Println(Grey(value))
	} else {
		fmt.Println(value)
	}
}

func printSplash(version string) {
	url := fmt.Sprintf("%s/api/%s", apiAddr, apiVersion)
	fmt.Println(Grey("Textile shell version v" + version))
	fmt.Println(Grey("Textile API: " + url))
	fmt.Println(Grey("type 'help' for available commands"))
}
