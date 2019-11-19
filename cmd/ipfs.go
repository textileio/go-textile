package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/textileio/go-textile/pb"
	"github.com/textileio/go-textile/util"
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

func IpfsPubsubPub(topic string, data string) error {
	res, err := executeStringCmd(http.MethodPost, "ipfs/pubsub/pub/"+topic, params{
		payload: strings.NewReader(data),
	})
	if err != nil {
		return err
	}
	output(res)
	return nil
}

func IpfsPubsubSubCommand(topic string) error {
	updates, err := IpfsPubsubSub(topic)
	if err != nil {
		return err
	}

	for {
		select {
		case update, ok := <-updates:
			if !ok {
				return nil
			}

			strings := new(pb.Strings)
			err = proto.Unmarshal(update.Data.Value.Value, strings)
			if err != nil {
				return err
			}
			fmt.Printf(strings.GetValues()[0])
		}
	}
}

func IpfsPubsubSub(topic string) (<-chan *pb.MobileQueryEvent, error) {
	updates := make(chan *pb.MobileQueryEvent, 10)
	go func() {
		defer close(updates)

		res, cancel, err := request(http.MethodGet, "ipfs/pubsub/sub/"+topic, params{})
		if err != nil {
			output(err.Error())
			return
		}
		defer res.Body.Close()
		defer cancel()

		if res.StatusCode >= 400 {
			body, err := util.UnmarshalString(res.Body)
			if err != nil {
				output(err.Error())
			} else {
				output(body)
			}
			return
		}

		decoder := json.NewDecoder(res.Body)
		for decoder.More() {
			var update pb.MobileQueryEvent
			if err := pbUnmarshaler.UnmarshalNext(decoder, &update); err == io.EOF {
				return
			} else if err != nil {
				output(err.Error())
				return
			}
			updates <- &update
		}
	}()

	return updates, nil
}
