package cmd

import "net/http"

func Ping(address string) error {
	res, err := executeStringCmd(http.MethodGet, "ping", params{args: []string{address}})
	if err != nil {
		return err
	}
	output(res)
	return nil
}
