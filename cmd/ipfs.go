package cmd

import (
	"net/http"
	"strconv"
)

func IpfsPeer() error {
	res, err := executeStringCmd(http.MethodGet, "ipfs/id", params{})
	if err != nil {
		return err
	}
	output(res)
	return nil
}

func IpfsSwarmConnect(address string) error {
	res, err := executeJsonCmd(http.MethodPost, "ipfs/swarm/connect", params{
		args: []string{address},
	}, nil)
	if err != nil {
		return err
	}
	output(res)
	return nil
}

func IpfsSwarmPeers(verbose bool, streams bool, latency bool, direction bool) error {
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

func IpfsCat(hash string, key string) error {
	return executeBlobCmd(http.MethodGet, "ipfs/cat/"+hash, params{
		opts: map[string]string{"key": key},
	})
}
