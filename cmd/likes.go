package cmd

import (
	"net/http"
)

func LikeAdd(blockID string) error {
	res, err := executeJsonCmd(http.MethodPost, "blocks/"+blockID+"/likes", params{}, nil)
	if err != nil {
		return err
	}
	output(res)
	return nil
}

func LikeList(blockID string) error {
	res, err := executeJsonCmd(http.MethodGet, "blocks/"+blockID+"/likes", params{}, nil)
	if err != nil {
		return err
	}
	output(res)
	return nil
}

func LikeGet(likeID string) error {
	res, err := executeJsonCmd(http.MethodGet, "blocks/"+likeID+"/like", params{}, nil)
	if err != nil {
		return err
	}
	output(res)
	return nil
}

func LikeIgnore(likeID string) error {
	return BlockIgnore(likeID)
}
