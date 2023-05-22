package tendermintlogger

import tmlog "github.com/cometbft/cometbft/libs/log"

var _ tmlog.Logger = (*DiscardLogger)(nil)
