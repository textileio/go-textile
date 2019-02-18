package util

import (
	"io"
	"io/ioutil"
	"strings"

	"github.com/golang/protobuf/ptypes/timestamp"
)

func UnmarshalString(body io.ReadCloser) (string, error) {
	data, err := ioutil.ReadAll(body)
	if err != nil {
		return "", err
	}
	return trimQuotes(string(data)), nil
}

func SplitString(in string, sep string) []string {
	list := make([]string, 0)
	for _, s := range strings.Split(in, sep) {
		t := strings.TrimSpace(s)
		if t != "" {
			list = append(list, t)
		}
	}
	return list
}

func ListContainsString(list []string, i string) bool {
	for _, v := range list {
		if v == i {
			return true
		}
	}
	return false
}

func ProtoNanos(ts *timestamp.Timestamp) int64 {
	return int64(ts.Nanos) + ts.Seconds*1e9
}

func ProtoTs(nsec int64) *timestamp.Timestamp {
	n := nsec / 1e9
	sec := n
	nsec -= n * 1e9
	if nsec < 0 {
		nsec += 1e9
		sec--
	}

	return &timestamp.Timestamp{
		Seconds: sec,
		Nanos:   int32(nsec),
	}
}

func trimQuotes(s string) string {
	if len(s) > 0 && s[0] == '"' {
		s = s[1:]
	}
	if len(s) > 0 && s[len(s)-1] == '"' {
		s = s[:len(s)-1]
	}
	return s
}
