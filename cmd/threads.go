package cmd

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/golang/protobuf/ptypes"
	"github.com/mitchellh/go-homedir"
	"github.com/textileio/go-textile/pb"
	"github.com/textileio/go-textile/schema/textile"
)

func ThreadAdd(name string, key string, tipe string, sharing string, whitelist []string, schema string, schemaFile string, blob bool, cameraRoll bool, media bool) error {
	var body []byte
	if schema == "" {
		if schemaFile != "" {
			path, err := homedir.Expand(string(schemaFile))
			if err != nil {
				return err
			}

			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			body, err = ioutil.ReadAll(file)
			if err != nil {
				return err
			}
		} else if blob {
			body = []byte(textile.Blob)
		} else if cameraRoll {
			body = []byte(textile.CameraRoll)
		} else if media {
			body = []byte(textile.Media)
		}
	}

	if body != nil {
		var schemaf pb.FileIndex
		if _, err := executeJsonPbCmd(http.MethodPost, "mills/schema", params{
			payload: bytes.NewReader(body),
			ctype:   "application/json",
		}, &schemaf); err != nil {
			return err
		}
		schema = schemaf.Hash
	}

	res, err := executeJsonCmd(http.MethodPost, "threads", params{
		args: []string{name},
		opts: map[string]string{
			"key":       key,
			"type":      tipe,
			"sharing":   sharing,
			"whitelist": strings.Join(whitelist, ","),
			"schema":    schema,
		},
	}, nil)
	if err != nil {
		return err
	}
	output(res)
	return nil
}

func ThreadList() error {
	res, err := executeJsonCmd(http.MethodGet, "threads", params{}, nil)
	if err != nil {
		return err
	}
	output(res)
	return nil
}

func ThreadGet(threadID string) error {
	res, err := executeJsonCmd(http.MethodGet, "threads/"+threadID, params{}, nil)
	if err != nil {
		return err
	}
	output(res)
	return nil
}

func ThreadPeer(threadID string) error {
	res, err := executeJsonCmd(http.MethodGet, "threads/"+threadID+"/peers", params{}, nil)
	if err != nil {
		return err
	}
	output(res)
	return nil
}

func ThreadRename(name string, threadID string) error {
	res, err := executeStringCmd(http.MethodPost, "threads/"+threadID+"/name", params{args: []string{name}})
	if err != nil {
		return err
	}
	output(res)
	return nil
}

func ThreadAbandon(threadID string) error {
	res, err := executeStringCmd(http.MethodDelete, "threads/"+threadID, params{})
	if err != nil {
		return err
	}
	output(res)
	return nil
}

func ThreadSnapshotCreate() error {
	res, err := createThreadSnapshot()
	if err != nil {
		return err
	}
	output(res)
	return nil
}

func createThreadSnapshot() (string, error) {
	return executeStringCmd(http.MethodPost, "snapshots", params{})
}

func ThreadSnapshotSearch(wait int) error {
	handleSearchStream("snapshots/search", params{
		opts: map[string]string{
			"wait": strconv.Itoa(wait),
		},
	})
	return nil
}

func ThreadSnapshotApply(id string, wait int) error {
	results := handleSearchStream("snapshots/search", params{
		opts: map[string]string{
			"wait": strconv.Itoa(wait),
		},
	})

	var result *pb.QueryResult
	for _, r := range results {
		if r.Id == id {
			result = &r
		}
	}

	if result == nil {
		output("Could not find snapshot with ID: " + id)
		return nil
	}

	if err := applyThreadSnapshot(result); err != nil {
		return err
	}

	return nil
}

func applyThreadSnapshot(result *pb.QueryResult) error {
	snap := new(pb.Thread)
	if err := ptypes.UnmarshalAny(result.Value, snap); err != nil {
		return err
	}
	data, err := pbMarshaler.MarshalToString(result.Value)
	if err != nil {
		return err
	}

	res, err := executeStringCmd(http.MethodPut, "threads/"+snap.Id, params{
		payload: strings.NewReader(data),
		ctype:   "application/json",
	})
	if err != nil {
		return err
	}
	if res == "" {
		output("applied " + result.Id)
	} else {
		output("error applying " + result.Id + ": " + res)
	}
	return nil
}
