package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	"github.com/textileio/textile-go/util"
)

var errMissingReplacement = errors.New("missing replacement value")
var errMissingPath = errors.New("missing key path")

func init() {
	register(&configCmd{})
}

type configCmd struct {
	Client ClientOptions `group:"Client Options"`
}

func (x *configCmd) Name() string {
	return "config"
}

func (x *configCmd) Short() string {
	return "Get and set config values"
}

func (x *configCmd) Long() string {
	return `
The config command controls configuration variables.
It works much like 'git config'. The configuration
values are stored in a config file inside your Textile
repository.

Getting config values will report the currently active
config settings. This may differ from the values specifed
when setting values.

When changing values, valid JSON types must be used.
For example, a string should be escaped or wrapped in
single quotes (e.g., \"127.0.0.1:40600\") and arrays and
objects work fine (e.g. '{"API": "127.0.0.1:40600"}')
but should be wrapped in single quotes. Be sure to restart
the daemon for changes to take effect.

Examples:

Get the value of the 'Addresses.API' key:

  $ textile config Addresses.API
  $ textile config Addresses/API # Alternative syntax

Print the entire Textile config file to console:

  $ textile config

Set the value of the 'Addresses.API' key:

  $ textile config Addresses.API \"127.0.0.1:40600\"
`
}

func (x *configCmd) Execute(args []string) error {
	setApi(x.Client)

	patchFmt := `[
  {"op": "replace", "path": "%s", "value": %s}
]`

	var path string
	if len(args) > 0 {
		path = "/" + strings.Replace(args[0], ".", "/", -1)
	}
	if len(args) > 1 {
		patch := []byte(fmt.Sprintf(patchFmt, path, args[1]))

		res, _, err := request(PATCH, "config", params{
			payload: bytes.NewBuffer(patch),
		})
		if err != nil {
			return err
		}
		defer res.Body.Close()

		if res.StatusCode >= 400 {
			body, err := util.UnmarshalString(res.Body)
			if err != nil {
				return err
			}
			return errors.New(body)
		}

		output("Updated! Restart daemon for changes to take effect.")
		return nil
	}

	res, err := executeJsonCmd(GET, "config"+path, params{}, nil)
	if err != nil {
		return err
	}

	output(res)
	return nil
}
