package cmd

import (
	"net/http"
	"strconv"

	"github.com/textileio/go-textile/pb"
)

func Feed(threadID string, offset string, limit int, mode string) error {
	var list pb.FeedItemList
	opts := map[string]string{
		"thread": threadID,
		"offset": offset,
		"limit":  strconv.Itoa(limit),
		"mode":   mode,
	}
	res, err := executeJsonPbCmd(http.MethodGet, "feed", params{opts: opts}, &list)
	if err != nil {
		return err
	}
	if list.Count > 0 {
		output(res)
	}

	if list.Next == "" {
		return nil
	}

	if err := nextPage(); err != nil {
		return err
	}

	return Feed(threadID, list.Next, limit, mode)
}
