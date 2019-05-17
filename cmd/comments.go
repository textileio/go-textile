package cmd

import (
	"net/http"
)


func CommentAdd(blockID string, commentBody string) error {
	res, err := executeJsonCmd(http.MethodGet, "blocks/"+blockID+"/comments", params{args: []string{commentBody}}, nil)
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

func CommentGet(commentID string) error {
	res, err := executeJsonCmd(http.MethodGet, "blocks/"+commentID+"/comment", params{}, nil)
	if err != nil {
		return err
	}
	output(res)
	return nil
}


func CommentIgnore(commentID string) error {
	return BlockRemove(commentID)
}
