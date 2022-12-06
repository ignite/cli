---
sidebar_position: 8
---

# CLI

A command line interface (CLI) lets you interact with your app after it is
running on a machine somewhere. Each module has its own namespace within the CLI
that gives it the ability to create and sign messages that are destined to be
handled by that module.

The CLI also comes with the ability to query the state of a module. When
combined with the rest of the app, the CLI lets you do things like generate keys
for a new account or check the status of an interaction you already had with the
application.

The CLI for the scavenge module is present in the `tx.go` and `query.go` files
in the `x/scavenge/client/cli/` directory.

- The `tx.go` file is for making transactions that contain messages that will
  ultimately update the state.
- The `query.go` file is for making queries let you read information from the
  state.

Both files use the [Cobra](https://github.com/spf13/cobra) library.

## Commit solution

Previously, you have scaffolded the code for handling the "commit solution"
message. You defined two fields for the message: "solution hash" and "solution
scavenger hash". This is the data you want to submit to the blockchain. On the
CLI, however, you want the user to be able to submit the solution as a string.
Make the appropriate changes to the CLI to calculate the hash of the solution
and the hash of the solution with the scavenger address.

```go title="x/scavenge/client/cli/tx_commit_solution.go"
package cli

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"

	"scavenge/x/scavenge/types"
)

func CmdCommitSolution() *cobra.Command {
	cmd := &cobra.Command{
		// pass a solution as the only argument
		// highlight-next-line
		Use:   "commit-solution [solution]",
		Short: "Broadcast message commit-solution",
		// set the number of arguments to 1
		// highlight-next-line
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			solution := args[0]

			// find a hash of the solution
			solutionHash := sha256.Sum256([]byte(solution))

			// convert the solution hash to string
			solutionHashString := hex.EncodeToString(solutionHash[:])

			// convert a scavenger address to string
			scavenger := clientCtx.GetFromAddress().String()

			// find the hash of solution and scavenger address
			solutionScavengerHash := sha256.Sum256([]byte(solution + scavenger))

			// convert the hash to string
			solutionScavengerHashString := hex.EncodeToString(solutionScavengerHash[:])

			// create a new message
			msg := types.NewMsgCommitSolution(
				clientCtx.GetFromAddress().String(),
				solutionHashString,
				solutionScavengerHashString,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			// broadcast the transaction with the message to the blockchain
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
```

Note that this file makes use of the `sha256` library for hashing the plain text
solutions into the scrambled hashes. This activity takes place on the client
side so the solutions are never leaked to any public entity that might want to
sneak a peak and steal the bounty reward associated with the scavenges. You can
also notice that the hashes are converted into hexadecimal representation to
make them easy to read as strings. Hashes are ultimately stored as hexadecimal
representations in the keeper.

## Submit scavenge

Modify the "submit scavenge" CLI to calculate the hash from a solution submitted
as a string.

```go title="x/scavenge/client/cli/tx_submit_scavenge.go"
package cli

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"

	"scavenge/x/scavenge/types"
)

func CmdSubmitScavenge() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "submit-scavenge [solution] [description] [reward]",
		Short: "Broadcast message submit-scavenge",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			// find a hash of the solution
			solutionHash := sha256.Sum256([]byte(args[0]))

			// convert the hash to string
			solutionHashString := hex.EncodeToString(solutionHash[:])
			argsDescription := args[1]
			argsReward := args[2]

			// create a new message
			msg := types.NewMsgSubmitScavenge(clientCtx.GetFromAddress().String(), solutionHashString, argsDescription, argsReward)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			// broadcast the transaction with the message to the blockchain
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
```
