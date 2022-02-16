package networktypes

import (
	"github.com/pkg/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	launchtypes "github.com/tendermint/spn/x/launch/types"
)

// GenesisInformation represents all information for a chain to construct the genesis.
// This structure indexes accounts and validators by their address for better performance
type GenesisInformation struct {
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
	gi.GenesisAccounts = make([]GenesisAccount, len(genAccs))
	gi.VestingAccounts = make([]VestingAccount, len(vestingAccs))
	gi.GenesisValidators = make([]GenesisValidator, len(genVals))
	copy(gi.GenesisAccounts, genAccs)
	copy(gi.VestingAccounts, vestingAccs)
	copy(gi.GenesisValidators, genVals)
	return gi
}

// GetGenesisAccounts converts into array and returns genesis accounts
func (gi GenesisInformation) GetGenesisAccounts() (accs []GenesisAccount) {
	accs = make([]GenesisAccount, len(gi.GenesisAccounts))
	copy(accs, gi.GenesisAccounts)
	return accs
}

// GetVestingAccounts converts into array and returns vesting accounts
func (gi GenesisInformation) GetVestingAccounts() (accs []VestingAccount) {
	accs = make([]VestingAccount, len(gi.VestingAccounts))
	copy(accs, gi.VestingAccounts)
	return accs
}

// GetGenesisValidators converts into array and returns genesis validators
func (gi GenesisInformation) GetGenesisValidators() (vals []GenesisValidator) {
	vals = make([]GenesisValidator, len(gi.GenesisValidators))
	copy(vals, gi.GenesisValidators)
	return vals
}

func (gi GenesisInformation) containsGenesisAccount(address string) bool {
	for _, account := range gi.GenesisAccounts {
		if account.Address == address {
			return true
		}
	}
	return false
}
func (gi GenesisInformation) containsVestingAccount(address string) bool {
	for _, account := range gi.VestingAccounts {
		if account.Address == address {
			return true
		}
	}
	return false
}
func (gi GenesisInformation) containsGenesisValidator(address string) bool {
	for _, account := range gi.GenesisValidators {
		if account.Address == address {
			return true
		}
	}
	return false
}

func (gi *GenesisInformation) addGenesisAccount(acc GenesisAccount) {
	gi.GenesisAccounts = append(gi.GenesisAccounts, acc)
}

func (gi *GenesisInformation) addVestingAccount(acc VestingAccount) {
	gi.VestingAccounts = append(gi.VestingAccounts, acc)
}

func (gi *GenesisInformation) addGenesisValidator(val GenesisValidator) {
	gi.GenesisValidators = append(gi.GenesisValidators, val)
}

func (gi *GenesisInformation) removeGenesisAccount(address string) {
	for i, account := range gi.GenesisAccounts {
		if account.Address == address {
			gi.GenesisAccounts = append(gi.GenesisAccounts[:i], gi.GenesisAccounts[i+1:]...)
		}
	}
}

func (gi *GenesisInformation) removeVestingAccount(address string) {
	for i, account := range gi.VestingAccounts {
		if account.Address == address {
			gi.VestingAccounts = append(gi.VestingAccounts[:i], gi.VestingAccounts[i+1:]...)
		}
	}
}

func (gi *GenesisInformation) removeGenesisValidator(address string) {
	for i, account := range gi.GenesisValidators {
		if account.Address == address {
			gi.GenesisValidators = append(gi.GenesisValidators[:i], gi.GenesisValidators[i+1:]...)
		}
	}
}

// ApplyRequest applies to the genesisInformation the changes implied by the approval of a request
func (gi GenesisInformation) ApplyRequest(request launchtypes.Request) (GenesisInformation, error) {
	switch requestContent := request.Content.Content.(type) {
	case *launchtypes.RequestContent_GenesisAccount:
		// new genesis account in the genesis
		ga := ToGenesisAccount(*requestContent.GenesisAccount)
		genExist := gi.containsGenesisAccount(ga.Address)
		vestingExist := gi.containsVestingAccount(ga.Address)
		if genExist || vestingExist {
			return gi, NewWrappedErrInvalidRequest(request.RequestID, "genesis account already in genesis")
		}
		gi.addGenesisAccount(ga)

	case *launchtypes.RequestContent_VestingAccount:
		// new vesting account in the genesis
		va, err := ToVestingAccount(*requestContent.VestingAccount)
		if err != nil {
			// we don't treat this error as errInvalidRequests
			// because it can occur if we don't support this format of vesting account
			// but the request is still correct
			return gi, err
		}

		genExist := gi.containsGenesisAccount(va.Address)
		vestingExist := gi.containsVestingAccount(va.Address)
		if genExist || vestingExist {
			return gi, NewWrappedErrInvalidRequest(request.RequestID, "vesting account already in genesis")
		}
		gi.addVestingAccount(va)

	case *launchtypes.RequestContent_AccountRemoval:
		// account removed from the genesis
		ar := requestContent.AccountRemoval
		genExist := gi.containsGenesisAccount(ar.Address)
		vestingExist := gi.containsVestingAccount(ar.Address)
		if !genExist && !vestingExist {
			return gi, NewWrappedErrInvalidRequest(request.RequestID, "account can't be removed because it doesn't exist")
		}
		gi.removeGenesisAccount(ar.Address)
		gi.removeVestingAccount(ar.Address)

	case *launchtypes.RequestContent_GenesisValidator:
		// new genesis validator in the genesis
		gv := ToGenesisValidator(*requestContent.GenesisValidator)
		if gi.containsGenesisValidator(gv.Address) {
			return gi, NewWrappedErrInvalidRequest(request.RequestID, "genesis validator already in genesis")
		}
		gi.addGenesisValidator(gv)

	case *launchtypes.RequestContent_ValidatorRemoval:
		// validator removed from the genesis
		vr := requestContent.ValidatorRemoval
		if gi.containsGenesisValidator(vr.ValAddress) {
			return gi, NewWrappedErrInvalidRequest(request.RequestID, "genesis validator can't be removed because it doesn't exist")
		}
	}

	return gi, nil
}
