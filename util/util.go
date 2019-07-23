package util

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
)

func UnmarshalString(body io.ReadCloser) (string, error) {
	data, err := ioutil.ReadAll(body)
	if err != nil {
		return "", err
	}
	return TrimQuotes(string(data)), nil
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

func EqualStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func ProtoTime(ts *timestamp.Timestamp) time.Time {
	return time.Unix(ts.Seconds, int64(ts.Nanos))
}

func ProtoNanos(ts *timestamp.Timestamp) int64 {
	if ts == nil {
		ts = ptypes.TimestampNow()
	}
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

func ProtoTsIsNewer(ts1 *timestamp.Timestamp, ts2 *timestamp.Timestamp) bool {
	return ProtoNanos(ts1) > ProtoNanos(ts2)
}

func TrimQuotes(s string) string {
	if len(s) > 0 && s[0] == '"' {
		s = s[1:]
	}
	if len(s) > 0 && s[len(s)-1] == '"' {
		s = s[:len(s)-1]
	}
	return s
}

func WriteFileByPath(name string, data []byte) error {
	dir := filepath.Dir(name)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return ioutil.WriteFile(name, data, 0644)
}

func Mkdirp(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}
