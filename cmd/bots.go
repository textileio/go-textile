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

func BotsDisable(id string) error {
	res, err := executeJsonCmd(http.MethodGet, "bots/disable", params{
		opts: map[string]string{"id": id},
	}, nil)
	if err != nil {
		return err
	}
	output(res)
	return nil
}

func BotsEnable(id string) error {
	res, err := executeJsonCmd(http.MethodGet, "bots/enable", params{
		opts: map[string]string{"id": id},
	}, nil)
	if err != nil {
		return err
	}
	output(res)
	return nil
}
