package testutil

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type (
	Gentx struct {
		Body Body `json:"body"`
	}

	Body struct {
		Messages []Message `json:"messages"`
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

func NewGentx(address, denom, amount, pubkey string) *Gentx {
	return &Gentx{Body: Body{
		Messages: []Message{
			{
				DelegatorAddress: address,
				PubKey:           MessagePubKey{Key: pubkey},
				Value:            MessageValue{Denom: denom, Amount: amount},
			},
		},
	}}
}

func (g *Gentx) SaveTo(dir string) (string, error) {
	encoded, err := json.Marshal(g)
	if err != nil {
		return "", err
	}
	savePath := filepath.Join(dir, "gentx0.json")
	return savePath, os.WriteFile(savePath, encoded, 0666)
}
func (g *Gentx) JSON() ([]byte, error) {
	return json.Marshal(g)
}
