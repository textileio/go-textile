package util

import (
	"bufio"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestURL(t *testing.T, addr string, method string, status int) {
	// prepare the request
	req, err := http.NewRequest(http.MethodGet, addr, nil)
	assert.NoError(t, err)

	// @NOTE cannot test cors:
	// 1. CORS server enforcement is optional
	// 2. it is only mandatory for the client to enforce
	// 3. Could not find any go server implementations that offer CORS enforcement
	// 4. Could not find any go client implementations that offer CORS enforcement
	// 5. All that exists in go, is the server sending the headers, for client enforcement

	// perform the response
	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)

	// body
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

	// headers
	t.Log("\nHEADERS:")
	for k, v := range resp.Header {
		t.Logf("%s: %s\n", k, v)
	}

	// status
	t.Logf("\nSTATUS: %d", resp.StatusCode)
	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("bad status: got %v want %v", resp.StatusCode, http.StatusNoContent)
		return
		}
}