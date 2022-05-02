package networktypes

import (
	"github.com/pkg/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	launchtypes "github.com/tendermint/spn/x/launch/types"
)

// GenesisInformation represents all information for a chain to construct the genesis.
// This structure indexes accounts and validators by their address for better performance
type GenesisInformation struct {
	// make sure to use slices for the following because slices are ordered.
	// they later used to create a Genesis so, having them ordered is important to
	// be able to produce a deterministic Genesis.

	GenesisAccounts   []GenesisAccount
	VestingAccounts   []VestingAccount
	GenesisValidators []GenesisValidator
}

// GenesisAccount represents an account with initial coin allocation for the chain for the chain genesis
type GenesisAccount struct {
	Address string
	Coins   string
}

// VestingAccount represents a vesting account with initial coin allocation  and vesting option for the chain genesis
// VestingAccount supports currently only delayed vesting option
type VestingAccount struct {
	Address      string
	TotalBalance string
	Vesting      string
	EndTime      int64
}

// GenesisValidator represents a genesis validator associated with a gentx in the chain genesis
type GenesisValidator struct {
	Address        string
	Gentx          []byte
	Peer           launchtypes.Peer
	SelfDelegation sdk.Coin
}

// ToGenesisAccount converts genesis account from SPN
func ToGenesisAccount(acc launchtypes.GenesisAccount) GenesisAccount {
	return GenesisAccount{
		Address: acc.Address,
		Coins:   acc.Coins.String(),
	}
}

// ToVestingAccount converts vesting account from SPN
func ToVestingAccount(acc launchtypes.VestingAccount) (VestingAccount, error) {
	delayedVesting := acc.VestingOptions.GetDelayedVesting()
	if delayedVesting == nil {
		return VestingAccount{}, errors.New("only delayed vesting option is supported")
	}

	return VestingAccount{
		Address:      acc.Address,
		TotalBalance: delayedVesting.TotalBalance.String(),
		Vesting:      delayedVesting.Vesting.String(),
		EndTime:      delayedVesting.EndTime,
	}, nil
}

// ToGenesisValidator converts genesis validator from SPN
func ToGenesisValidator(val launchtypes.GenesisValidator) GenesisValidator {
	return GenesisValidator{
		Address:        val.Address,
		Gentx:          val.GenTx,
		Peer:           val.Peer,
		SelfDelegation: val.SelfDelegation,
	}
}

// NewGenesisInformation initializes a new GenesisInformation
func NewGenesisInformation(
	genAccs []GenesisAccount,
	vestingAccs []VestingAccount,
	genVals []GenesisValidator,
) (gi GenesisInformation) {
	return GenesisInformation{
		GenesisAccounts:   genAccs,
		VestingAccounts:   vestingAccs,
		GenesisValidators: genVals,
	}
}

func (gi GenesisInformation) ContainsGenesisAccount(address string) bool {
	for _, account := range gi.GenesisAccounts {
		if account.Address == address {
			return true
		}
	}
	return false
}
func (gi GenesisInformation) ContainsVestingAccount(address string) bool {
	for _, account := range gi.VestingAccounts {
		if account.Address == address {
			return true
		}
	}
	return false
}
func (gi GenesisInformation) ContainsGenesisValidator(address string) bool {
	for _, account := range gi.GenesisValidators {
		if account.Address == address {
			return true
		}
	}
	return false
}

func (gi *GenesisInformation) AddGenesisAccount(acc GenesisAccount) {
	gi.GenesisAccounts = append(gi.GenesisAccounts, acc)
}

func (gi *GenesisInformation) AddVestingAccount(acc VestingAccount) {
	gi.VestingAccounts = append(gi.VestingAccounts, acc)
}

func (gi *GenesisInformation) AddGenesisValidator(val GenesisValidator) {
	gi.GenesisValidators = append(gi.GenesisValidators, val)
}

func (gi *GenesisInformation) RemoveGenesisAccount(address string) {
	for i, account := range gi.GenesisAccounts {
		if account.Address == address {
			gi.GenesisAccounts = append(gi.GenesisAccounts[:i], gi.GenesisAccounts[i+1:]...)
		}
	}
}

func (gi *GenesisInformation) RemoveVestingAccount(address string) {
	for i, account := range gi.VestingAccounts {
		if account.Address == address {
			gi.VestingAccounts = append(gi.VestingAccounts[:i], gi.VestingAccounts[i+1:]...)
		}
	}
}

func (gi *GenesisInformation) RemoveGenesisValidator(address string) {
	for i, account := range gi.GenesisValidators {
		if account.Address == address {
			gi.GenesisValidators = append(gi.GenesisValidators[:i], gi.GenesisValidators[i+1:]...)
		}
	}
}

// ApplyRequest applies to the genesisInformation the changes implied by the approval of a request
func (gi GenesisInformation) ApplyRequest(request Request) (GenesisInformation, error) {
	switch requestContent := request.Content.Content.(type) {
	case *launchtypes.RequestContent_GenesisAccount:
		// new genesis account in the genesis
		ga := ToGenesisAccount(*requestContent.GenesisAccount)
		genExist := gi.ContainsGenesisAccount(ga.Address)
		vestingExist := gi.ContainsVestingAccount(ga.Address)
		if genExist || vestingExist {
			return gi, NewWrappedErrInvalidRequest(request.RequestID, "genesis account already in genesis")
		}
		gi.AddGenesisAccount(ga)

	case *launchtypes.RequestContent_VestingAccount:
		// new vesting account in the genesis
		va, err := ToVestingAccount(*requestContent.VestingAccount)
		if err != nil {
			// we don't treat this error as errInvalidRequests
			// because it can occur if we don't support this format of vesting account
			// but the request is still correct
			return gi, err
		}

		genExist := gi.ContainsGenesisAccount(va.Address)
		vestingExist := gi.ContainsVestingAccount(va.Address)
		if genExist || vestingExist {
			return gi, NewWrappedErrInvalidRequest(request.RequestID, "vesting account already in genesis")
		}
		gi.AddVestingAccount(va)

	case *launchtypes.RequestContent_AccountRemoval:
		// account removed from the genesis
		ar := requestContent.AccountRemoval
		genExist := gi.ContainsGenesisAccount(ar.Address)
		vestingExist := gi.ContainsVestingAccount(ar.Address)
		if !genExist && !vestingExist {
			return gi, NewWrappedErrInvalidRequest(request.RequestID, "account can't be removed because it doesn't exist")
		}
		gi.RemoveGenesisAccount(ar.Address)
		gi.RemoveVestingAccount(ar.Address)

	case *launchtypes.RequestContent_GenesisValidator:
		// new genesis validator in the genesis
		gv := ToGenesisValidator(*requestContent.GenesisValidator)
		if gi.ContainsGenesisValidator(gv.Address) {
			return gi, NewWrappedErrInvalidRequest(request.RequestID, "genesis validator already in genesis")
		}
		gi.AddGenesisValidator(gv)

	case *launchtypes.RequestContent_ValidatorRemoval:
		// validator removed from the genesis
		vr := requestContent.ValidatorRemoval
		if !gi.ContainsGenesisValidator(vr.ValAddress) {
			return gi, NewWrappedErrInvalidRequest(request.RequestID, "genesis validator can't be removed because it doesn't exist")
		}
	}

	return gi, nil
}
