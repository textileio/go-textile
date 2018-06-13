package util

import (
	"bufio"
	"github.com/op/go-logging"
	"os"
)

type StdOutLogger struct {
	log    *logging.Logger
	stdout *os.File
	r      *os.File
	w      *os.File
}

func NewStdOutLogger(log *logging.Logger) (*StdOutLogger, error) {
	r, w, err := os.Pipe()
	if err != nil {
		return nil, err
	}
	return &StdOutLogger{log: log, stdout: os.Stdout, r: r, w: w}, nil
}

func (s *StdOutLogger) Start() {
	os.Stdout = s.w
	go func() {
		scanner := bufio.NewScanner(s.r)
		for scanner.Scan() {
			s.log.Info(scanner.Text())
		}
	}()
}

func (s *StdOutLogger) Stop() {
	s.w.Close()
	os.Stdout = s.stdout
}
