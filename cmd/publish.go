package cmd

import (
	"net/http"
	"os"
)

func Publish(topic string) error {
	res, err := executeStringCmd(http.MethodPost, "publish", params{
		args:    []string{topic},
		payload: os.Stdin,
	})
	if err != nil {
		return err
	}
	output(res)
	return nil
}
