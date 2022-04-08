package tendermintlogger

import tmlog "github.com/tendermint/tendermint/libs/log"

var _ tmlog.Logger = (*DiscardLogger)(nil)
