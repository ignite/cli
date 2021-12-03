package cosmosutil

import (
	"strings"
)

// VerifyPeerFormat checks if the peer address format is valid
func VerifyPeerFormat(peer string) bool {
	// Check the format of the peer
	nodeHost := strings.Split(peer, "@")
	if len(nodeHost) != 2 ||
		len(nodeHost[0]) == 0 ||
		len(nodeHost[1]) == 0 {
		return false
	}
	return true
}
