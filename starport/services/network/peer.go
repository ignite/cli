package network

import (
	"fmt"

	launchtypes "github.com/tendermint/spn/x/launch/types"
)

func PeerAddress(peer launchtypes.Peer) (string, error) {
	var peerAddr string
	switch conn := peer.Connection.(type) {
	case *launchtypes.Peer_TcpAddress:
		peerAddr = fmt.Sprintf("%s@%s", peer.Id, conn.TcpAddress)
	case *launchtypes.Peer_HttpTunnel:
		peerAddr = fmt.Sprintf("%s@%s", peer.Id, conn.HttpTunnel.Address)
	default:
		return peerAddr, fmt.Errorf("invalid peer type")
	}
	return peerAddr, nil
}
