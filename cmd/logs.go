package cmd

import (
	"net/http"
	"strconv"
)

func Logs(subsystem string, level string, texOnly bool) error {
	if subsystem != "" {
		subsystem = "/" + subsystem
	}
	var method method
	method = http.MethodGet
	if level != "" {
		method = http.MethodPost
	}

	opts := map[string]string{
		"subsystem": subsystem,
		"level":     level,
		"tex-only":  strconv.FormatBool(texOnly),
	}
	res, err := executeJsonCmd(method, "logs"+subsystem, params{opts: opts}, nil)
	if err != nil {
		return err
	}
	output(res)
	return nil
}
