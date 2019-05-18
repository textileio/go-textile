package cmd

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/textileio/go-textile/util"
)

func IPFSId() error {
	res, err := executeStringCmd(http.MethodGet, "ipfs/id", params{})
	if err != nil {
		return err
	}
	output(res)
	return nil
}

func IPFSSwarmConnect(address string) error {
	res, err := executeJsonCmd(http.MethodPost, "ipfs/swarm/connect", params{
		args: []string{address},
	}, nil)
	if err != nil {
		return err
	}
	output(res)
	return nil
}

func IPFSSwarmPeers(verbose bool, streams bool, latency bool, direction bool) error {
	res, err := executeJsonCmd(http.MethodGet, "ipfs/swarm/peers", params{
		opts: map[string]string{
			"verbose":   strconv.FormatBool(verbose),
			"streams":   strconv.FormatBool(streams),
			"latency":   strconv.FormatBool(latency),
			"direction": strconv.FormatBool(direction),
		},
	}, nil)
	if err != nil {
		return err
	}
	output(res)
	return nil
}


func IPFSCat(hash string, key string) error {
	res, _, err := request(http.MethodGet, "ipfs/cat/"+hash, params{
		opts: map[string]string{"key": key},
	})
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		body, err := util.UnmarshalString(res.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf(body)
	}

	if _, err := io.Copy(os.Stdout, res.Body); err != nil {
		return err
	}

	return nil
}
