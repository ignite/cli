package networktypes

import (
	"github.com/pkg/errors"

	launchtypes "github.com/tendermint/spn/x/launch/types"
)

// errInvalidRequest is an error returned in methods manipulating requests when they are invalid
var errInvalidRequest = errors.New("request is invalid")

// GenesisInformation represents all information for a chain to construct the genesis.
type GenesisInformation struct {
	GenesisAccounts   map[string]GenesisAccount
	VestingAccounts   map[string]VestingAccount
	GenesisValidators map[string]GenesisValidator
}

// GenesisAccount represents an account with initial coin allocation for the chain for the chain genesis
type GenesisAccount struct {
	Address string
	Coins   string
}

// VestingAccount represents a vesting account with initial coin allocation  and vesting option for the chain genesis
// VestingAccount supports currently only delayed vesting option
type VestingAccount struct {
	Address         string
	StartingBalance string
	Vesting         string
	EndTime         int64
}

// GenesisValidator represents a genesis validator associated with a gentx in the chain genesis
type GenesisValidator struct {
	Address string
	Gentx   []byte
	Peer    string
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
		Address:         acc.Address,
		StartingBalance: acc.StartingBalance.String(),
		Vesting:         delayedVesting.Vesting.String(),
		EndTime:         delayedVesting.EndTime,
	}, nil
}

// ToGenesisValidator converts genesis validator from SPN
func ToGenesisValidator(val launchtypes.GenesisValidator) GenesisValidator {
	return GenesisValidator{
		Address: val.Address,
		Gentx:   val.GenTx,
		Peer:    val.Peer,
	}
}

// NewGenesisInformation initializes a new GenesisInformation
func NewGenesisInformation(
	genAccs []GenesisAccount,
	vestingAccs []VestingAccount,
	genVals []GenesisValidator,
) (gi GenesisInformation) {
	// convert account arrays into maps
	for _, genAcc := range genAccs {
		gi.GenesisAccounts[genAcc.Address] = genAcc
	}
	for _, vestingAcc := range vestingAccs {
		gi.VestingAccounts[vestingAcc.Address] = vestingAcc
	}
	for _, genVal := range genVals {
		gi.GenesisValidators[genVal.Address] = genVal
	}
	return gi
}

// GetGenesisAccounts converts into array and returns genesis accounts
func (gi GenesisInformation) GetGenesisAccounts() (accs []GenesisAccount) {
	for _, genAcc := range gi.GenesisAccounts {
		accs = append(accs, genAcc)
	}
	return accs
}

// GetVestingAccounts converts into array and returns vesting accounts
func (gi GenesisInformation) GetVestingAccounts() (accs []VestingAccount) {
	for _, vestingAcc := range gi.VestingAccounts {
		accs = append(accs, vestingAcc)
	}
	return accs
}

// GetGenesisValidators converts into array and returns genesis validators
func (gi GenesisInformation) GetGenesisValidators() (vals []GenesisValidator) {
	for _, genVal := range gi.GenesisValidators {
		vals = append(vals, genVal)
	}
	return vals
}

// ApplyRequest applies to the genesisInformation the changes implied by the approval of a request
func (gi GenesisInformation) ApplyRequest(request launchtypes.Request) (GenesisInformation, error) {
	switch requestContent := request.Content.Content.(type) {
	case *launchtypes.RequestContent_GenesisAccount:
		// new genesis account in the genesis
		ga := ToGenesisAccount(*requestContent.GenesisAccount)
		if _, ok := gi.GenesisAccounts[ga.Address]; ok {
			return gi, errors.Wrapf(errInvalidRequest, "genesis account %s already in genesis", ga.Address)
		}
		gi.GenesisAccounts[ga.Address] = ga

	case *launchtypes.RequestContent_VestingAccount:
		// new vesting account in the genesis
		va, err := ToVestingAccount(*requestContent.VestingAccount)
		if err != nil {
			// we don't treat this error as errInvalidRequests
			// because it can occur if we don't support this format of vesting account
			// but the request is still correct
			return gi, err
		}

		if _, ok := gi.VestingAccounts[va.Address]; ok {
			return gi, errors.Wrapf(errInvalidRequest, "vesting account %s already in genesis", va.Address)
		}
		gi.VestingAccounts[va.Address] = va

	case *launchtypes.RequestContent_AccountRemoval:
		// account removed from the genesis
		ar := requestContent.AccountRemoval
		_, genExist := gi.GenesisAccounts[ar.Address]
		_, vestingExist := gi.VestingAccounts[ar.Address]
		if !genExist && !vestingExist {
			return gi, errors.Wrapf(errInvalidRequest, "account %s can't be removed because it doesn't exist", ar.Address)
		}
		delete(gi.GenesisAccounts, ar.Address)
		delete(gi.VestingAccounts, ar.Address)

	case *launchtypes.RequestContent_GenesisValidator:
		// new genesis validator in the genesis
		gv := ToGenesisValidator(*requestContent.GenesisValidator)
		if _, ok := gi.GenesisValidators[gv.Address]; ok {
			return gi, errors.Wrapf(errInvalidRequest, "genesis validator %s already in genesis", gv.Address)
		}
		gi.GenesisValidators[gv.Address] = gv

	case *launchtypes.RequestContent_ValidatorRemoval:
		// validator removed from the genesis
		vr := requestContent.ValidatorRemoval
		if _, ok := gi.GenesisValidators[vr.ValAddress]; !ok {
			return gi, errors.Wrapf(errInvalidRequest, "genesis validator %s can't be removed because it doesn't exist", vr.ValAddress)
		}
	}

	return gi, nil
}
