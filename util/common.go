package util

import (
	"errors"
	"fmt"
	"strings"
)

func BuildExternalInviteLink(id string, key string, name string) string {
	return fmt.Sprintf("https://www.textile.io/clients/#invite=%s&key=%s&name=%s", id, key, name)
}

func ParseExternalInviteLink(link string) (id string, key string, name string, err error) {
	parts := strings.Split(link, "#")
	if len(parts) == 1 {
		err = errors.New("missing anchor")
		return
	}
	parts = strings.Split(parts[1], "&")
	if len(parts) == 1 {
		err = errors.New("missing params")
		return
	}
	for _, part := range parts {
		kv := strings.Split(part, "=")
		if len(kv) != 2 {
			continue
		}
		switch kv[0] {
		case "invite":
			id = kv[1]
		case "key":
			key = kv[1]
		case "name":
			name = kv[1]
		}
	}
	if id == "" {
		err = errors.New("missing invite id")
		return
	}
	if key == "" {
		err = errors.New("missing invite key")
		return
	}
	if name == "" {
		err = errors.New("missing invite name")
		return
	}
	return
}
