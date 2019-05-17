package cmd

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"

	"github.com/textileio/go-textile/util"
)

func Config(name string, value string) error {
	patchFmt := `[
  {"op": "replace", "path": "%s", "value": %s}
]`

	var path string
	if name != "" {
		path = "/" + strings.Replace(name, ".", "/", -1)
	}
	if value != "" {
		patch := []byte(fmt.Sprintf(patchFmt, path, value))

		res, _, err := request(http.MethodPatch, "config", params{
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
			return fmt.Errorf(body)
		}

		output("Updated! Restart daemon for changes to take effect.")
		return nil
	}

	res, err := executeJsonCmd(http.MethodGet, "config"+path, params{}, nil)
	if err != nil {
		return err
	}

	output(res)
	return nil
}
