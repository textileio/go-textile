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
	if value == "" {
		// get
		res, err := executeJsonCmd(http.MethodGet, "config"+path, params{}, nil)
		if err != nil {
			return err
		}
		output(res)
	} else {
		// set
		patchString := fmt.Sprintf(patchFmt, path, value)
		patchBytes := []byte(patchString)
		patchBuffer := bytes.NewBuffer(patchBytes)

		// request
		res, _, err := request(http.MethodPatch, "config", params{
			payload: patchBuffer,
		})
		if err != nil {
			return err
		}
		defer res.Body.Close()

		// check
		if res.StatusCode >= 400 {
			body, err := util.UnmarshalString(res.Body)
			if err != nil {
				return err
			}
			suggestion := fmt.Sprintf(`textile config '%s' '"%s"'`, name, value)
			return fmt.Errorf("Applying the configuration patch failed, you may have forgotten to JSON escape your input value.\n\nError: %s\n\n%s\n\nYou could try this instead:\n\n%s", body, patchString, suggestion)
		}
		output("Updated! Restart daemon for changes to take effect.")
		return nil
	}

	return nil
}
