package cosmosutil

import (
	"strings"

	launchtypes "github.com/tendermint/spn/x/launch/types"
)

// VerifyPeerFormat checks if the peer address format is valid
func VerifyPeerFormat(peer launchtypes.Peer) bool {
	// Check the format of the peer
	switch conn := peer.Connection.(type) {
	case *launchtypes.Peer_TcpAddress:
		nodeHost := strings.Split(conn.TcpAddress, "@")
		if len(nodeHost) != 2 ||
			len(nodeHost[0]) == 0 ||
			len(nodeHost[1]) == 0 {
			return false
		}
	case *launchtypes.Peer_HttpTunnel:
		if conn.HttpTunnel.Address == "" {
			return false
		}
	default:
		return false
	}
	return true
}
