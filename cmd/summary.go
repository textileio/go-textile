package cmd

import (
	"github.com/textileio/go-textile/pb"
	"net/http"
)

func Summary() error {
	var info pb.Summary
	res, err := executeJsonPbCmd(http.MethodGet, "summary", params{}, &info)
	if err != nil {
		return err
	}
	output(res)
	return nil
}
