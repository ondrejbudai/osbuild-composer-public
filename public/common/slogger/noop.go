package slogger

import (
	"github.com/ondrejbudai/osbuild-composer-public/pkg/jobqueue"
)

type noopLogger struct {
}

func NewNoopLogger() jobqueue.SimpleLogger {
	return &noopLogger{}
}

func (s *noopLogger) Info(_ string, _ ...string) {
}

func (s *noopLogger) Error(_ error, _ string, _ ...string) {
}
