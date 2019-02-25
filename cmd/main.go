package cmd

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"

	"github.com/fatih/color"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/textileio/textile-go/pb"
	"github.com/textileio/textile-go/util"
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
	res, _, err := request(meth, pth, pars)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, err := util.UnmarshalString(res.Body)
	if err != nil {
		return "", err
	}

	return body, nil
}

func executeJsonCmd(meth method, pth string, pars params, target interface{}) (string, error) {
	res, _, err := request(meth, pth, pars)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		body, err := util.UnmarshalString(res.Body)
		if err != nil {
			return "", err
		}
		return "", errors.New(body)
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	if target == nil {
		target = new(interface{})
	}
	if err := json.Unmarshal(data, target); err != nil {
		return "", err
	}
	jsn, err := json.MarshalIndent(target, "", "    ")
	if err != nil {
		return "", err
	}

	return string(jsn), nil
}

func executeJsonPbCmd(meth method, pth string, pars params, target proto.Message) (string, error) {
	res, _, err := request(meth, pth, pars)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		body, err := util.UnmarshalString(res.Body)
		if err != nil {
			return "", err
		}
		return "", errors.New(body)
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	if err := pbUnmarshaler.Unmarshal(bytes.NewReader(data), target); err != nil {
		return "", err
	}
	jsn, err := pbMarshaler.MarshalToString(target)
	if err != nil {
		return "", err
	}

	return jsn, nil
}

func request(meth method, pth string, pars params) (*http.Response, func(), error) {
	apiUrl := fmt.Sprintf("%s/api/%s/%s", apiAddr, apiVersion, pth)
	req, err := http.NewRequest(string(meth), apiUrl, pars.payload)
	if err != nil {
		return nil, nil, err
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

	tr := &http.Transport{}
	client := &http.Client{Transport: tr}
	res, err := client.Do(req)
	cancel := func() {
		tr.CancelRequest(req)
	}

	return res, cancel, err
}

var errMissingSearchInfo = errors.New("missing search info")

var pbMarshaler = jsonpb.Marshaler{
	OrigName: true,
	Indent:   "    ",
}

var pbUnmarshaler = jsonpb.Unmarshaler{
	AllowUnknownFields: true,
}

func handleSearchStream(pth string, param params) []pb.QueryResult {
	var results []pb.QueryResult
	outputCh := make(chan interface{})

	cancel := func() {}
	done := make(chan struct{})
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)

	go func() {
		defer func() {
			cancel()
			done <- struct{}{}
		}()

		var res *http.Response
		var err error
		res, cancel, err = request(POST, pth, param)
		if err != nil {
			outputCh <- err.Error()
			return
		}
		defer res.Body.Close()

		if res.StatusCode >= 400 {
			body, err := util.UnmarshalString(res.Body)
			if err != nil {
				outputCh <- err.Error()
			} else {
				outputCh <- body
			}
			return
		}

		decoder := json.NewDecoder(res.Body)
		for decoder.More() {
			var result pb.QueryResult
			if err := pbUnmarshaler.UnmarshalNext(decoder, &result); err == io.EOF {
				return
			} else if err != nil {
				outputCh <- err.Error()
				return
			}
			results = append(results, result)

			out, err := pbMarshaler.MarshalToString(&result)
			if err != nil {
				outputCh <- err.Error()
				return
			}
			outputCh <- out
		}
	}()

	for {
		select {
		case val := <-outputCh:
			output(val)

		case <-quit:
			fmt.Println("Interrupted")
			if cancel != nil {
				fmt.Printf("Canceling...")
				cancel()
			}
			fmt.Print("done\n")
			os.Exit(1)

		case <-done:
			return results
		}
	}
}

// https://gist.github.com/r0l1/3dcbb0c8f6cfe9c66ab8008f55f8f28b
func confirm(q string) bool {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("%s [y/n]: ", q)

		response, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		response = strings.ToLower(strings.TrimSpace(response))
		if response == "y" || response == "yes" {
			return true
		} else if response == "n" || response == "no" {
			return false
		}
	}
}

func output(value interface{}) {
	fmt.Println(value)
}
