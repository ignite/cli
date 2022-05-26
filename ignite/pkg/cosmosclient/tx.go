package cosmosclient

import (
	"encoding/json"
	"time"

	ctypes "github.com/tendermint/tendermint/rpc/core/types"
)

// TX defines a block transaction.
type TX struct {
	// BlockTime returns the time of the block that contains the transaction.
	BlockTime time.Time

	// Raw contains the transaction as returned by the Tendermint API.
	Raw *ctypes.ResultTx
}

// TXEvent defines a transaction event.
type TXEvent struct {
	Type       string             `json:"type"`
	Attributes []TXEventAttribute `json:"attributes"`
}

// TXEventAttribute defines a transaction event attribute.
type TXEventAttribute struct {
	Key   string `json:"key"`
	Value any    `json:"value"`
}

// UnmarshallEvents parses the JSON encoded transactions events.
func UnmarshallEvents(tx TX) ([]TXEvent, error) {
	// The transaction's event log contains a list where each item is an object
	// with a single "events" property which in turn contains the list of events
	var log []struct {
		Events []TXEvent `json:"events"`
	}

	raw := tx.Raw.TxResult.GetLog()
	if err := json.Unmarshal([]byte(raw), &log); err != nil {
		return nil, err
	}

	if len(log) > 0 {
		return log[0].Events, nil
	}

	return nil, nil
}
