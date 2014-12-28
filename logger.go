package main

import "github.com/ivpusic/golog"

type DebugLogWriter struct {
	*golog.Logger
}

func (l *DebugLogWriter) Write(p []byte) (n int, err error) {
	l.Debug(string(p))
	return len(p), nil
}
