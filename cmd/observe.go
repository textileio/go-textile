package cmd

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/textileio/go-textile/pb"
	"github.com/textileio/go-textile/util"
)

func ObserveCommand(threadID string, types []string) error {
	updates, err := Observe(threadID, types)
	if err != nil {
		return err
	}

	for {
		select {
		case update, ok := <-updates:
			if !ok {
				return nil
			}

			out, err := pbMarshaler.MarshalToString(update)
			if err == io.EOF {
				break
			} else if err != nil {
				return err
			}
			output(out)
		}
	}
}

func Observe(threadID string, types []string) (<-chan *pb.FeedItem, error) {
	if threadID != "" {
		threadID = "/" + threadID
	}

	updates := make(chan *pb.FeedItem, 10)
	go func() {
		defer close(updates)

		res, cancel, err := request(http.MethodGet, "observe"+threadID, params{
			opts: map[string]string{"type": strings.Join(types, "|")},
		})
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
			var update pb.FeedItem
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
