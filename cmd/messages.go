package cmd

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/textileio/go-textile/pb"
	"github.com/textileio/go-textile/util"
)

func MessageAdd(threadID string, body string) error {
	res, err := AddMessage(threadID, body)
	if err != nil {
		return err
	}
	output(res)
	return nil
}

func AddMessage(threadID string, body string) (string, error) {
	if threadID == "" {
		threadID = "default"
	}

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
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("next page...")
	if _, err := reader.ReadString('\n'); err != nil {
		return err
	}

	return MessageList(threadID, list.Items[len(list.Items)-1].Block, limit)
}

func MessageGet(messageID string) error {
	res, err := executeJsonCmd(http.MethodGet, "messages/"+util.TrimQuotes(messageID), params{}, nil)
	if err != nil {
		return err
	}
	output(res)
	return nil
}

func MessageIgnore(messageID string) error {
	return BlockRemove(messageID)
}
