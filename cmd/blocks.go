package cmd

import (
	"net/http"
	"strconv"

	"github.com/textileio/go-textile/pb"
)

func BlockList(threadID string, offset string, limit int, dots bool) error {
	var nextOffset string
	opts := map[string]string{
		"thread": threadID,
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

	if err := nextPage(); err != nil {
		return err
	}

	return BlockList(threadID, nextOffset, limit, dots)
}

func BlockMeta(blockID string) error {
	_, res, err := getBlockMeta(blockID)
	if err != nil {
		return err
	}
	output(res)
	return nil
}

func getBlockMeta(blockID string) (*pb.Block, string, error) {
	var block pb.Block
	res, err := executeJsonPbCmd(http.MethodGet, "blocks/"+blockID, params{}, &block)
	if err != nil {
		return nil, "", err
	}
	return &block, res, nil
}

// Adds new block to the thread to indicate that this block should be ignored, essentially removing the block
func BlockIgnore(blockID string) error {
	res, err := executeJsonCmd(http.MethodDelete, "blocks/"+blockID, params{}, nil)
	if err != nil {
		return err
	}
	output(res)
	return nil
}
