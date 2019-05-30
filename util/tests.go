package util

import (
	"bufio"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/textileio/go-textile/pb"
)

var (
	TestContact = &pb.Contact{
		Address: "address1",
		Name:    "joe",
		Avatar:  "Qm123",
		Peers: []*pb.Peer{
			{
				Id:      "abcde",
				Address: "address1",
				Name:    "joe",
				Avatar:  "Qm123",
				Inboxes: []*pb.Cafe{{
					Peer:     "peer",
					Address:  "address",
					Api:      "v0",
					Protocol: "/textile/cafe/1.0.0",
					Node:     "v1.0.0",
					Url:      "https://mycafe.com",
				}},
			},
		},
	}
)

const (
	TestLogSchema = `
	{
	  "pin": true,
	  "mill": "/json",
	  "json_schema": {
		"$schema": "http://json-schema.org/draft-04/schema#",
		"$ref": "#/definitions/Log",
		"definitions": {
		  "Log": {
			"required": [
			  "priority",
			  "timestamp",
			  "hostname",
			  "application",
			  "pid",
			  "message"
			],
			"properties": {
			  "application": {
				"type": "string"
			  },
			  "hostname": {
				"type": "string"
			  },
			  "message": {
				"type": "string"
			  },
			  "pid": {
				"type": "integer"
			  },
			  "priority": {
				"type": "integer"
			  },
			  "timestamp": {
				"type": "string"
			  }
			},
			"additionalProperties": false,
			"type": "object"
		  }
		}
	  }
	}`
)

func TestURL(t *testing.T, addr string) {
	req, err := http.NewRequest(http.MethodGet, addr, nil)
	assert.NoError(t, err)

	// @NOTE cannot test cors:
	// 1. CORS server enforcement is optional
	// 2. it is only mandatory for the client to enforce
	// 3. Could not find any go server implementations that offer CORS enforcement
	// 4. Could not find any go client implementations that offer CORS enforcement
	// 5. All that exists in go, is the server sending the headers, for client enforcement

	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)

	t.Log("\nBODY:")
	defer func() {
		err := resp.Body.Close()
		assert.NoError(t, err)
	}()
	scanner := bufio.NewScanner(resp.Body)
	scanner.Split(bufio.ScanBytes)
	for scanner.Scan() {
		t.Log(scanner.Text())
	}

	t.Log("\nHEADERS:")
	for k, v := range resp.Header {
		t.Logf("%s: %s\n", k, v)
	}

	t.Logf("\nSTATUS: %d", resp.StatusCode)
	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("bad status: got %v want %v", resp.StatusCode, http.StatusNoContent)
	}
}
