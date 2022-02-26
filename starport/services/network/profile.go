package network

import (
	"context"
	"fmt"
	"io/ioutil"

	valtypes "github.com/tendermint/spn/pkg/types"
	profiletypes "github.com/tendermint/spn/x/profile/types"

	"github.com/tendermint/starport/starport/pkg/events"
	"github.com/tendermint/starport/starport/services/network/networktypes"
)

// SetValidatorConsAddress associates a Tendermint consensus address to a specific validator address on SPN
func (n Network) SetValidatorConsAddress(ctx context.Context, validatorKeyPath string) error {
	// TODO: set the correct chainID dynamically
	chainID := "spn-1"

	n.ev.Send(events.New(events.StatusOngoing,
		fmt.Sprintf("Reading the validator key %s", validatorKeyPath)))

	// Read and parse the validator private key file
	valConsKey, valPubKey, err := parseValidatorKey(validatorKeyPath)
	if err != nil {
		return err
	}
	consAddress := valPubKey.GetConsAddress().Bytes()

	n.ev.Send(events.New(events.StatusOngoing,
		fmt.Sprintf("Setting the validator consensus address chain %s",
			valPubKey.GetConsAddress().String()),
	))

	// Get the current consensus key nonce and sign the message
	nonce, err := n.ConsensusKeyNonce(ctx, consAddress)
	if err != nil {
		nonce = 0
	}
	signature, err := valConsKey.Sign(nonce, chainID)
	if err != nil {
		return err
	}

	// Create and broadcast the transaction
	msg := profiletypes.NewMsgSetValidatorConsAddress(
		n.account.Address(networktypes.SPN),
		signature,
		valConsKey.PubKey.Type(),
		chainID,
		nonce,
		valConsKey.PubKey.Bytes(),
	)
	res, err := n.cosmos.BroadcastTx(n.account.Name, msg)
	if err != nil {
		return err
	}

	var consAddrRes profiletypes.MsgSetValidatorConsAddressResponse
	if err := res.Decode(&consAddrRes); err != nil {
		return err
	}

	n.ev.Send(events.New(events.StatusDone,
		fmt.Sprintf("Validator consensus address added %s",
			valPubKey.GetConsAddress().String()),
	))
	return nil
}

// ConsensusKeyNonce fetches the consensus key nonce from Starport Network
func (n Network) ConsensusKeyNonce(ctx context.Context, consensusAddress []byte) (uint64, error) {
	n.ev.Send(events.New(events.StatusOngoing, "Fetching consensus key nonce"))
	res, err := profiletypes.NewQueryClient(n.cosmos.Context).
		ConsensusKeyNonce(ctx,
			&profiletypes.QueryGetConsensusKeyNonceRequest{
				ConsensusAddress: consensusAddress,
			},
		)
	if err != nil {
		return 0, err
	}
	return res.ConsensusKeyNonce.Nonce, nil
}

// parseValidatorKey read and parse the validator private key file from path
func parseValidatorKey(validatorKeyPath string) (valtypes.ValidatorKey, valtypes.ValidatorConsPubKey, error) {
	valConsKeyBytes, err := ioutil.ReadFile(validatorKeyPath)
	if err != nil {
		return valtypes.ValidatorKey{}, valtypes.ValidatorConsPubKey{}, err
	}
	valConsKey, err := valtypes.LoadValidatorKey(valConsKeyBytes)
	if err != nil {
		return valConsKey, valtypes.ValidatorConsPubKey{}, err
	}

	// Convert to consensus pub key type to fetch the address
	valPubKey, err := valtypes.NewValidatorConsPubKey(valConsKey.PubKey.Bytes(), valConsKey.PubKey.Type())
	return valConsKey, valPubKey, err
}
