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

func BlockList(threadId string, offset string, limit int, dots bool) error {
	if threadId == "" {
		threadId = "default"
	}

	var nextOffset string
	opts := map[string]string{
		"thread": threadId,
		"offset": offset,
		"limit":  strconv.Itoa(limit),
		"dots":   strconv.FormatBool(dots),
	}

	if dots {
		var viz pb.BlockViz
		_, err := executeJsonPbCmd(http.MethodGet, "blocks", params{opts: opts}, &viz)
		if err != nil {
			return err
		}
		if viz.Count > 0 {
			output(viz.Dots)
		}

		if viz.Next == "" {
			return nil
		}

		nextOffset = viz.Next
	} else {
		var list pb.BlockList
		res, err := executeJsonPbCmd(http.MethodGet, "blocks", params{opts: opts}, &list)
		if err != nil {
			return err
		}
		if len(list.Items) > 0 {
			output(res)
		}

		if len(list.Items) < limit {
			return nil
		}
		nextOffset = list.Items[len(list.Items)-1].Id
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("next page...")
	if _, err := reader.ReadString('\n'); err != nil {
		return err
	}

	return BlockList(threadId, nextOffset, limit, dots)
}

func BlockMeta(id string) error {
	_, res, err := getBlockMeta(id)
	if err != nil {
		return err
	}
	output(res)
	return nil
}

func getBlockMeta(id string) (*pb.Block, string, error) {
	var block pb.Block
	res, err := executeJsonPbCmd(http.MethodGet, "blocks/"+id, params{}, &block)
	if err != nil {
		return nil, "", err
	}
	return &block, res, nil
}

func BlockRemove(id string) error {
	res, err := executeJsonCmd(http.MethodDelete, "blocks/"+util.TrimQuotes(id), params{}, nil)
	if err != nil {
		return err
	}
	output(res)
	return nil
}
