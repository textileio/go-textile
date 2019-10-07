package cmd

import (
	"net/http"
)

func BotsList() error {
	res, err := executeJsonCmd(http.MethodGet, "bots/list", params{}, nil)
	if err != nil {
		return err
	}
	output(res)
	return nil
}
