package testutil

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

type (
	Gentx struct {
		Body Body `json:"body"`
	}

	Body struct {
		Messages []Message `json:"messages"`
		Memo     string    `json:"memo"`
	}

	Message struct {
		DelegatorAddress string        `json:"delegator_address"`
		ValidatorAddress string        `json:"validator_address"`
		PubKey           MessagePubKey `json:"pubkey"`
		Value            MessageValue  `json:"value"`
	}

	MessageValue struct {
		Denom  string `json:"denom"`
		Amount string `json:"amount"`
	}

	MessagePubKey struct {
		Key string `json:"key"`
	}
)

// // NewGenesis creates easily modifiable gentx object for testing purposes
func NewGentx(address, denom, amount, pubkey, memo string) *Gentx {
	return &Gentx{Body: Body{
		Memo: memo,
		Messages: []Message{
			{
				DelegatorAddress: address,
				PubKey:           MessagePubKey{Key: pubkey},
				Value:            MessageValue{Denom: denom, Amount: amount},
			},
		},
	}}
}

// SaveTo saves gentx json representation to the specified directory and returns full path
func (g *Gentx) SaveTo(t *testing.T, dir string) string {
	encoded, err := json.Marshal(g)
	assert.NoError(t, err)
	savePath := filepath.Join(dir, "gentx0.json")
	err = os.WriteFile(savePath, encoded, 0666)
	assert.NoError(t, err)
	return savePath
}

// JSON returns json representation of the gentx
func (g *Gentx) JSON(t *testing.T) []byte {
	data, err := json.Marshal(g)
	assert.NoError(t, err)
	return data
}
