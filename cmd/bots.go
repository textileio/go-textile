package cmd

import (
	"net/http"
)

// BotsList lists all enabled bots
func BotsList() error {
	res, err := executeJsonCmd(http.MethodGet, "bots/list", params{}, nil)
	if err != nil {
		return err
	}
	output(res)
	return nil
}

// BotsDisable disables a abot
func BotsDisable(id string) error {
	res, err := executeJsonCmd(http.MethodPost, "bots/disable", params{
		opts: map[string]string{"id": id},
	}, nil)
	if err != nil {
		return err
	}
	output(res)
	return nil
}

// BotsEnable enables a known bot
func BotsEnable(id string) error {
	res, err := executeJsonCmd(http.MethodPost, "bots/enable", params{
		opts: map[string]string{"id": id},
	}, nil)
	if err != nil {
		return err
	}
	output(res)
	return nil
}
