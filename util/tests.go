package util

import (
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
