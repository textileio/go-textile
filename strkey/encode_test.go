package strkey_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	. "github.com/textileio/go-textile/strkey"
)

func TestEncode(t *testing.T) {
	cases := []struct {
		Name        string
		VersionByte VersionByte
		Payload     []byte
		Expected    string
	}{
		{
			Name:        "AccountID",
			VersionByte: VersionByteAccountID,
			Payload: []byte{
				0x69, 0xa8, 0xc4, 0xcb, 0xb9, 0xf6, 0x4e, 0x8a,
				0x07, 0x98, 0xf6, 0xe1, 0xac, 0x65, 0xd0, 0x6c,
				0x31, 0x62, 0x92, 0x90, 0x56, 0xbc, 0xf4, 0xcd,
				0xb7, 0xd3, 0x73, 0x8d, 0x18, 0x55, 0xf3, 0x63,
			},
			Expected: "P6NKHguuQnGraNM58WZTh2tPmRAnQQn1YzHQZYPmRt8WABDF",
		},
		{
			Name:        "Seed",
			VersionByte: VersionByteSeed,
			Payload: []byte{
				0x69, 0xa8, 0xc4, 0xcb, 0xb9, 0xf6, 0x4e, 0x8a,
				0x07, 0x98, 0xf6, 0xe1, 0xac, 0x65, 0xd0, 0x6c,
				0x31, 0x62, 0x92, 0x90, 0x56, 0xbc, 0xf4, 0xcd,
				0xb7, 0xd3, 0x73, 0x8d, 0x18, 0x55, 0xf3, 0x63,
			},
			Expected: "SV8k5RKcUg1ZtN6qvcUGXfvqjv2nGeBhxy7sUnG1AxM5jB23",
		},
	}

	for _, kase := range cases {
		actual, err := Encode(kase.VersionByte, kase.Payload)
		if assert.NoError(t, err, "An error occured in case %s", kase.Name) {
			assert.Equal(t, kase.Expected, actual, "Output mismatch in case %s", kase.Name)
		}
	}

	// test bad version byte
	_, err := Encode(VersionByte(2), cases[0].Payload)
	assert.Error(t, err)
}
