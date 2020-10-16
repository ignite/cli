package chain

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	starportsecretconf "github.com/tendermint/starport/starport/services/chain/secretconf"
)

type relayerInfo struct {
	ChainID    string
	Mnemonic   string
	RPCAddress string
}

// RelayerInfo initializes or updates relayer setup for the chain itself and returns
// a meta info to share with other chains so they can connect.
// TODO only stargate
func (s *Chain) RelayerInfo() (base64Info string, err error) {
	sconf, err := starportsecretconf.Open(s.app.Path)
	if err != nil {
		return "", err
	}
	relayerAcc, found := sconf.SelfRelayerAccount(s.app.n())
	if !found {
		if err := sconf.SetSelfRelayerAccount(s.app.n()); err != nil {
			return "", err
		}
		relayerAcc, _ = sconf.SelfRelayerAccount(s.app.n())
		if err := starportsecretconf.Save(s.app.Path, sconf); err != nil {
			return "", err
		}
	}
	rpcAddress, err := s.rpcAddress()
	if err != nil {
		return "", err
	}
	info := relayerInfo{
		ChainID:    s.app.n(),
		Mnemonic:   relayerAcc.Mnemonic,
		RPCAddress: rpcAddress,
	}
	data, err := json.Marshal(info)
	if err != nil {
		return "", err
	}
	return base64.RawStdEncoding.EncodeToString(data), nil
}

func (s *Chain) RelayerAdd(base64Info string) error {
	data, err := base64.RawStdEncoding.DecodeString(base64Info)
	if err != nil {
		return err
	}
	var info relayerInfo
	if err := json.Unmarshal(data, &info); err != nil {
		return err
	}
	sconf, err := starportsecretconf.Open(s.app.Path)
	if err != nil {
		return err
	}
	sconf.UpsertRelayerAccount(starportsecretconf.RelayerAccount{
		ID:         info.ChainID,
		Mnemonic:   info.Mnemonic,
		RPCAddress: info.RPCAddress,
	})
	if err := starportsecretconf.Save(s.app.Path, sconf); err != nil {
		return err
	}
	fmt.Fprint(s.stdLog(logStarport).out, "\n💫  Chain added\n")
	return nil
}
