package main

import (
	"github.com/op/go-logging"
)

var stdoutLogFormat = logging.MustStringFormatter(
	`%{color:reset}%{color}%{time:15:04:05.000} [%{shortfunc}] [%{level}] %{message}`,
)

var logger logging.Backend

type MobileConfig struct {

	// Path for the node's data directory
	RepoPath string
}
