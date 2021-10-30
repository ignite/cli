package typed

import (
	"fmt"

	"github.com/tendermint/starport/starport/pkg/multiformatname"
	"github.com/tendermint/starport/starport/pkg/placeholder"
)

func ModuleSimulationMsgModify(
	replacer placeholder.Replacer,
	content string,
	typeName multiformatname.Name,
	msgs ...string,
) string {
	if len(msgs) == 0 {
		msgs = append(msgs, "")
	}
	for _, msg := range msgs {
		// simulation constants
		templateConst := `opWeightMsg%[2]v%[3]v = "op_weight_msg_create_chain"
	// TODO: Determine the simulation weight value
	defaultWeightMsg%[2]v%[3]v int = 100

	%[1]v`
		replacementConst := fmt.Sprintf(templateConst, PlaceholderSimapConst, msg, typeName.UpperCamel)
		content = replacer.Replace(content, PlaceholderSimapConst, replacementConst)

		// simulation operations
		templateOp := `var weightMsg%[2]v%[3]v int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsg%[2]v%[3]v, &weightMsg%[2]v%[3]v, nil,
		func(_ *rand.Rand) {
			weightMsg%[2]v%[3]v = defaultWeightMsg%[2]v%[3]v
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsg%[2]v%[3]v,
		func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accounts []simtypes.Account, chainID string) (OperationMsg simtypes.OperationMsg, futureOps []simtypes.FutureOperation, err error) {

			// TODO: Handling the simulation

			simAccount, _ := simtypes.RandomAcc(r, accounts)
			msg := &types.Msg%[2]v%[3]v{}
			
			skipSimulation := true
			if skipSimulation {
				return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "skip %[2]v simulation message"), nil, nil
			}

			txCtx := simulation.OperationInput{
				R:               r,
				App:             app,
				TxGen:           simappparams.MakeTestEncodingConfig().TxConfig,
				Cdc:             nil,
				Msg:             msg,
				MsgType:         msg.Type(),
				Context:         ctx,
				SimAccount:      simAccount,
				ModuleName:      types.ModuleName,
				CoinsSpentInMsg: sdk.NewCoins(),
			}
			return simulation.GenAndDeliverTxWithRandFees(txCtx)
		},
	))

	%[1]v`
		replacementOp := fmt.Sprintf(templateOp, PlaceholderSimapOperation, msg, typeName.UpperCamel)
		content = replacer.Replace(content, PlaceholderSimapOperation, replacementOp)
	}
	return content
}
