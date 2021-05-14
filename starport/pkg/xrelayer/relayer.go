package xrelayer

import (
	"context"

	tsrelayer "github.com/tendermint/starport/starport/pkg/nodetime/ts-relayer"
)

// Link links all chains that has a path to each other.
// paths are optional and acts as a filter to only link some chains.
// calling Link multiple times for the same paths does not have any side effects.
type LinkStatus struct {
	ID       string `json:"pathName"`
	ErrorMsg string `json:"error"`
}

func Link(ctx context.Context, paths ...string) (linkedPaths, alreadyLinkedPaths []string, failedToLinkPaths []LinkStatus, err error) {
	var reply struct {
		LinkedPaths        []string     `json:"linkedPaths"`
		AlreadyLinkedPaths []string     `json:"alreadyLinkedPaths"`
		FailedToLinkPaths  []LinkStatus `json:"failedToLinkPaths"`
	}
	err = tsrelayer.Call(ctx, "link", []interface{}{paths}, &reply)
	linkedPaths = reply.LinkedPaths
	alreadyLinkedPaths = reply.AlreadyLinkedPaths
	failedToLinkPaths = reply.FailedToLinkPaths
	return
}

// Start relays tx packets for paths until ctx is canceled.
func Start(ctx context.Context, paths ...string) error {
	var reply interface{}
	return tsrelayer.Call(ctx, "start", []interface{}{paths}, &reply)
}

// Path represents a path between two chains.
type Path struct {
	// ID is id of the path.
	ID string `json:"id"`

	// IsLinked indicates that chains of these paths are linked or not.
	IsLinked bool `json:"isLinked"`

	// Src end of the path.
	Src PathEnd `json:"src"`

	// Dst end of the path.
	Dst PathEnd `json:"dst"`
}

// PathEnd represents the chain at one side of a Path.
type PathEnd struct {
	ChannelID string `json:"channelID"`
	ChainID   string `json:"chainID"`
	PortID    string `json:"portID"`
}

// GetPath returns a path by its id.
func GetPath(ctx context.Context, id string) (Path, error) {
	var path Path
	err := tsrelayer.Call(ctx, "getPath", []interface{}{id}, &path)
	return path, err
}

// ListPaths list all the paths.
func ListPaths(ctx context.Context) ([]Path, error) {
	var paths []Path
	err := tsrelayer.Call(ctx, "listPaths", nil, &paths)
	return paths, err
}
