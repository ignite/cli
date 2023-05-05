package cosmosclient

import (
	"encoding/json"
	"fmt"
	"time"

	ctypes "github.com/cometbft/cometbft/rpc/core/types"
)

// TX defines a block transaction.
type TX struct {
	// BlockTime returns the time of the block that contains the transaction.
	BlockTime time.Time

	// Raw contains the transaction as returned by the Tendermint API.
	Raw *ctypes.ResultTx
}

// GetEvents returns the transaction events.
func (t TX) GetEvents() (events []TXEvent, err error) {
	for _, e := range t.Raw.TxResult.Events {
		evt := TXEvent{Type: e.Type}

		for _, a := range e.Attributes {
			// Make sure that the attribute value is a valid JSON encoded string.
			// Tendermint event attribute values contain JSON encoded values without quotes
			// so string values need to be encoded to be quoted and saved as valid JSONB.
			v, err := formatAttributeValue([]byte(a.Value))
			if err != nil {
				return nil, fmt.Errorf("error encoding event attr '%s.%s': %w", e.Type, a.Key, err)
			}

			evt.Attributes = append(evt.Attributes, TXEventAttribute{
				Key:   a.Key,
				Value: v,
			})
		}

		events = append(events, evt)
	}

	return events, nil
}

// TXEvent defines a transaction event.
type TXEvent struct {
	Type       string             `json:"type"`
	Attributes []TXEventAttribute `json:"attributes"`
}

// TXEventAttribute defines a transaction event attribute.
type TXEventAttribute struct {
	Key   string `json:"key"`
	Value []byte `json:"value"`
}

func formatAttributeValue(v []byte) ([]byte, error) {
	if json.Valid(v) {
		return v, nil
	}

	// Encode all string or invalid values
	return json.Marshal(string(v))
}
