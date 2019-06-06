package cmd

import (
	"net/http"
	"strconv"

	"github.com/textileio/go-textile/pb"
)

func MessageAdd(threadID string, body string) error {
	res, err := addMessage(threadID, body)
	if err != nil {
		return err
	}
	output(res)
	return nil
}

func addMessage(threadID string, body string) (string, error) {
	res, err := executeJsonCmd(http.MethodPost, "threads/"+threadID+"/messages", params{
		args: []string{body},
	}, nil)

	if err != nil {
		return "", err
	}

	return res, nil
}

func MessageList(threadID string, offset string, limit int) error {
	var list pb.TextList
	opts := map[string]string{
		"thread": threadID,
		"offset": offset,
		"limit":  strconv.Itoa(limit),
	}
	res, err := executeJsonPbCmd(http.MethodGet, "messages", params{opts: opts}, &list)
	if err != nil {
		return err
	}
	if len(list.Items) > 0 {
		output(res)
	}

	if len(list.Items) < limit {
		return nil
	}

	if err := nextPage(); err != nil {
		return err
	}

	return MessageList(threadID, list.Items[len(list.Items)-1].Block, limit)
}

func MessageGet(blockID string) error {
	res, err := executeJsonCmd(http.MethodGet, "messages/"+blockID, params{}, nil)
	if err != nil {
		return err
	}
	output(res)
	return nil
}

func MessageIgnore(blockID string) error {
	return BlockIgnore(blockID)
}
