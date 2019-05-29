package cmd

import (
	"net/http"
)

func CafeAdd(peerId string, token string) error {
	res, err := executeJsonCmd(http.MethodPost, "cafes", params{
		args: []string{peerId},
		opts: map[string]string{"token": token},
	}, nil)
	if err != nil {
		return err
	}
	output(res)
	return nil
}

func CafeList() error {
	res, err := executeJsonCmd(http.MethodGet, "cafes", params{}, nil)
	if err != nil {
		return err
	}
	output(res)
	return nil
}

func CafeGet(cafeID string) error {
	res, err := executeJsonCmd(http.MethodGet, "cafes/"+cafeID, params{}, nil)
	if err != nil {
		return err
	}
	output(res)
	return nil
}

func CafeDelete(cafeID string) error {
	res, err := executeStringCmd(http.MethodDelete, "cafes/"+cafeID, params{})
	if err != nil {
		return err
	}
	output(res)
	return nil
}

func CafeMessages() error {
	res, err := executeStringCmd(http.MethodPost, "cafes/messages", params{})
	if err != nil {
		return err
	}
	output(res)
	return nil
}
