package plugintypes

import (
	"context"

	"github.com/hashicorp/go-plugin"
	"github.com/spf13/cobra"
)

type InitArgs struct {
	ctx context.Context
}

type ExecArgs struct {
	Cmd  *cobra.Command
	Args []string
}

var HandshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "BASIC_PLUGIN",
	MagicCookieValue: "hello",
}
