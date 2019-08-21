package cmd

import (
	"net/http"
)

func CommentAdd(blockID string, commentBody string) error {
	res, err := executeJsonCmd(http.MethodPost, "blocks/"+blockID+"/comments", params{args: []string{commentBody}}, nil)
	if err != nil {
		return err
	}
	output(res)
	return nil
}

func CommentList(blockID string) error {
	res, err := executeJsonCmd(http.MethodGet, "blocks/"+blockID+"/comments", params{}, nil)
	if err != nil {
		return err
	}
	output(res)
	return nil
}

func CommentGet(blockID string) error {
	res, err := executeJsonCmd(http.MethodGet, "blocks/"+blockID+"/comment", params{}, nil)
	if err != nil {
		return err
	}
	output(res)
	return nil
}

func CommentIgnore(blockID string) error {
	return BlockIgnore(blockID)
}
