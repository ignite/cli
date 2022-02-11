package cosmosutil

import (
	"net"
	"strings"

	launchtypes "github.com/tendermint/spn/x/launch/types"
	"github.com/tendermint/starport/starport/pkg/xurl"
)

// VerifyPeerFormat checks if the peer address format is valid
func VerifyPeerFormat(peer launchtypes.Peer) bool {
	// Check the format of the peer
	switch conn := peer.Connection.(type) {
	case *launchtypes.Peer_TcpAddress:
		nodeHost := strings.Split(conn.TcpAddress, ":")
		if len(nodeHost) != 2 ||
			net.ParseIP(nodeHost[0]) == nil ||
			len(nodeHost[1]) == 0 {
			return false
		}
		return true
	case *launchtypes.Peer_HttpTunnel:
		return xurl.IsHTTP(conn.HttpTunnel.Address) && conn.HttpTunnel.Name == "chisel"
	default:
		return false
	}
}
