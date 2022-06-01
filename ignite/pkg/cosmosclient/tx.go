package cosmosclient

import (
	"encoding/json"
	"fmt"
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

// GetLog returns the event log encoded as JSON.
func (t TX) GetEventLog() []byte {
	return []byte(t.Raw.TxResult.GetLog())
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

// UnmarshallEvents parses JSON encoded transaction event logs.
func UnmarshallEvents(b []byte) ([]TXEvent, error) {
	// The transaction's event log contains a list where each item is an object
	// with a single "events" property which in turn contains the list of events
	var log []struct {
		Events []TXEvent `json:"events"`
	}

	if err := json.Unmarshal(b, &log); err != nil {
		return nil, fmt.Errorf("error decoding transaction events: %w", err)
	}

	if len(log) > 0 {
		return log[0].Events, nil
	}

	return nil, nil
}
