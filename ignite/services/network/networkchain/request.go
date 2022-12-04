package networkchain

import (
	"fmt"
	"github.com/pkg/errors"

	cosmosgenesis "github.com/ignite/cli/ignite/pkg/cosmosutil/genesis"
	"github.com/ignite/cli/ignite/services/network/networktypes"
)

// CheckRequestChangeParam builds the genesis for the chain from the launch approved requests
func (c Chain) CheckRequestChangeParam(
	module, param string, value []byte,
) error {
	genesisPath, err := c.chain.GenesisPath()
	if err != nil {
		return errors.Wrap(err, "genesis of the blockchain can't be read")
	}

	genesis, err := cosmosgenesis.FromPath(genesisPath)
	if err != nil {
		return errors.Wrap(err, "genesis of the blockchain can't be parsed")
	}

	pc := []networktypes.ParamChange{
		{
			Module: module,
			Param: param,
			Value: value,
		},
	}

	if err := applyParamChanges(genesis, pc); err != nil {
		return fmt.Errorf("error applying param changes to genesis: %w", err)
	}

	return nil
}
