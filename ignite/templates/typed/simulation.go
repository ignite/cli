package typed

import (
	"fmt"
	"strings"

	"github.com/ignite/cli/v29/ignite/pkg/multiformatname"
	"github.com/ignite/cli/v29/ignite/pkg/xast"
)

func ModuleSimulationMsgModify(
	content,
	modulePath,
	moduleName string,
	typeName, msgSigner multiformatname.Name,
	msgs ...string,
) (string, error) {
	if len(msgs) == 0 {
		msgs = append(msgs, "")
	}

	// Import
	content, err := xast.AppendImports(
		content,
		xast.WithLastNamedImport(
			fmt.Sprintf("%[1]vsimulation", moduleName),
			fmt.Sprintf("%[1]v/x/%[2]v/simulation", modulePath, moduleName),
		),
		xast.WithImport("math/rand", 0),
	)
	if err != nil {
		return "", err
	}

	for _, msg := range msgs {
		// simulation operations
		replacementOp := fmt.Sprintf(`
	const (
		opWeightMsg%[1]v%[2]v = "op_weight_msg_%[3]v"
		defaultWeightMsg%[1]v%[2]v int = 100 // TODO: Determine the simulation weight value for your use case
	)

	var weightMsg%[1]v%[2]v int
	simState.AppParams.GetOrGenerate(opWeightMsg%[1]v%[2]v, &weightMsg%[1]v%[2]v, nil,
		func(_ *rand.Rand) {
			weightMsg%[1]v%[2]v = defaultWeightMsg%[1]v%[2]v
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsg%[1]v%[2]v,
		%[3]vsimulation.SimulateMsg%[1]v%[2]v(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))

`, msg, typeName.UpperCamel, moduleName)

		content, err = xast.ModifyFunction(content, "WeightedOperations", xast.AppendFuncCode(replacementOp))
		if err != nil {
			return "", err
		}

		// add proposal simulation operations for msgs having an authority as signer.
		if strings.Contains(content, "ProposalMsgs") && strings.EqualFold(msgSigner.Original, "authority") {
			replacementOpMsg := fmt.Sprintf(`simulation.NewWeightedProposalMsg(
	opWeightMsg%[1]v%[2]v,
	defaultWeightMsg%[1]v%[2]v,
	func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
		%[3]vsimulation.SimulateMsg%[1]v%[2]v(am.authKeeper, am.bankKeeper, am.keeper)
		return nil
	},
),`, msg, typeName.UpperCamel, moduleName)
			content, err = xast.ModifyFunction(content, "ProposalMsgs", xast.AppendFuncCode(replacementOpMsg))
			if err != nil {
				return "", err
			}
		}
	}

	return content, nil
}
