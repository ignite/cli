package tendermintlogger

import tmlog "github.com/cometbft/cometbft/libs/log"

type DiscardLogger struct{}

func (l DiscardLogger) Debug(msg string, keyvals ...interface{}) {}
func (l DiscardLogger) Info(msg string, keyvals ...interface{})  {}
func (l DiscardLogger) Error(msg string, keyvals ...interface{}) {}
func (l DiscardLogger) With(keyvals ...interface{}) tmlog.Logger { return l }
