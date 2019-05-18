package cmd

import (
	"net/http"
	"strconv"
)

func TokenCreate(token string, noStore bool) error {
	opts := map[string]string{
		"token": token,
		"store": strconv.FormatBool(!noStore),
	}

	res, err := executeStringCmd(http.MethodPost, "tokens", params{opts: opts})
	if err != nil {
		return err
	}
	output(res)
	return nil
}

func TokenList() error {
	res, err := executeJsonCmd(http.MethodGet, "tokens", params{}, nil)
	if err != nil {
		return err
	}
	output(res)
	return nil
}

func TokenValidate(token string) error {
	res, err := executeStringCmd(http.MethodGet, "tokens/"+token, params{})
	if err != nil {
		return err
	}
	output(res)
	return nil
}

func TokenRemove(token string) error {
	res, err := executeStringCmd(http.MethodDelete, "tokens/"+token, params{})
	if err != nil {
		return err
	}
	output(res)
	return nil
}
