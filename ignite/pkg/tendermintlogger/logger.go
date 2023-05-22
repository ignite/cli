package tendermintlogger

import tmlog "github.com/cometbft/cometbft/libs/log"

type DiscardLogger struct{}

func (l DiscardLogger) Debug(string, ...interface{})     {}
func (l DiscardLogger) Info(string, ...interface{})      {}
func (l DiscardLogger) Error(string, ...interface{})     {}
func (l DiscardLogger) With(...interface{}) tmlog.Logger { return l }
