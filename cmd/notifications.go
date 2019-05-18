package cmd

import (
	"net/http"
)

func NotificationList() error {
	res, err := executeJsonCmd(http.MethodGet, "notifications", params{}, nil)
	if err != nil {
		return err
	}
	output(res)
	return nil
}

func NotificationRead(id string) error {
	res, err := executeStringCmd(http.MethodPost, "notifications/"+id+"/read", params{})
	if err != nil {
		return err
	}
	output(res)
	return nil
}
